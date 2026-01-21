# pglogalyze

A simple command-line tool for analyzing PostgreSQL log files.

## Overview

pglogalyze is a lightweight utility that parses and analyzes PostgreSQL log files, providing insights into database activity, query performance, and potential issues within a specified time range.

## Features

- Parse PostgreSQL log files
- Filter logs by time range and severity
- Simple command-line interface

## Installation

### Prerequisites

- Go 1.16 or higher

## Usage

### Basic Command

```bash
go run ./cmd/pglogalyze/main.go -f
```

### Example

```bash
go run ./cmd/pglogalyze/main.go \
  -f ../../../Desktop/psql/log/postgresql-2026-01-18_124734.log \
  -l INFO \
  -st "2026-01-18 12:51:00" \
  -et "2026-01-18 14:30:34"
```

### Command-line Options

| Option | Description | Required | Format |
|--------|-------------|----------|--------|
| `-f` | Path to PostgreSQL log file | Yes | File path |
| `-l` | Level of severity to analyze | No | `INFO - LOG - ...`
| `-st` | Start time for analysis | No | `YYYY-MM-DD HH:MM:SS` |
| `-et` | End time for analysis | No | `YYYY-MM-DD HH:MM:SS` |
