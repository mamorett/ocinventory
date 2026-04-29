package report

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/mamorett/ocinventory/internal/inventory"
)

// WriteVMsCSV writes all VM records as CSV to w.
func WriteVMsCSV(w io.Writer, results []inventory.CompartmentResult) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	if err := cw.Write([]string{
		"Compartment", "Name", "State", "Shape", "OCPU", "RAM_GB", "OS", "Architecture",
	}); err != nil {
		return err
	}

	for _, cr := range results {
		for _, vm := range cr.VMs {
			if err := cw.Write([]string{
				cr.Compartment.Name,
				vm.Name,
				vm.State,
				vm.Shape,
				fmt.Sprintf("%.1f", vm.OCPU),
				fmt.Sprintf("%.1f", vm.MemoryGB),
				vm.OS,
				vm.Architecture,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

// WriteVolumesCSV writes all volume records (block + boot) as CSV to w.
func WriteVolumesCSV(w io.Writer, results []inventory.CompartmentResult) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	if err := cw.Write([]string{
		"Compartment", "Type", "Name", "Size_GB", "State",
	}); err != nil {
		return err
	}

	for _, cr := range results {
		for _, v := range cr.Volumes {
			if err := cw.Write([]string{
				cr.Compartment.Name,
				v.Type,
				v.Name,
				fmt.Sprintf("%d", v.SizeGB),
				v.State,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

// WriteVCNsCSV writes all VCN records as CSV to w.
func WriteVCNsCSV(w io.Writer, results []inventory.CompartmentResult) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	if err := cw.Write([]string{
		"Compartment", "VCN_Name", "CIDR_Blocks", "State",
	}); err != nil {
		return err
	}

	for _, cr := range results {
		for _, vcn := range cr.VCNs {
			if err := cw.Write([]string{
				cr.Compartment.Name,
				vcn.Name,
				strings.Join(vcn.CIDRBlocks, ";"),
				vcn.State,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}
