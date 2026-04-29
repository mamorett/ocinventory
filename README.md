# ☁️ ocinventory

[![Go Version](https://img.shields.io/github/go-mod/go-version/mamorett/ocinventory)](https://golang.org)
[![License](https://img.shields.io/github/license/mamorett/ocinventory)](LICENSE)

**ocinventory** is a lightweight, concurrent CLI tool designed to scan and report on your Oracle Cloud Infrastructure (OCI) resources. It traverses your entire compartment hierarchy to discover Compute instances, Block Volumes, and Virtual Cloud Networks (VCNs), generating clean, actionable reports in Markdown and CSV formats.

---

## ✨ Features

- 🚀 **Fast & Concurrent:** Scans multiple compartments in parallel using Go goroutines.
- 📂 **Recursive Discovery:** Automatically walks the compartment tree to find resources everywhere.
- 📊 **Multiple Formats:**
  - **Markdown:** A beautiful summary report perfect for documentation.
  - **CSV:** Detailed spreadsheets for VMs, Volumes, and VCNs, ready for Excel or analysis.
- 🛠️ **Cross-Platform:** Binaries available for Linux (AMD64/ARM64) and macOS (Apple Silicon).
- 🔒 **Secure:** Uses your existing OCI configuration and API keys.

---

## 🚀 Getting Started

### Prerequisites

- [Go 1.19+](https://golang.org/dl/) (if building from source).
- An OCI Account and [API Signing Key](https://docs.oracle.com/en-us/iaas/Content/API/Concepts/apisigningkey.htm).
- A valid OCI config file (usually at `~/.oci/config`).

### Installation

Clone the repository and build the binary:

```bash
git clone https://github.com/mamorett/ocinventory.git
cd ocinventory
make build
```

The binaries will be available in the `dist/` directory.

---

## 🛠️ Usage

Run the tool by specifying your OCI profile:

```bash
./dist/ocinventory-linux-amd64 -profile DEFAULT
```

### Options

| Flag | Description | Default |
| :--- | :--- | :--- |
| `-profile` | **(Required)** OCI config profile name | - |
| `-config` | Path to OCI config file | `~/.oci/config` |
| `-output` | Directory for output report files | `.` |
| `-concurrency` | Max parallel compartment scans | `5` |
| `-version` | Print version and exit | `false` |

### Example

```bash
# Scan using the 'PROD' profile with high concurrency
./ocinventory -profile PROD -concurrency 10 -output ./reports
```

---

## 📈 Reports

After a successful scan, `ocinventory` generates the following files in your output directory:

- `INVENTORY_<profile>_<timestamp>.md`: A high-level summary of your tenancy inventory.
- `INVENTORY_<profile>_<timestamp>_VMs.csv`: Details on all Compute instances (State, Shape, IP, etc.).
- `INVENTORY_<profile>_<timestamp>_Volumes.csv`: Details on Block Storage (Size, Attachment, etc.).
- `INVENTORY_<profile>_<timestamp>_VCNs.csv`: Network configuration overview.

---

## 🏗️ Building for Other Platforms

The included `Makefile` makes it easy to cross-compile:

```bash
make linux-amd64   # Linux x86_64
make linux-arm64   # Linux ARM64 (e.g., OCI Ampere A1)
make darwin-arm64  # macOS Apple Silicon
make all           # Build all targets
```

---

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

---

<p align="center">
  Built with ❤️ for OCI Users.
</p>
