package inventory

import (
	"context"
	"fmt"
	"sync"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
)

// imageCache avoids redundant GetImage calls across goroutines.
var imageCache sync.Map // map[string]string  (imageOCID -> "OS Name")

// ListInstances fetches all compute instances in one compartment and enriches
// them with the OS display name from the image (cached per image OCID).
func ListInstances(ctx context.Context, provider common.ConfigurationProvider, comp CompartmentInfo) ([]VMRecord, error) {
	client, err := core.NewComputeClientWithConfigurationProvider(provider)
	if err != nil {
		return nil, fmt.Errorf("create compute client: %w", err)
	}

	var vms []VMRecord
	var page *string

	for {
		resp, err := client.ListInstances(ctx, core.ListInstancesRequest{
			CompartmentId: common.String(comp.ID),
			Page:          page,
			Limit:         common.Int(100),
		})
		if err != nil {
			return nil, fmt.Errorf("list instances in %s: %w", comp.Name, err)
		}

		for _, inst := range resp.Items {
			// Skip terminated instances.
			if inst.LifecycleState == core.InstanceLifecycleStateTerminated {
				continue
			}

			name := "(unnamed)"
			if inst.DisplayName != nil {
				name = *inst.DisplayName
			}

			var ocpu float32
			var memGB float32
			if inst.ShapeConfig != nil {
				if inst.ShapeConfig.Ocpus != nil {
					ocpu = *inst.ShapeConfig.Ocpus
				}
				if inst.ShapeConfig.MemoryInGBs != nil {
					memGB = *inst.ShapeConfig.MemoryInGBs
				}
			}

			arch := "x86_64"
			if inst.ShapeConfig != nil && inst.ShapeConfig.ProcessorDescription != nil {
				arch = *inst.ShapeConfig.ProcessorDescription
			}

			shape := ""
			if inst.Shape != nil {
				shape = *inst.Shape
			}

			imageID := ""
			if inst.ImageId != nil {
				imageID = *inst.ImageId
			}

			osName := resolveImageOS(ctx, client, imageID)

			vms = append(vms, VMRecord{
				CompartmentID:   comp.ID,
				CompartmentName: comp.Name,
				Name:            name,
				State:           string(inst.LifecycleState),
				Shape:           shape,
				OCPU:            ocpu,
				MemoryGB:        memGB,
				OS:              osName,
				Architecture:    arch,
				ImageID:         imageID,
			})
		}

		if resp.OpcNextPage == nil {
			break
		}
		page = resp.OpcNextPage
	}

	return vms, nil
}

// resolveImageOS looks up the OS display name for an image OCID.
// Results are memoised in imageCache.
func resolveImageOS(ctx context.Context, client core.ComputeClient, imageID string) string {
	if imageID == "" {
		return "unknown"
	}

	if cached, ok := imageCache.Load(imageID); ok {
		return cached.(string)
	}

	resp, err := client.GetImage(ctx, core.GetImageRequest{
		ImageId: common.String(imageID),
	})
	if err != nil {
		imageCache.Store(imageID, "unknown")
		return "unknown"
	}

	name := "unknown"
	if resp.Image.DisplayName != nil {
		name = *resp.Image.DisplayName
	} else if resp.Image.OperatingSystem != nil {
		ver := ""
		if resp.Image.OperatingSystemVersion != nil {
			ver = " " + *resp.Image.OperatingSystemVersion
		}
		name = *resp.Image.OperatingSystem + ver
	}

	imageCache.Store(imageID, name)
	return name
}
