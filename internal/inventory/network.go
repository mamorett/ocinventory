package inventory

import (
	"context"
	"fmt"
	"strings"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
)

// ListVCNs fetches all VCNs in a compartment.
func ListVCNs(ctx context.Context, provider common.ConfigurationProvider, comp CompartmentInfo) ([]VCNRecord, error) {
	client, err := core.NewVirtualNetworkClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, fmt.Errorf("create network client: %w", err)
	}

	var records []VCNRecord
	var page *string

	for {
		resp, err := client.ListVcns(ctx, core.ListVcnsRequest{
			CompartmentId: common.String(comp.ID),
			Page:          page,
			Limit:         common.Int(100),
		})
		if err != nil {
			return nil, fmt.Errorf("list VCNs in %s: %w", comp.Name, err)
		}

		for _, vcn := range resp.Items {
			if vcn.LifecycleState == core.VcnLifecycleStateTerminating ||
				vcn.LifecycleState == core.VcnLifecycleStateTerminated {
				continue
			}

			name := "(unnamed)"
			if vcn.DisplayName != nil {
				name = *vcn.DisplayName
			}

			// Collect all CIDR blocks; CidrBlocks (slice) is preferred over
			// the legacy CidrBlock (single string).
			var cidrs []string
			if len(vcn.CidrBlocks) > 0 {
				cidrs = vcn.CidrBlocks
			} else if vcn.CidrBlock != nil && *vcn.CidrBlock != "" {
				cidrs = []string{*vcn.CidrBlock}
			}

			records = append(records, VCNRecord{
				CompartmentID:   comp.ID,
				CompartmentName: comp.Name,
				Name:            name,
				CIDRBlocks:      cidrs,
				State:           strings.ToUpper(string(vcn.LifecycleState)),
			})
		}

		if resp.OpcNextPage == nil {
			break
		}
		page = resp.OpcNextPage
	}

	return records, nil
}
