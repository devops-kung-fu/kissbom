// Package models contains any structs used throughout kissbom
package models

import (
	"bytes"
	"encoding/json"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/gocarina/gocsv"
	"gopkg.in/yaml.v3"
)

// Enumeration of valid output formats for kissbom.
const (
	OptionJSON       = "json"       // OptionDefault represents the default output format (kissbom format, json encoding).
	OptionYAML       = "yaml"       // OptionYAML represents the YAML output format (kissbom format, json encoding).
	OptionMinimal    = "minimal"    // OptionMinimal represents a minimal output format (kissbom format, but only Purls).
	OptionCompatible = "compatible" // OptionCompatible represents a compatible output format (CycloneDX formatted output, but only Purls).
	OptionCSV        = "csv"        // OptionCSV represents a CSV (Comma-Separated Values) output format (all kissbom elements)
)

// KissBOM represents a collection of packages.
type KissBOM struct {
	Packages []Package `json:"packages"` // Packages is a slice of Package structs, serialized as "packages" in JSON.
}

// Package represents information about a software package.
type Package struct {
	Purl      string `json:"purl" csv:"purl" yaml:"purl"`                                    // Purl is the Package URL, a unique identifier for the package.
	License   string `json:"license,omitempty" csv:"license" yaml:"license,omitempty"`       // License is the software license associated with the package, omitempty allows for optional serialization.
	Copyright string `json:"copyright,omitempty" csv:"copyright" yaml:"copyright,omitempty"` // Copyright is information about the package's copyright, omitempty allows for optional serialization.
	Notes     string `json:"notes,omitempty" csv:"notes" yaml:"notes,omitempty"`             // Notes is additional notes or comments about the package, omitempty allows for optional serialization.
}

// NewKissBOMFromCycloneDX creates a new KissBOM (Keep It Simple Software Bill of Materials)
// from a CycloneDX Bill of Materials (BOM).
// It iterates over the components in the CycloneDX BOM and constructs a simplified representation
// in the KissBOM format, with essential information such as Package URL and Description.
//
// Parameters:
//   - cdx: A pointer to a CycloneDX BOM containing information about software components.
//
// Returns:
//   - kissbom: A KissBOM representation derived from the CycloneDX BOM.
//
// NewKissBOMFromCycloneDX converts a CycloneDX BOM (Bill of Materials) to a KissBOM
// (KISS Build of Materials) by extracting relevant information from each component.
func NewKissBOMFromCycloneDX(cdx *cyclonedx.BOM) (kissbom KissBOM) {
	// Check if the Components list is nil
	if cdx.Components == nil {
		// If nil, return an empty KissBOM
		return
	}

	// Iterate through each component and populate the KissBOM Packages
	for _, component := range *cdx.Components {
		if component.PackageURL != "" {
			kissbom.Packages = append(kissbom.Packages, Package{
				Purl:      component.PackageURL,
				License:   extractLicense(component),
				Copyright: component.Copyright,
				Notes:     component.Description,
			})
		}
	}

	// Return the populated KissBOM
	return
}

// extractLicense extracts the license expression from a CycloneDX component.
// If the component has licenses, it returns the expression of the first license;
// otherwise, it returns an empty string.
func extractLicense(component cyclonedx.Component) string {
	if component.Licenses != nil && len(*component.Licenses) != 0 {
		return (*component.Licenses)[0].Expression
	}
	// If no licenses or licenses are nil, return an empty string
	return ""
}

// JSON converts the KissBOM struct to JSON format
func (k *KissBOM) JSON() ([]byte, error) {
	return json.MarshalIndent(k, "", "    ")
}

// YAML converts the KissBOM struct to YAML format
func (k *KissBOM) YAML() ([]byte, error) {
	return yaml.Marshal(k)
}

// CSV converts the KissBOM struct to CSV format using gocsv
func (k *KissBOM) CSV() ([]byte, error) {
	// Encode KissBOM to CSV
	c, err := gocsv.MarshalString(&k.Packages)
	return []byte(c), err
}

// Minimal converts the KissBOM struct to a JSON format with only the PURLs
func (k *KissBOM) Minimal() ([]byte, error) {
	kissbom := KissBOM{}
	for _, p := range k.Packages {
		kissbom.Packages = append(kissbom.Packages, Package{
			Purl: p.Purl,
		})
	}
	return json.MarshalIndent(kissbom, "", "    ")
}

// Compatible generates a CycloneDX Bill of Materials (BOM) based on the packages
// stored in the KissBOM instance. Each package's PackageURL is used to create
// corresponding CycloneDX components, and these components are added to the BOM.
// The resulting BOM is then encoded to a byte slice using the JSON format.
//
// Returns:
//   - The encoded BOM as a byte slice.
//   - An error if there was any issue during encoding.
func (k *KissBOM) Compatible() ([]byte, error) {
	bom := cyclonedx.NewBOM()
	components := []cyclonedx.Component{}
	for _, c := range k.Packages {
		component := cyclonedx.Component{
			PackageURL: c.Purl,
		}
		components = append(components, component)
	}
	bom.Components = &components

	return k.encodeBOM(bom)
}

// encodeBOM encodes a given CycloneDX BOM to a byte slice using the JSON format.
//
// Parameters:
//   - bom: The CycloneDX BOM to be encoded.
//
// Returns:
//   - The encoded BOM as a byte slice.
//   - An error if there was any issue during encoding.
func (k *KissBOM) encodeBOM(bom *cyclonedx.BOM) ([]byte, error) {
	var buf bytes.Buffer
	encoder := cyclonedx.NewBOMEncoder(&buf, cyclonedx.BOMFileFormatJSON)
	err := encoder.Encode(bom)
	return buf.Bytes(), err
}
