// Package report renders inventory results to Markdown documents.
package report

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/mamorett/ocinventory/internal/inventory"
)

// WriteMarkdown writes the full Markdown inventory report to w.
func WriteMarkdown(w io.Writer, result inventory.InventoryResult) error {
	now := time.Now().Format("2006-01-02 15:04:05 MST")

	// ---- Header ----
	fmt.Fprintf(w, "# OCI Inventory: %s\n", result.Profile)
	fmt.Fprintf(w, "Generated: %s\n", now)
	fmt.Fprintf(w, "Tenancy: %s\n\n", result.TenancyID)

	// ---- Summary table ----
	fmt.Fprintln(w, "## Summary")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "| Compartment | VMs | OCPUs | RAM (GB) | Block Vol (GB) | Boot Vol (GB) | VCNs |")
	fmt.Fprintln(w, "| :--- | :---: | :---: | :---: | :---: | :---: | :---: |")

	for _, cr := range result.Results {
		// Only emit rows where something exists.
		if len(cr.VMs) == 0 && len(cr.Volumes) == 0 && len(cr.VCNs) == 0 {
			continue
		}

		var totalOCPU float32
		var totalRAM float32
		for _, vm := range cr.VMs {
			totalOCPU += vm.OCPU
			totalRAM += vm.MemoryGB
		}

		var blockGB, bootGB int64
		for _, v := range cr.Volumes {
			if v.Type == "Block" {
				blockGB += v.SizeGB
			} else {
				bootGB += v.SizeGB
			}
		}

		fmt.Fprintf(w, "| %s | %d | %.1f | %.1f | %d | %d | %d |\n",
			cr.Compartment.Name,
			len(cr.VMs),
			totalOCPU,
			totalRAM,
			blockGB,
			bootGB,
			len(cr.VCNs),
		)
	}

	// ---- Per-compartment detail ----
	fmt.Fprintln(w, "\n## Details")

	for _, cr := range result.Results {
		if len(cr.VMs) == 0 && len(cr.Volumes) == 0 && len(cr.VCNs) == 0 {
			continue
		}

		fmt.Fprintf(w, "\n### Compartment: %s\n\n", cr.Compartment.Name)

		// VMs
		fmt.Fprintln(w, "#### Virtual Machines")
		fmt.Fprintln(w, "| Name | State | Shape | OCPU | RAM (GB) | OS | Architecture |")
		fmt.Fprintln(w, "| :--- | :--- | :--- | :---: | :---: | :--- | :--- |")
		if len(cr.VMs) == 0 {
			fmt.Fprintln(w, "| — | — | — | — | — | — | — |")
		} else {
			for _, vm := range cr.VMs {
				fmt.Fprintf(w, "| %s | %s | %s | %.1f | %.1f | %s | %s |\n",
					vm.Name, vm.State, vm.Shape, vm.OCPU, vm.MemoryGB, vm.OS, vm.Architecture)
			}
		}
		fmt.Fprintln(w)

		// Boot Volumes
		boots := filterVolumes(cr.Volumes, "Boot")
		fmt.Fprintln(w, "#### Storage — Boot Volumes")
		fmt.Fprintln(w, "| Name | Size (GB) | State |")
		fmt.Fprintln(w, "| :--- | :---: | :--- |")
		if len(boots) == 0 {
			fmt.Fprintln(w, "| — | — | — |")
		} else {
			for _, v := range boots {
				fmt.Fprintf(w, "| %s | %d | %s |\n", v.Name, v.SizeGB, v.State)
			}
		}
		fmt.Fprintln(w)

		// Block Volumes
		blocks := filterVolumes(cr.Volumes, "Block")
		fmt.Fprintln(w, "#### Storage — Block Volumes")
		fmt.Fprintln(w, "| Name | Size (GB) | State |")
		fmt.Fprintln(w, "| :--- | :---: | :--- |")
		if len(blocks) == 0 {
			fmt.Fprintln(w, "| — | — | — |")
		} else {
			for _, v := range blocks {
				fmt.Fprintf(w, "| %s | %d | %s |\n", v.Name, v.SizeGB, v.State)
			}
		}
		fmt.Fprintln(w)

		// VCNs
		fmt.Fprintln(w, "#### Networking (VCNs)")
		fmt.Fprintln(w, "| VCN Name | CIDR Blocks | State |")
		fmt.Fprintln(w, "| :--- | :--- | :--- |")
		if len(cr.VCNs) == 0 {
			fmt.Fprintln(w, "| — | — | — |")
		} else {
			for _, vcn := range cr.VCNs {
				fmt.Fprintf(w, "| %s | %s | %s |\n",
					vcn.Name,
					strings.Join(vcn.CIDRBlocks, ", "),
					vcn.State)
			}
		}

		fmt.Fprintln(w, "\n---")
	}

	return nil
}

func filterVolumes(vols []inventory.VolumeRecord, volType string) []inventory.VolumeRecord {
	var out []inventory.VolumeRecord
	for _, v := range vols {
		if v.Type == volType {
			out = append(out, v)
		}
	}
	return out
}
