package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"mash/pkg/auth"
	"mash/pkg/organization"
	db_organization "mash/pkg/organization/db"
)

type Service struct {
	organizationRepository       db_organization.Repository
	organizationMemberRepository db_organization.MemberRepository
}

func New(
	organizationRepository db_organization.Repository,
	organizationMemberRepository db_organization.MemberRepository,
) *Service {
	return &Service{
		organizationRepository:       organizationRepository,
		organizationMemberRepository: organizationMemberRepository,
	}
}

func (svc *Service) ListByUserID(ctx context.Context, userID string) ([]*organization.Organization, error) {
	members, err := svc.organizationMemberRepository.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var res []*organization.Organization

	for _, m := range members {
		org, err := svc.organizationRepository.Get(ctx, m.OrganizationID)
		if err != nil {
			return nil, err
		}
		res = append(res, org)
	}

	return res, nil
}

func (svc *Service) Members(ctx context.Context, organizationID string) ([]*organization.Member, error) {
	members, err := svc.organizationMemberRepository.ListByOrganizationID(ctx, organizationID)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (svc *Service) Create(ctx context.Context, name string) (*organization.Organization, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, err
	}

	org := organization.Organization{
		ID:        uuid.NewString(),
		Name:      name,
		CreatedAt: time.Now(),
		CreatedBy: userID,
	}

	if err := svc.organizationRepository.Create(ctx, org); err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	// add the creator as a member
	if _, err := svc.AddMember(ctx, org.ID, userID); err != nil {
		return nil, fmt.Errorf("failed to invite creator to organization: %w", err)
	}

	return &org, nil
}

func (svc *Service) AddMember(ctx context.Context, orgID, userID string) (*organization.Member, error) {
	addedByUserID, err := auth.UserID(ctx)
	if err != nil {
		return nil, err
	}

	member := organization.Member{
		ID:             uuid.NewString(),
		OrganizationID: orgID,
		UserID:         userID,
		CreatedAt:      time.Now(),
		CreatedBy:      addedByUserID,
	}

	if err := svc.organizationMemberRepository.Create(ctx, member); err != nil {
		return nil, fmt.Errorf("failed to create member: %w", err)
	}

	return &member, nil
}

func (svc *Service) GetByID(ctx context.Context, organizationID string) (*organization.Organization, error) {
	member, err := svc.organizationRepository.Get(ctx, organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}
	return member, nil
}

func (svc *Service) GetMemberByUserIDAndOrganizationID(ctx context.Context, userID, organizationID string) (*organization.Member, error) {
	member, err := svc.organizationMemberRepository.GetByUserIDAndOrganizationID(ctx, userID, organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get member: %w", err)
	}
	return member, nil
}
