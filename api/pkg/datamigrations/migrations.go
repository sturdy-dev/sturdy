package datamigrations

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"getsturdy.com/api/pkg/changes"
	service_changes "getsturdy.com/api/pkg/changes/service"
	"getsturdy.com/api/pkg/codebases"
	service_codebases "getsturdy.com/api/pkg/codebases/service"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type Service struct {
	logger     *zap.Logger
	migrations []migration
}

func NewService(
	logger *zap.Logger,
	db *sqlx.DB,
	changesService *service_changes.Service,
	codebasesService *service_codebases.Service,
) *Service {
	return &Service{
		logger: logger.Named("datamigrations"),
		migrations: []migration{
			newChangesCodebaseIDCommitIDUniq(logger, db, changesService, codebasesService),
		},
	}
}

type migration interface {
	Run(context.Context) error
	Name() string
	Skip(context.Context) (bool, error)
}

func (s *Service) Run(ctx context.Context) error {
	for _, migration := range s.migrations {
		if skip, err := migration.Skip(ctx); err != nil {
			return fmt.Errorf("failed to check if migration %s should be skipped: %w", migration.Name(), err)
		} else if skip {
			s.logger.Warn("skipping migration", zap.String("migration", migration.Name()))
			continue
		}

		s.logger.Warn("running migration", zap.String("migration", migration.Name()))
		if err := migration.Run(ctx); err != nil {
			return fmt.Errorf("failed to run migration %s: %w", migration.Name(), err)
		}
		s.logger.Warn("migration completed", zap.String("migration", migration.Name()))
	}
	return nil
}

type changesCodebaseIDCommitIDUniq struct {
	logger *zap.Logger

	db               *sqlx.DB
	changesService   *service_changes.Service
	codebasesService *service_codebases.Service

	codebaseIDs      []codebases.ID
	codebaseIDsError error
	codebaseIDsOnce  sync.Once
}

func newChangesCodebaseIDCommitIDUniq(
	logger *zap.Logger,
	db *sqlx.DB,
	changesService *service_changes.Service,
	codebasesService *service_codebases.Service,
) *changesCodebaseIDCommitIDUniq {
	c := &changesCodebaseIDCommitIDUniq{
		logger:           logger,
		db:               db,
		changesService:   changesService,
		codebasesService: codebasesService,
	}
	c.logger = c.logger.With(zap.String("migration", c.Name()))
	return c
}

func (m *changesCodebaseIDCommitIDUniq) Skip(ctx context.Context) (bool, error) {
	codebaseIDs, err := m.loadCodebaseIDsWithDuplicates(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to load codebase ids: %w", err)
	}
	return len(codebaseIDs) == 0, nil
}

func (*changesCodebaseIDCommitIDUniq) Name() string {
	return "prepare changed_codebase_id_commit_id_uniq_idx"
}

func (c *changesCodebaseIDCommitIDUniq) Run(ctx context.Context) error {
	codebaseIDs, err := c.loadCodebaseIDsWithDuplicates(ctx)
	if err != nil {
		return fmt.Errorf("failed to load codebase ids: %w", err)
	}

	wg, ctx := errgroup.WithContext(ctx)
	for _, id := range codebaseIDs {
		id := id
		wg.Go(func() error {
			if err := c.migrateCodebase(ctx, id); err != nil {
				return fmt.Errorf("failed to migrate codebase %s: %w", id, err)
			}
			return nil
		})
	}

	return wg.Wait()
}

func (c *changesCodebaseIDCommitIDUniq) migrateCodebase(ctx context.Context, id codebases.ID) error {
	ambiguousChangeIDs, err := c.ambiguousChangeIDs(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to load duplicated changes for codebase %s: %w", id, err)
	}

	dbChangeSet, err := c.codebaseChangeSet(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to load changes for codebase %s: %w", id, err)
	}

	toDelete := []changes.ID{}
	for _, changeID := range ambiguousChangeIDs {
		if !dbChangeSet[changeID] {
			toDelete = append(toDelete, changeID)
		}
	}

	if len(toDelete) == 0 {
		return nil
	}

	c.logger.Warn("deleting duplicated changes", zap.Stringer("codebase", id), zap.Int("count", len(toDelete)))

	if _, err := c.db.ExecContext(ctx, `
		DELETE FROM changes
		WHERE id = ANY($1)
	`, pq.Array(toDelete)); err != nil {
		return fmt.Errorf("failed to delete duplicated changes for codebase %s: %w", id, err)
	}

	return nil
}

func (c *changesCodebaseIDCommitIDUniq) ambiguousChangeIDs(ctx context.Context, codebaseID codebases.ID) ([]changes.ID, error) {
	changeIDs := []changes.ID{}
	if err := c.db.SelectContext(ctx, &changeIDs, `
		SELECT changes.id
		FROM (
			SELECT codebase_id, commit_id
			FROM changes
			GROUP BY codebase_id, commit_id
			HAVING COUNT(1) > 1
		) as t
		JOIN changes ON 
			changes.codebase_id = t.codebase_id 
			AND changes.commit_id = t.commit_id
		WHERE changes.codebase_id = $1
	`, codebaseID); err != nil {
		return nil, fmt.Errorf("failed to load duplicated change ids: %w", err)
	}
	return changeIDs, nil
}

func (c *changesCodebaseIDCommitIDUniq) codebaseChangeSet(ctx context.Context, id codebases.ID) (map[changes.ID]bool, error) {
	codebase, err := c.codebasesService.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to load codebase %s: %w", id, err)
	}

	change, err := c.changesService.HeadChange(ctx, codebase)
	if err != nil {
		return nil, fmt.Errorf("failed to load head commit id: %w", err)
	}

	changeSet := map[changes.ID]bool{}
	for change.ParentChangeID != nil {
		if next, err := c.changesService.GetChangeByID(ctx, *change.ParentChangeID); errors.Is(err, sql.ErrNoRows) {
			break
		} else if err != nil {
			return nil, fmt.Errorf("failed to load parent change: %w", err)
		} else {
			changeSet[change.ID] = true
			changeSet[next.ID] = true
			change = next
		}
	}
	return changeSet, nil
}

func (c *changesCodebaseIDCommitIDUniq) loadCodebaseIDsWithDuplicates(ctx context.Context) ([]codebases.ID, error) {
	c.codebaseIDsOnce.Do(func() {
		codebaseIDs := []codebases.ID{}
		if err := c.db.SelectContext(ctx, &codebaseIDs, `
			SELECT codebase_id
			FROM (
				SELECT codebase_id, commit_id
				FROM changes
				GROUP BY codebase_id, commit_id
				HAVING COUNT(1) > 1
			) AS t
			WHERE codebase_id = 'b47c077a-e098-4eec-bb64-318f5a8571e3'
			GROUP by codebase_id
		`); err != nil {
			c.codebaseIDsError = err
		} else {
			c.codebaseIDs = codebaseIDs
		}
	})

	if c.codebaseIDsError != nil {
		return nil, c.codebaseIDsError
	}

	return c.codebaseIDs, nil
}
