# plex_exporter
[Plex](https://www.plex.tv) prometheus exporter

![Golang CI](https://github.com/sfragata/plex_exporter/workflows/Golang%20CI/badge.svg)

## Installation

### Mac

```
brew tap sfragata/tap

brew install sfragata/tap/plex_exporter
```

### Linux and Windows

get latest release [here](https://github.com/sfragata/plex_exporter/releases)

## Usage

```
plex_exporter - Prometheus exporter for plex

  Flags: 
        --version       Displays the program version string.
    -h  --help          Displays help with available flag, subcommand, and positional value parameters.
    -H  --host          Plex address
    -p  --port          Plex port (default: 32400)
    -t  --token         Plex token
    -pm --port-metrics  Plex exporter metrics port (default: 2112)
```    