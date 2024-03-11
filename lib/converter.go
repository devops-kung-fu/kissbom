// Package lib contains all of the core functionality of the application
package lib

import (
	"bytes"
	"fmt"
	"log"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/spf13/afero"

	"github.com/devops-kung-fu/kissbom/models"
)

// Converter represents a utility for file conversion.
type Converter struct {
	Afs            *afero.Afero // Afero file system abstraction for file operations.
	OutputFileName string       // Name of the output file.
	OutputFormat   string       // Desired output format.
}

// NewConverter creates a new instance of the Converter with default settings.
// It initializes the Afs field with an Afero instance using the default operating system file system.
//
// Returns:
//   - A pointer to the newly created Converter instance.
func NewConverter() *Converter {
	return &Converter{
		Afs: &afero.Afero{Fs: afero.NewOsFs()},
	}
}

// Convert executes the conversion of the provided CycloneDX file to a KissBOM
func (c *Converter) Convert(filename string) error {
	log.Printf("converting: %v", filename)

	source, err := c.Afs.ReadFile(filename)
	if err != nil {
		return err
	}

	log.Printf("bytes: %v", len(source))

	kissbom, err := c.transform(source)
	if err != nil {
		return err
	}

	return c.writeToFile(kissbom)
}

// transform takes a byte slice representing a CycloneDX Bill of Materials (BOM) in JSON format,
// decodes it into a CycloneDX BOM object, and then transforms it into a KissBOM object along
// with a filename. Any decoding errors are returned as an error.
func (c *Converter) transform(source []byte) (kissbom models.KissBOM, err error) {
	var cdx cyclonedx.BOM

	err = cyclonedx.NewBOMDecoder(bytes.NewReader(source), cyclonedx.BOMFileFormatJSON).Decode(&cdx)
	if err != nil {
		return
	}

	log.Println("transformed to kissbom")

	c.OutputFileName = c.buildOutputFilename(&cdx)

	return models.NewKissBOMFromCycloneDX(&cdx), nil
}

// buildOutputFilename builds the output filename from the provided CycloneDX BOM
//
// The filename should be used to document the subject of the SBoM including optionally
// the product or component name, the SBoM author name, and an ISO 8601 timestamp of when
// this SBoM was last modified. Filenames should be lowercase and contain no space and should prefer using "-", "_" and "." as separator between words.
func (c *Converter) buildOutputFilename(cdx *cyclonedx.BOM) string {
	if cdx.Metadata != nil && cdx.Metadata.Component != nil {
		subject := cdx.Metadata.Component.Name
		publisher := cdx.Metadata.Component.Publisher
		timestamp := cdx.Metadata.Timestamp
		c.OutputFileName = fmt.Sprintf("%s_%s_%s", subject, publisher, timestamp)
	}
	return c.OutputFileName
}

// Function to write the KissBOM to a file based on the specified output format
func (c *Converter) writeToFile(kissbom models.KissBOM) error {
	var outputData []byte
	var err error

	switch c.OutputFormat {
	case models.OptionJSON:
		outputData, err = kissbom.JSON()
		c.OutputFileName += ".json"
	case models.OptionYAML:
		outputData, err = kissbom.YAML()
		c.OutputFileName += ".yaml"
	case models.OptionCSV:
		outputData, err = kissbom.CSV()
		c.OutputFileName += ".csv"
	case models.OptionMinimal:
		outputData, err = kissbom.Minimal()
		c.OutputFileName += ".json"
	case models.OptionCompatible:
		outputData, err = kissbom.Compatible()
		c.OutputFileName += ".cyclonedx.json"
	default:
		err = fmt.Errorf("unsupported output format: %s", c.OutputFormat)
	}

	if err != nil {
		return err
	}

	log.Printf("final bytes: %v", len(outputData))

	// Use afero to write the output data to the file
	err = afero.WriteFile(c.Afs, c.OutputFileName, outputData, 0644)
	log.Printf("saved: %v", c.OutputFileName)
	return err
}
