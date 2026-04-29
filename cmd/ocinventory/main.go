package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mamorett/ocinventory/internal/config"
	"github.com/mamorett/ocinventory/internal/inventory"
	"github.com/mamorett/ocinventory/internal/report"
)

var version = "dev"

func main() {
	profileFlag := flag.String("profile", "", "OCI config profile name (required)")
	configFlag := flag.String("config", "", "Path to OCI config file (default: ~/.oci/config)")
	outputFlag := flag.String("output", ".", "Directory for output report files")
	concFlag := flag.Int("concurrency", 5, "Max parallel compartment goroutines")
	versionFlag := flag.Bool("version", false, "Print version and exit")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "ocinventory %s\n\n", version)
		fmt.Fprintf(os.Stderr, "Usage: ocinventory -profile <name> [options]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *versionFlag {
		fmt.Println("ocinventory", version)
		os.Exit(0)
	}

	if *profileFlag == "" {
		fmt.Fprintln(os.Stderr, "error: -profile is required")
		flag.Usage()
		os.Exit(1)
	}

	// Build OCI config provider.
	provider, err := config.NewProvider(*profileFlag, *configFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	tenancyID, _ := provider.TenancyOCID()
	fmt.Printf("Profile  : %s\n", *profileFlag)
	fmt.Printf("Tenancy  : %s\n", tenancyID)

	ctx := context.Background()

	// 1. List all compartments.
	fmt.Println("\nDiscovering compartments...")
	compartments, err := inventory.ListAllCompartments(ctx, provider)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error listing compartments: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Found %d compartment(s).\n", len(compartments))

	// 2. Scan all compartments concurrently.
	fmt.Printf("\nScanning resources (concurrency=%d)...\n", *concFlag)
	results, err := inventory.ScanAll(ctx, provider, compartments, *concFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error scanning: %v\n", err)
		os.Exit(1)
	}

	fullResult := inventory.InventoryResult{
		Profile:   *profileFlag,
		TenancyID: tenancyID,
		Results:   results,
	}

	// 3. Write reports.
	ts := time.Now().Format("20060102_1504")
	base := fmt.Sprintf("INVENTORY_%s_%s", *profileFlag, ts)

	if err := os.MkdirAll(*outputFlag, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "error creating output dir: %v\n", err)
		os.Exit(1)
	}

	// Markdown report.
	mdPath := filepath.Join(*outputFlag, base+".md")
	if err := writeFile(mdPath, func(f *os.File) error {
		return report.WriteMarkdown(f, fullResult)
	}); err != nil {
		fmt.Fprintf(os.Stderr, "error writing markdown: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\nMarkdown report : %s\n", mdPath)

	// CSV — VMs.
	vmsPath := filepath.Join(*outputFlag, base+"_VMs.csv")
	if err := writeFile(vmsPath, func(f *os.File) error {
		return report.WriteVMsCSV(f, results)
	}); err != nil {
		fmt.Fprintf(os.Stderr, "error writing VMs CSV: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("VMs CSV         : %s\n", vmsPath)

	// CSV — Volumes.
	volsPath := filepath.Join(*outputFlag, base+"_Volumes.csv")
	if err := writeFile(volsPath, func(f *os.File) error {
		return report.WriteVolumesCSV(f, results)
	}); err != nil {
		fmt.Fprintf(os.Stderr, "error writing Volumes CSV: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Volumes CSV     : %s\n", volsPath)

	// CSV — VCNs.
	vcnsPath := filepath.Join(*outputFlag, base+"_VCNs.csv")
	if err := writeFile(vcnsPath, func(f *os.File) error {
		return report.WriteVCNsCSV(f, results)
	}); err != nil {
		fmt.Fprintf(os.Stderr, "error writing VCNs CSV: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("VCNs CSV        : %s\n", vcnsPath)

	fmt.Println("\nDone.")
}

func writeFile(path string, fn func(*os.File) error) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return fn(f)
}
