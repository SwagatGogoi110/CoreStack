package organizations

import (
	"context"
	"fmt"

	"github.com/hectorvent/cloudstack/internal/services/aws/organizations/model"
	"github.com/hectorvent/cloudstack/internal/storage"
)

type OrganizationsService struct {
	orgStore     storage.Backend[string, *model.Organization]
	accountStore storage.Backend[string, *model.Account]
}

func NewOrganizationsService(factory *storage.Factory) (*OrganizationsService, error) {
	orgStore, _ := storage.CreateAccountAware[*model.Organization](factory, "organizations", "organizations.json", "wal")
	accountStore, _ := storage.CreateAccountAware[*model.Account](factory, "organizations", "accounts.json", "wal")

	return &OrganizationsService{
		orgStore:     orgStore,
		accountStore: accountStore,
	}, nil
}

func (s *OrganizationsService) CreateOrganization(ctx context.Context, featureSet string) (*model.Organization, error) {
	if orgs, _ := s.orgStore.Scan(ctx, func(k string) bool { return true }); len(orgs) > 0 {
		return nil, fmt.Errorf("AlreadyInOrganizationException: The AWS account is already part of an organization")
	}

	id := "o-rootorg"
	org := &model.Organization{
		Id:               id,
		Arn:              fmt.Sprintf("arn:aws:organizations::000000000000:organization/%s", id),
		FeatureSet:       featureSet,
		MasterAccountId:  "000000000000",
		MasterAccountArn: "arn:aws:organizations::000000000000:account/o-rootorg/000000000000",
		MasterAccountEmail: "admin@localhost",
	}

	s.orgStore.Put(ctx, id, org)
	return org, nil
}

func (s *OrganizationsService) DescribeOrganization(ctx context.Context) (*model.Organization, error) {
	orgs, _ := s.orgStore.Scan(ctx, func(k string) bool { return true })
	if len(orgs) == 0 {
		return nil, fmt.Errorf("AWSOrganizationsNotInUseException: Your account is not a member of an organization")
	}
	return orgs[0], nil
}

func (s *OrganizationsService) ListAccounts(ctx context.Context) ([]*model.Account, error) {
	return s.accountStore.Scan(ctx, func(k string) bool { return true })
}
