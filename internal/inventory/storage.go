package inventory

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
)

// ListVolumes fetches all block volumes and boot volumes in a compartment.
func ListVolumes(ctx context.Context, provider common.ConfigurationProvider, comp CompartmentInfo) ([]VolumeRecord, error) {
	client, err := core.NewBlockstorageClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, fmt.Errorf("create block storage client: %w", err)
	}

	var records []VolumeRecord

	// --- Block Volumes ---
	var page *string
	for {
		resp, err := client.ListVolumes(ctx, core.ListVolumesRequest{
			CompartmentId: common.String(comp.ID),
			Page:          page,
			Limit:         common.Int(100),
		})
		if err != nil {
			return nil, fmt.Errorf("list block volumes in %s: %w", comp.Name, err)
		}

		for _, v := range resp.Items {
			if v.LifecycleState == core.VolumeLifecycleStateTerminated {
				continue
			}
			name := "(unnamed)"
			if v.DisplayName != nil {
				name = *v.DisplayName
			}
			var sizeGB int64
			if v.SizeInGBs != nil {
				sizeGB = *v.SizeInGBs
			}
			records = append(records, VolumeRecord{
				CompartmentID:   comp.ID,
				CompartmentName: comp.Name,
				Type:            "Block",
				Name:            name,
				SizeGB:          sizeGB,
				State:           fmt.Sprintf("%s", v.LifecycleState),
			})
		}

		if resp.OpcNextPage == nil {
			break
		}
		page = resp.OpcNextPage
	}

	// --- Boot Volumes ---
	// Boot volumes require an availability domain filter; we iterate all ADs.
	adClient, err := newIdentityClient(provider)
	if err != nil {
		return nil, err
	}

	tenancyID, _ := provider.TenancyOCID()
	ads, err := listAvailabilityDomains(ctx, adClient, tenancyID)
	if err != nil {
		return nil, err
	}

	for _, ad := range ads {
		var bvPage *string
		for {
			resp, err := client.ListBootVolumes(ctx, core.ListBootVolumesRequest{
				CompartmentId:      common.String(comp.ID),
				AvailabilityDomain: common.String(ad),
				Page:               bvPage,
				Limit:              common.Int(100),
			})
			if err != nil {
				// Some compartments may not be visible in certain ADs; skip gracefully.
				break
			}

			for _, bv := range resp.Items {
				if bv.LifecycleState == core.BootVolumeLifecycleStateTerminated {
					continue
				}
				name := "(unnamed)"
				if bv.DisplayName != nil {
					name = *bv.DisplayName
				}
				var sizeGB int64
				if bv.SizeInGBs != nil {
					sizeGB = *bv.SizeInGBs
				}
				records = append(records, VolumeRecord{
					CompartmentID:   comp.ID,
					CompartmentName: comp.Name,
					Type:            "Boot",
					Name:            name,
					SizeGB:          sizeGB,
					State:           fmt.Sprintf("%s", bv.LifecycleState),
				})
			}

			if resp.OpcNextPage == nil {
				break
			}
			bvPage = resp.OpcNextPage
		}
	}

	return records, nil
}
