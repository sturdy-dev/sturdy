package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/organization"
	db_organization "getsturdy.com/api/pkg/organization/db"
	"getsturdy.com/api/pkg/shortid"
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
		return nil, fmt.Errorf("could not list organizations by user id: %w", err)
	}

	var res []*organization.Organization

	for _, m := range members {
		org, err := svc.organizationRepository.Get(ctx, m.OrganizationID)
		if err != nil {
			return nil, fmt.Errorf("could not get organization by membership: %w", err)
		}
		res = append(res, org)
	}

	return res, nil
}

func (svc *Service) Members(ctx context.Context, organizationID string) ([]*organization.Member, error) {
	members, err := svc.organizationMemberRepository.ListByOrganizationID(ctx, organizationID)
	if err != nil {
		return nil, fmt.Errorf("could not get members: %w", err)
	}
	return members, nil
}

func (svc *Service) GetMember(ctx context.Context, organizationID, userID string) (*organization.Member, error) {
	member, err := svc.organizationMemberRepository.GetByUserIDAndOrganizationID(ctx, userID, organizationID)
	if err != nil {
		return nil, fmt.Errorf("could not get member: %w", err)
	}
	return member, nil
}

func (svc *Service) Create(ctx context.Context, name string) (*organization.Organization, error) {
	userID, err := auth.UserID(ctx)
	if err != nil {
		return nil, err
	}

	org := organization.Organization{
		ID:        uuid.NewString(),
		ShortID:   organization.ShortOrganizationID(shortid.New()),
		Name:      name,
		CreatedAt: time.Now(),
		CreatedBy: userID,
	}

	if err := svc.organizationRepository.Create(ctx, org); err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	// add the creator as a member
	if _, err := svc.AddMember(ctx, org.ID, userID, userID); err != nil {
		return nil, fmt.Errorf("failed to invite creator to organization: %w", err)
	}

	return &org, nil
}

func (svc *Service) AddMember(ctx context.Context, orgID, userID, addedByUserID string) (*organization.Member, error) {
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

func (svc *Service) GetByShortID(ctx context.Context, shortID organization.ShortOrganizationID) (*organization.Organization, error) {
	member, err := svc.organizationRepository.GetByShortID(ctx, shortID)
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

// TODO: Only allow calls to GetFirst from self hosted installations
func (svc *Service) GetFirst(ctx context.Context) (*organization.Organization, error) {
	org, err := svc.organizationRepository.GetFirst(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get first organization: %w", err)
	}
	return org, nil
}
