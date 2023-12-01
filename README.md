![kissbom](img/kissbom128x128.png)

# kissbom

[![](https://img.shields.io/badge/Status-ALPHA-orange)](CONTRIBUTING.md)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/devops-kung-fu/kissbom) 
[![Go Report Card](https://goreportcard.com/badge/github.com/devops-kung-fu/kissbom)](https://goreportcard.com/report/github.com/devops-kung-fu/kissbom) 

Converts a CycloneDX file into a KissBOM. Implements the [kissbom-spec](https://github.com/kissbom/kissbom-spec). 

## Overview

Let's face it - SBOMs can be massive in size. On an episode of [daBOM](https://dabom.show/philippe-ombredanne/), [Philippe Ombredanne](https://github.com/pombredanne) mentioned that he had created a very minimal SBOM specification called KissBOM. KissBOMs are pretty much the bare minimum that one would need to describe software in an SBOM format. 

We thought it may be cool to implement a CLI that would convert a CycloneDX file to a KissBOM. ```kissbom``` will take a CycloneDX file, remove all non-essential fields, and lets you save it to a variety of formats - including a CycloneDX formatted kissbom.

Using a sample CycloneDX SBOM for [juiceshop](./_TESTDATA_/juiceshop.cyclonedx.json), we found that a generated kissbom in JSON format was **1/10th of the size** of the original file.

## KissBOM Content

KissBOMs contain a similar collection of packages that are defined in a CycloneDX format, but only the essential fields. The fields that are maintained from the CycloneDX spec are:

| Field | Description | Required |
|---|---|---|
| PURL | The package url | YES |
| License | The defined License of the package | NO |
| Copyright | The copyright for the package | NO |
| Notes | Any notes available for the package | NO |

## Installation

### Mac

You can use [Homebrew](https://brew.sh) to install ```kissbom``` using the following:

``` bash
brew tap devops-kung-fu/homebrew-tap
brew install devops-kung-fu/homebrew-tap/kissbom
```

If you do not have Homebrew, you can still [download the latest release](https://github.com/devops-kung-fu/kissbom/releases) (ex: ```kissbom.1.0_darwin_all.tar.gz```), extract the files from the archive, and use the ```kissbom``` binary.  

If you wish, you can move the ```kissbom``` binary to your ```/usr/local/bin``` directory or anywhere on your path.

### Linux

To install ```kissbom```,  [download the latest release](https://github.com/devops-kung-fu/kisbbom/releases) for your platform and install locally. For example, install ```kissbom``` on Ubuntu:

```bash
dpkg -i kissbom_0.4.1_linux_arm64.deb
```
## Usage

```kissbom``` is a really simple CLI with only a small number of options. To quickly convert a CycloneDX SBOM to a JSON formatted KissBOM, run the following:

``` bash
kissbom convert test.cyclonedx.json //where test.cyclonedx.json is a valid CycloneDX SBOM
```

### Output Formats

```kissbom``` can output a KissBOM in a variety of formats using the ```--format``` flag. Valid options are:

| Option | Description |
|---|---|
|```--format=json``` | Outputs all 4 KissBOM fields in JSON format. This is the default output format |
|```--format=yaml``` | Outputs all 4 KissBOM fields in YAML format |
|```--format=csv``` | Outputs all 4 KissBOM fields into a CSV formatted file |
|```--format=minimal``` | Outputs just the KissBOM required fields into a JSON formatted file (Purl) |
|```--format=compatible``` | Outputs all 4 KissBOM fields in a CycloneDX formatted JSON file |

### Debugging

To enable verbose logging in ```kissbom```, use the ```--debug``` flag.

## Credits

A big thank-you to our friends at [Good Ware](https://www.flaticon.com/authors/good-ware) for the ```kissbom``` logo.