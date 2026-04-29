package inventory

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/identity"
)

// ListAllCompartments returns every active compartment under the tenancy root,
// including the root tenancy itself, using a recursive subtree query.
func ListAllCompartments(ctx context.Context, provider common.ConfigurationProvider) ([]CompartmentInfo, error) {
	tenancyID, err := provider.TenancyOCID()
	if err != nil {
		return nil, fmt.Errorf("get tenancy OCID: %w", err)
	}

	client, err := identity.NewIdentityClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, fmt.Errorf("create identity client: %w", err)
	}

	var compartments []CompartmentInfo

	// Always add the root tenancy first.
	compartments = append(compartments, CompartmentInfo{
		ID:   tenancyID,
		Name: "Root_Tenancy",
	})

	var page *string
	for {
		resp, err := client.ListCompartments(ctx, identity.ListCompartmentsRequest{
			CompartmentId:          common.String(tenancyID),
			CompartmentIdInSubtree: common.Bool(true),
			LifecycleState:         identity.CompartmentLifecycleStateActive,
			Page:                   page,
			Limit:                  common.Int(100),
		})
		if err != nil {
			return nil, fmt.Errorf("list compartments: %w", err)
		}

		for _, c := range resp.Items {
			name := "(unnamed)"
			if c.Name != nil {
				name = *c.Name
			}
			compartments = append(compartments, CompartmentInfo{
				ID:   *c.Id,
				Name: name,
			})
		}

		if resp.OpcNextPage == nil {
			break
		}
		page = resp.OpcNextPage
	}

	return compartments, nil
}
