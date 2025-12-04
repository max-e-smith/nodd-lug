# cruise-lug

CLI tool for downloading marine geophysics survey data from the NOAA Open Data Dissemenation (NODD) hosted cloud.

## installation

### Install module with Go

#### Requirements
Install Go: https://go.dev/doc/install

#### Instructions
TODO

### Install for your platform
TODO -- https://goreleaser.com/customization/builds/go/

## usage
clug [command] [options] [survey] [local_path]

Commands
- get
- config
- update

Options
`-b | --bathy`: downloads bathymetry data from specified survey
`-s | --summary`: provides a summary of the survey request instead of download
