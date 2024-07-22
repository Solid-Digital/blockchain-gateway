package factory

import (
	"context"
	"database/sql"
	"fmt"

	"bitbucket.org/unchain/ares/gen/orm"
	"bitbucket.org/unchain/ares/pkg/ares"
	"github.com/volatiletech/null"
)

func (f *Factory) OrganizationAndAWSBillingProvider(create bool, productCode, customerID string) (*orm.Organization, *orm.OrganizationBillingProvider) {
	org := f.Organization(create)

	billingInfoString := fmt.Sprintf(`{"awsCustomerId": "%s", "productCode": "%s"}`, customerID, productCode)

	obp := &orm.OrganizationBillingProvider{
		OrganizationID: org.ID,
		ProviderName:   ares.ProviderAWS,
		BillingInfo:    null.JSONFrom([]byte(billingInfoString)),
		CreatedByID:    org.CreatedByID,
		UpdatedByID:    org.UpdatedByID,
	}

	if !create {
		return org, obp
	}

	err := f.ares.DB.WrapTx(func(ctx context.Context, tx *sql.Tx) error {
		err := org.AddOrganizationBillingProviders(ctx, tx, true, obp)
		if err != nil {
			return err
		}

		return nil
	})

	f.suite.Require().NoError(err)

	return org, obp

}
