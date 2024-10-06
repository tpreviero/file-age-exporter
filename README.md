
# File Age Exporter

The **File Age Exporter** is a Prometheus exporter that tracks how long files in specified directories have remained unchanged (since last modification). 
It aggregates file modification times and exposes them in a metrics-friendly format that Prometheus can scrape.

The metrics are presented in the following format:

```
file_since_total{year="YYYY", month="MMMM", week="xxx"} 123456
```

- The exporter generates a series for each **week**, **month**, and **year**.
- Files modified within the same week of a given month and year are aggregated under the same series.

## Building the Exporter

To build the exporter using Docker, use the following command:

```bash
docker run \
  --rm \
  -v $(pwd):/app \
  -w /app \
  golang:1.23 \
  go build -o file_age_exporter
```

## Running the Exporter

You can run the exporter using Docker with this command:

```bash
docker run \
  --rm \
  -it \
  -v $(pwd):/go/src/file-age-exporter \
  -v $(pwd):/a-directory \
  -w /go/src/file-age-exporter \
  -p 9123:9123 \
  golang:1.23 \
  go run . \
  --dir /a-directory
```

In this example, the exporter will collect file modification times from `/a-directory` and expose metrics on port `9123`.

## Usage

The following CLI options are available:

- `--dir`
  - Specifies the directories to scan for file metrics.
  - You can use this flag multiple times for different directories.
  - Note: Overlapping directories will be counted only once.

- `--exclude`
  - Exclude files or directories based on a glob pattern (e.g., `*.log`).
  - You can use this flag multiple times to exclude multiple patterns.

- `--listen-address` (default `:9123`)
  - The address on which the exporter listens for HTTP requests.

- `--walk-interval` (default `60` seconds)
  - The interval (in seconds) at which the exporter will scan the directories to update the metrics.

- `--help`
  - Displays help information about the available options.
