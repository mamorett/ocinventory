package inventory

import (
	"context"
	"fmt"
	"sync"

	"github.com/oracle/oci-go-sdk/v65/common"
)

// ScanAll concurrently scans all compartments.
// concurrency controls the maximum number of parallel compartment goroutines.
func ScanAll(ctx context.Context, provider common.ConfigurationProvider, compartments []CompartmentInfo, concurrency int) ([]CompartmentResult, error) {
	if concurrency <= 0 {
		concurrency = 5
	}

	results := make([]CompartmentResult, len(compartments))
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for i, comp := range compartments {
		wg.Add(1)
		sem <- struct{}{}

		go func(idx int, c CompartmentInfo) {
			defer wg.Done()
			defer func() { <-sem }()

			res := CompartmentResult{Compartment: c}

			vms, err := ListInstances(ctx, provider, c)
			if err != nil {
				fmt.Printf("  [WARN] instances in %s: %v\n", c.Name, err)
			}
			res.VMs = vms

			vols, err := ListVolumes(ctx, provider, c)
			if err != nil {
				fmt.Printf("  [WARN] volumes in %s: %v\n", c.Name, err)
			}
			res.Volumes = vols

			vcns, err := ListVCNs(ctx, provider, c)
			if err != nil {
				fmt.Printf("  [WARN] VCNs in %s: %v\n", c.Name, err)
			}
			res.VCNs = vcns

			results[idx] = res
		}(i, comp)
	}

	wg.Wait()
	return results, nil
}
