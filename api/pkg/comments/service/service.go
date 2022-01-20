package service

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/change"
	db_comments "getsturdy.com/api/pkg/comments/db"
)

type Service struct {
	commentRepo db_comments.Repository
}

func New(commentRepo db_comments.Repository) *Service {
	return &Service{
		commentRepo: commentRepo,
	}
}

func (s *Service) MoveCommentsFromWorkspaceToChange(ctx context.Context, workspaceID string, changeID change.ID) error {
	// Move all live comments on this workspace to the new change
	comments, err := s.commentRepo.GetByWorkspace(workspaceID)
	if err != nil {
		return fmt.Errorf("failed to get comments in workspace: %w", err)
	}
	for _, comment := range comments {
		comment.WorkspaceID = nil
		comment.ChangeID = &changeID
		if err := s.commentRepo.Update(comment); err != nil {
			return fmt.Errorf("failed to update comment: %w", err)
		}
	}
	return nil
}
