package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"getsturdy.com/api/pkg/analytics"
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	"getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/organization"
	db_organization "getsturdy.com/api/pkg/organization/db"
	"getsturdy.com/api/pkg/shortid"
	"getsturdy.com/api/pkg/users"

	"github.com/google/uuid"
)

type Service struct {
	organizationRepository       db_organization.Repository
	organizationMemberRepository db_organization.MemberRepository
	analyticsService             *service_analytics.Service
}

func New(
	organizationRepository db_organization.Repository,
	organizationMemberRepository db_organization.MemberRepository,
	analyticsService *service_analytics.Service,
) *Service {
	return &Service{
		organizationRepository:       organizationRepository,
		organizationMemberRepository: organizationMemberRepository,
		analyticsService:             analyticsService,
	}
}

func (svc *Service) ListByUserID(ctx context.Context, userID users.ID) ([]*organization.Organization, error) {
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

func (svc *Service) CanAccess(ctx context.Context, userID users.ID, organizationID string) (bool, error) {
	_, err := svc.organizationMemberRepository.GetByUserIDAndOrganizationID(ctx, userID, organizationID)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return false, nil
	case err == nil:
		return true, nil
	default:
		return false, fmt.Errorf("could not get member: %w", err)
	}
}

func (svc *Service) GetMember(ctx context.Context, organizationID string, userID users.ID) (*organization.Member, error) {
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

	svc.analyticsService.IdentifyOrganization(ctx, &org)
	svc.analyticsService.Capture(ctx, "create organization", analytics.OrganizationID(org.ID))

	return &org, nil
}

func (svc *Service) Update(ctx context.Context, organizationID string, newName string) (*organization.Organization, error) {

	org, err := svc.organizationRepository.Get(ctx, organizationID)
	if err != nil {
		return nil, fmt.Errorf("could not get organization by this ID: %w", err)
	}

	org.Name = newName

	err = svc.organizationRepository.Update(ctx, org)
	if err != nil {
		return nil, fmt.Errorf("could not update organization: %w", err)
	}
	return org, nil
}

func (svc *Service) AddMember(ctx context.Context, orgID string, userID, addedByUserID users.ID) (*organization.Member, error) {
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

	svc.analyticsService.Capture(ctx, "add member to organization",
		analytics.OrganizationID(orgID),
		analytics.Property("user_id", userID),
	)

	return &member, nil
}

func (svc *Service) RemoveMember(ctx context.Context, orgID string, userID users.ID, deletedByUserID users.ID) error {
	member, err := svc.organizationMemberRepository.GetByUserIDAndOrganizationID(ctx, userID, orgID)
	if err != nil {
		return fmt.Errorf("could not get member: %w", err)
	}

	t := time.Now()
	member.DeletedAt = &t
	member.DeletedBy = &deletedByUserID

	if err := svc.organizationMemberRepository.Update(ctx, member); err != nil {
		return fmt.Errorf("could not update member: %w", err)
	}

	svc.analyticsService.Capture(ctx, "remove member from organization",
		analytics.OrganizationID(orgID),
		analytics.Property("user_id", userID),
	)

	return nil
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

func (svc *Service) GetMemberByUserIDAndOrganizationID(ctx context.Context, userID users.ID, organizationID string) (*organization.Member, error) {
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
