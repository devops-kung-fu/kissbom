package lib

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/devops-kung-fu/kissbom/models"
)

func TestConvert_Success(t *testing.T) {
	jsonContent := `
	{
		"bomFormat": "CycloneDX",
		"specVersion": "1.3",
		"components": [
			{
				"type": "library",
				"group": "",
				"name": "",
				"version": "",
				"purl": "pkg:pypi/requests@2.26.0",
				"licenses": [
					{
						"expression": "(AFL-2.1 OR BSD-3-Clause)"
					}
				]
			}
		]
	}
	`
	converter := Converter{
		Afs: &afero.Afero{Fs: afero.NewMemMapFs()},
	}

	e := func() error {
		var data []byte = []byte(jsonContent)
		f, err := (converter.Afs).Fs.OpenFile("test.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0644))
		if err != nil {
			return err
		}
		n, err := f.Write(data)
		if err == nil && n < len(data) {
			err = io.ErrShortWrite
		}
		if err1 := f.Close(); err == nil {
			err = err1
		}
		return err
	}()

	assert.NoError(t, e)

	converter.OutputFormat = "json" // Choose a valid output format for testing
	converter.OutputFileName = "test_output"
	err := converter.Convert("test.json")
	assert.NoError(t, err, "Expected no error")

	converter.OutputFormat = "yaml" // Choose a valid output format for testing
	err = converter.Convert("test.json")
	assert.NoError(t, err, "Expected no error")

	converter.OutputFormat = "csv" // Choose a valid output format for testing
	err = converter.Convert("test.json")
	assert.NoError(t, err, "Expected no error")

	converter.OutputFormat = "minimal" // Choose a valid output format for testing
	err = converter.Convert("test.json")
	assert.NoError(t, err, "Expected no error")

	converter.OutputFormat = "compatible" // Choose a valid output format for testing
	err = converter.Convert("test.json")
	assert.NoError(t, err, "Expected no error")

	converter.OutputFormat = "barf" // Choose a valid output format for testing
	err = converter.Convert("test.json")
	assert.Error(t, err, "Expected no error")

	e = func() error {
		var data []byte = []byte("<>test")
		f, err := (*converter.Afs).Fs.OpenFile("test.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0644))
		if err != nil {
			return err
		}
		n, err := f.Write(data)
		if err == nil && n < len(data) {
			err = io.ErrShortWrite
		}
		if err1 := f.Close(); err == nil {
			err = err1
		}
		return err
	}()

	assert.NoError(t,e)
	converter.OutputFormat = "csv" // Choose a valid output format for testing
	err = converter.Convert("test.json")
	assert.Error(t, err, "Expected no error")

}

func TestConvert_FileReadError(t *testing.T) {
	converter := Converter{
		Afs: &afero.Afero{Fs: afero.NewMemMapFs()},
	}
	err := converter.Convert("nonexistent_file.json")

	assert.Error(t, err, "Expected an error due to nonexistent file")
}

func TestTransform_Success(t *testing.T) {
	jsonContent := `
	{
		"bomFormat": "CycloneDX",
		"specVersion": "1.3",
		"components": [
			{
				"type": "library",
				"group": "",
				"name": "",
				"version": "",
				"purl": "pkg:pypi/requests@2.26.0"
			}
		]
	}`

	converter := Converter{
		Afs:            &afero.Afero{Fs: afero.NewMemMapFs()},
		OutputFileName: "test",
	}

	kissBom, err := converter.transform([]byte(jsonContent))

	assert.NoError(t, err, "Expected no error")
	assert.NotNil(t, kissBom, "Expected KissBOM object to be not nil")
	assert.NotEmpty(t, converter.OutputFileName, "Expected filename to be not empty")
	assert.Equal(t, converter.OutputFileName, "test")
	assert.Len(t, kissBom.Packages, 1)
	assert.Equal(t, kissBom.Packages[0].Purl, "pkg:pypi/requests@2.26.0")
}

func TestTransform_DecodeError(t *testing.T) {
	converter := Converter{
		Afs: &afero.Afero{Fs: afero.NewMemMapFs()},
	}

	invalidCycloneDxJSON := []byte(`{"fake""}`)

	_, err := converter.transform(invalidCycloneDxJSON)

	assert.Error(t, err, "Expected an error due to invalid JSON")
}

func TestBuildOutputFilename(t *testing.T) {
	converter := NewConverter()
	currentTime := time.Now().UTC()

	// Format the time as a string using a specific layout
	timeString := currentTime.Format("2006-01-02T15:04:05Z")

	// Create a CycloneDX BOM with test metadata
	testBOM := &cyclonedx.BOM{
		Metadata: &cyclonedx.Metadata{
			Component: &cyclonedx.Component{
				Name:      "TestComponent",
				Version:   "1.0.0",
				Publisher: "TestPublisher",
			},
			Timestamp: timeString,
		},
	}

	outputFilename := converter.buildOutputFilename(testBOM)

	assert.NotEmpty(t, outputFilename, "Expected output filename to be not empty")
}

func TestConverter_writeToFile(t *testing.T) {
	converter := Converter{
		Afs:            &afero.Afero{Fs: afero.NewMemMapFs()},
		OutputFileName: "&*%/\"",
		OutputFormat:   models.OptionMinimal,
	}

	kissBOM := models.KissBOM{
		Packages: []models.Package{
			{Purl: "pkg:pypi/requests@2.26.0", License: "MIT", Copyright: "Copyright 2023", Notes: "Some notes"},
		},
	}

	err := converter.writeToFile(kissBOM)
	assert.NoError(t, err)

}
