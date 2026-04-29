// Package inventory contains the core data types shared across all sub-packages.
package inventory

// CompartmentInfo holds the minimal info needed for scanning.
type CompartmentInfo struct {
	ID   string
	Name string
}

// VMRecord represents a single compute instance.
type VMRecord struct {
	CompartmentID   string
	CompartmentName string
	Name            string
	State           string
	Shape           string
	OCPU            float32
	MemoryGB        float32
	OS              string
	Architecture    string
	ImageID         string
}

// VolumeRecord represents either a Block Volume or a Boot Volume.
type VolumeRecord struct {
	CompartmentID   string
	CompartmentName string
	Type            string // "Boot" | "Block"
	Name            string
	SizeGB          int64
	State           string
}

// VCNRecord represents a Virtual Cloud Network.
type VCNRecord struct {
	CompartmentID   string
	CompartmentName string
	Name            string
	CIDRBlocks      []string
	State           string
}

// CompartmentResult is the collected inventory for one compartment.
type CompartmentResult struct {
	Compartment CompartmentInfo
	VMs         []VMRecord
	Volumes     []VolumeRecord
	VCNs        []VCNRecord
	Error       error
}

// InventoryResult is the full inventory across all compartments.
type InventoryResult struct {
	Profile    string
	TenancyID  string
	Results    []CompartmentResult
}
