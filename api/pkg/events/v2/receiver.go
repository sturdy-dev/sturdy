package events

import (
	"context"
	"fmt"

	db_codebase "getsturdy.com/api/pkg/codebase/db"
	db_organization "getsturdy.com/api/pkg/organization/db"
	"getsturdy.com/api/pkg/users"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
)

type receiver struct {
	UserIDs         []users.ID
	WorkspaceIDs    []string
	CodebaseIDs     []string
	OrganizationIDs []string
}

func Workspace(workspaceID string) *receiver {
	return &receiver{
		WorkspaceIDs: []string{workspaceID},
	}
}

func User(userID users.ID) *receiver {
	return &receiver{
		UserIDs: []users.ID{userID},
	}
}

func Codebase(codebaseID string) *receiver {
	return &receiver{
		CodebaseIDs: []string{codebaseID},
	}
}

func Organization(organizationID string) *receiver {
	return &receiver{
		OrganizationIDs: []string{organizationID},
	}
}

func (r *receiver) Topics(
	ctx context.Context,
	codebaseUserRepo db_codebase.CodebaseUserRepository,
	workspaceRepo db_workspaces.WorkspaceReader,
	organizationMemberRepoRepo db_organization.MemberRepository,
) (map[Topic]bool, error) {
	topics := map[Topic]bool{}
	for _, userID := range r.UserIDs {
		topics[userTopic(userID)] = true
	}

	for _, workspaceID := range r.WorkspaceIDs {
		ws, err := workspaceRepo.Get(workspaceID)
		if err != nil {
			return nil, fmt.Errorf("failed to get workspace %s: %w", workspaceID, err)
		}
		r.CodebaseIDs = append(r.CodebaseIDs, ws.CodebaseID)
	}

	for _, codebaseID := range r.CodebaseIDs {
		members, err := codebaseUserRepo.GetByCodebase(codebaseID)
		if err != nil {
			return nil, fmt.Errorf("failed to get codebase members: %w", err)
		}
		for _, member := range members {
			topics[userTopic(member.UserID)] = true
		}
	}

	for _, organizationID := range r.OrganizationIDs {
		members, err := organizationMemberRepoRepo.ListByOrganizationID(ctx, organizationID)
		if err != nil {
			return nil, fmt.Errorf("failed to get organization members: %w", err)
		}
		for _, member := range members {
			topics[userTopic(member.UserID)] = true
		}
	}

	return topics, nil
}
