# pglogalyze

A fast CLI tool for searching and filtering PostgreSQL log files — with binary search for time-range queries on large files.

---

## Features

- **Time-range filtering** — find logs between a `startTime` and `endTime`
- **Severity filtering** — filter by log level (ERROR, WARNING, INFO, etc.)
- **Log type filtering** — filter by category (APPLI, QUERY, CONNECTION, DURATION, CHECKPOINT, STARTUP, SHUTDOWN)
- **Binary search** — dichotomous offset search to efficiently locate time ranges in large files without scanning from the start
- **Reverse block reading** — reads from the end of the file backwards, returning the most recent N matching lines
- **Multi-line log support** — correctly groups continuation lines with their parent log entry
- **Application log detection** — identifies `user@database` log lines

---

## Installation

```bash
git clone https://github.com/Antoine256/pglogalyze
cd pglogalyze
go build -o pglogalyze ./cmd/pglogalyze
```

Requires Go 1.21+.

---

## Usage

```
./pglogalyze -f <logfile> [options]
```

### Options

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-f` | string | *(required)* | Path to the PostgreSQL log file |
| `-n` | int | `20` | Number of log entries to return |
| `-l` | string | | Severity level filter (e.g. `ERROR`, `WARNING`, `INFO`) |
| `-t` | string | | Log type filter (`APPLI`, `QUERY`, `CONNECTION`, `DURATION`, `CHECKPOINT`, `STARTUP`, `SHUTDOWN`) |
| `-st` | string | | Start time filter, format: `YYYY-MM-DDTHH:MM:SS` |
| `-et` | string | | End time filter, format: `YYYY-MM-DDTHH:MM:SS` |

### Examples

Get the last 20 log entries:
```bash
./pglogalyze -f /var/log/postgresql/postgresql.log
```

Get the last 50 ERROR logs:
```bash
./pglogalyze -f postgresql.log -n 50 -l ERROR
```

Get logs between two timestamps:
```bash
./pglogalyze -f postgresql.log -st 2024-01-15T10:00:00 -et 2024-01-15T11:00:00
```

Get the last 10 query logs after a given time:
```bash
./pglogalyze -f postgresql.log -st 2024-01-15T09:00:00 -t QUERY -n 10
```

---

## How it works

### Without time filters
The file is read **backwards** in 4KB blocks using `ReadAt`. Multi-line log entries (continuation lines without a timestamp prefix) are accumulated into a buffer and grouped with their parent entry. Reading stops once N entries are collected.

### With time filters
1. The **first and last lines** of the file are checked to quickly discard files with no matching entries.
2. A **binary search** (dichotomous offset search) locates the position in the file closest to `endTime`, by finding two consecutive timestamped lines that bracket the target time.
3. The parser then reads **backwards from that offset**, collecting entries until it has N lines or reaches `startTime`.

This avoids scanning the entire file for time-range queries, making it efficient even on multi-GB log files.

---

## Log format

pglogalyze expects standard PostgreSQL log lines in the format:

```
YYYY-MM-DD HH:MM:SS.mmm TZ [PID] LEVEL: message
YYYY-MM-DD HH:MM:SS.mmm TZ [PID] user@database LEVEL: message
```

---