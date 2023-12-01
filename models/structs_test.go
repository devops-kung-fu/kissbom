package models

import (
	"encoding/json"
	"testing"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestNewKissBOMFromCycloneDX(t *testing.T) {
	// Mock CycloneDX BOM for testing
	mockComponents := []cyclonedx.Component{
		{
			PackageURL:  "pkg:pypi/requests@2.26.0",
			Description: "Package 1 description",
		},
		{
			PackageURL:  "pkg:pypi/requests@2.26.1",
			Description: "Package 2 description",
		},
	}

	mockBOM := &cyclonedx.BOM{
		Components: &mockComponents,
	}

	kissBOM := NewKissBOMFromCycloneDX(mockBOM)

	// Check if KissBOM Packages are correctly populated
	assert.Len(t, kissBOM.Packages, 2)
	assert.Equal(t, "pkg:pypi/requests@2.26.0", kissBOM.Packages[0].Purl)
	assert.Equal(t, "Package 1 description", kissBOM.Packages[0].Notes)
	assert.Equal(t, "pkg:pypi/requests@2.26.1", kissBOM.Packages[1].Purl)
	assert.Equal(t, "Package 2 description", kissBOM.Packages[1].Notes)

	mockBOM.Components = nil

	kissBOM = NewKissBOMFromCycloneDX(mockBOM)
	assert.Len(t, kissBOM.Packages, 0)
}

func TestExtractLicense(t *testing.T) {
	// Test case 1: Component with a valid l

	component1 := cyclonedx.Component{
		Licenses: &cyclonedx.Licenses{
			cyclonedx.LicenseChoice{Expression: "(AFL-2.1 OR BSD-3-Clause)"},
		},
	}
	result1 := extractLicense(component1)
	if result1 != "(AFL-2.1 OR BSD-3-Clause)" {
		t.Errorf("Expected license 'MIT', got '%s'", result1)
	}

	// Test case 2: Component with no licenses
	component2 := cyclonedx.Component{
		Licenses: nil,
	}
	result2 := extractLicense(component2)
	if result2 != "" {
		t.Errorf("Expected empty license, got '%s'", result2)
	}

	// Test case 3: Component with an empty license list
	component3 := cyclonedx.Component{
		Licenses: nil,
	}
	result3 := extractLicense(component3)
	if result3 != "" {
		t.Errorf("Expected empty license, got '%s'", result3)
	}

}

func TestKissBOM_JSON(t *testing.T) {
	kissBOM := KissBOM{
		Packages: []Package{
			{Purl: "pkg:pypi/requests@2.26.0", License: "MIT", Copyright: "Copyright 2023", Notes: "Some notes"},
			{Purl: "pkg:pypi/requests@2.26.1", License: "Apache 2.0", Copyright: "Copyright 2023", Notes: ""},
		},
	}

	jsonData, err := kissBOM.JSON()
	assert.NoError(t, err)

	// Unmarshal JSON to verify its correctness
	var unmarshalledBOM KissBOM
	err = json.Unmarshal(jsonData, &unmarshalledBOM)
	assert.NoError(t, err)

	// Check if data is correctly preserved
	assert.Equal(t, kissBOM, unmarshalledBOM)
}

func TestKissBOM_CSV(t *testing.T) {
	kissBOM := KissBOM{
		Packages: []Package{
			{Purl: "pkg:pypi/requests@2.26.0", License: "MIT", Copyright: "Copyright 2023", Notes: "Some notes"},
			{Purl: "pkg:pypi/requests@2.26.1", License: "Apache 2.0", Copyright: "Copyright 2023", Notes: ""},
		},
	}

	csvData, err := kissBOM.CSV()
	assert.NoError(t, err)

	// Add assertions based on the expected CSV format
	// For simplicity, we'll just check if certain keywords are present in the CSV data
	assert.Contains(t, string(csvData), "purl,license,copyright,notes")
	assert.Contains(t, string(csvData), "pkg:pypi/requests@2.26.0,MIT,Copyright 2023,Some notes")
	assert.Contains(t, string(csvData), "pkg:pypi/requests@2.26.1,Apache 2.0,Copyright 2023,")

	// You may also consider using a CSV parsing library to validate the CSV structure more precisely
}

func TestKissBOM_YAML(t *testing.T) {
	kissBOM := KissBOM{
		Packages: []Package{
			{Purl: "pkg:pypi/requests@2.26.0", License: "MIT", Copyright: "Copyright 2023", Notes: "Some notes"},
			{Purl: "pkg:pypi/requests@2.26.1", License: "Apache 2.0", Copyright: "Copyright 2023", Notes: ""},
		},
	}

	yamlData, err := kissBOM.YAML()
	assert.NoError(t, err)

	// Unmarshal YAML to verify its correctness
	var unmarshalledBOM KissBOM
	err = yaml.Unmarshal(yamlData, &unmarshalledBOM)
	assert.NoError(t, err)

	// Check if data is correctly preserved
	assert.Equal(t, kissBOM, unmarshalledBOM)
}

func TestKissBOM_Minimal(t *testing.T) {
	kissBOM := KissBOM{
		Packages: []Package{
			{Purl: "pkg:pypi/requests@2.26.0", License: "MIT", Copyright: "Copyright 2023", Notes: "Some notes"},
			{Purl: "pkg:pypi/requests@2.26.1", License: "Apache 2.0", Copyright: "Copyright 2023", Notes: ""},
		},
	}

	minimalData, err := kissBOM.Minimal()
	assert.NoError(t, err)

	// Unmarshal JSON to verify its correctness
	var unmarshalledBOM KissBOM
	err = json.Unmarshal(minimalData, &unmarshalledBOM)
	assert.NoError(t, err)

	// Check if data is correctly preserved (only Purl should be present)
	assert.Len(t, unmarshalledBOM.Packages, 2)
	assert.Equal(t, "pkg:pypi/requests@2.26.0", unmarshalledBOM.Packages[0].Purl)
	assert.Empty(t, unmarshalledBOM.Packages[0].License)
	assert.Empty(t, unmarshalledBOM.Packages[0].Copyright)
	assert.Empty(t, unmarshalledBOM.Packages[0].Notes)
	assert.Equal(t, "pkg:pypi/requests@2.26.1", unmarshalledBOM.Packages[1].Purl)
	assert.Empty(t, unmarshalledBOM.Packages[1].License)
	assert.Empty(t, unmarshalledBOM.Packages[1].Copyright)
	assert.Empty(t, unmarshalledBOM.Packages[1].Notes)
}

func TestCompatible_Success(t *testing.T) {
	// Create a KissBOM with some packages for testing
	kissBOM := &KissBOM{
		Packages: []Package{
			{Purl: "pkg:pypi/requests@2.26.0"},
			{Purl: "pkg:nuget/Newtonsoft.Json@12.0.3"},
		},
	}

	// Call the Compatible method
	result, err := kissBOM.Compatible()

	// Assertions
	assert.NoError(t, err, "Expected no error")
	assert.NotNil(t, result, "Expected result to be not nil")
	assert.Contains(t, string(result), "cyclonedx.org/schema")
}

func TestEncodeBOM_Success(t *testing.T) {
	// Create a KissBOM with some packages for testing
	kissBOM := &KissBOM{
		Packages: []Package{
			{Purl: "pkg:pypi/requests@2.26.0"},
			{Purl: "pkg:nuget/Newtonsoft.Json@12.0.3"},
		},
	}

	bom := cyclonedx.NewBOM()
	components := []cyclonedx.Component{
		{
			PackageURL: "pkg:pypi/requests@2.26.0",
		},
		{
			PackageURL: "pkg:nuget/Newtonsoft.Json@12.0.3",
		},
	}
	bom.Components = &components

	// Call the encodeBOM method
	result, err := kissBOM.encodeBOM(bom)

	// Assertions
	assert.NoError(t, err, "Expected no error")
	assert.NotNil(t, result, "Expected result to be not nil")
	// Add additional assertions as needed
}
