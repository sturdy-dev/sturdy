package service

import (
	"context"

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
