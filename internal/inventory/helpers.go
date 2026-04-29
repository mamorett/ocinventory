package inventory

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/identity"
)

// newIdentityClient is a small helper used by storage.go to list ADs.
func newIdentityClient(provider common.ConfigurationProvider) (identity.IdentityClient, error) {
	client, err := identity.NewIdentityClientWithConfigurationProvider(provider)
	if err != nil {
		return identity.IdentityClient{}, fmt.Errorf("create identity client: %w", err)
	}
	return client, nil
}

// listAvailabilityDomains returns the names of all availability domains for the tenancy.
func listAvailabilityDomains(ctx context.Context, client identity.IdentityClient, tenancyID string) ([]string, error) {
	resp, err := client.ListAvailabilityDomains(ctx, identity.ListAvailabilityDomainsRequest{
		CompartmentId: common.String(tenancyID),
	})
	if err != nil {
		return nil, fmt.Errorf("list availability domains: %w", err)
	}

	var ads []string
	for _, ad := range resp.Items {
		if ad.Name != nil {
			ads = append(ads, *ad.Name)
		}
	}
	return ads, nil
}
