# pmm-dump-qan-sql

Convert [pmm-dump](https://github.com/percona/pmm-dump) QAN (Query Analytics) ClickHouse chunks (`ch/*.tsv`) into `INSERT` SQL for manual load via `clickhouse-client`.

License: [GPL-2.0](LICENSE)

## Prerequisites

- Go 1.26+ (to build)
- A QAN chunk from `pmm-dump export --dump-qan` (files under `ch/` in the dump, e.g. `0.tsv`)
- For loading SQL: `clickhouse-client` on the PMM server host or inside the `pmm-server` container

## Install

Clone this repo:

```bash
git clone <this-repo> pmm-dump-qan-sql
cd pmm-dump-qan-sql
go build -o pmm-dump-qan-sql ./cmd/pmm-dump-qan-sql
```

Install to `$GOBIN` or `/usr/local/bin`:

```bash
go install ./cmd/pmm-dump-qan-sql
```

## Extract a TSV chunk from a pmm-dump archive

If you only have the `.tar.gz` dump:

```bash
tar -xzf dump.tar.gz ch/0.tsv
# or extract the whole ch/ directory
tar -xzf dump.tar.gz ch/
```

## Usage

Converts pmm-dump’s custom tab-separated format (not native ClickHouse `TabSeparated`) into batched `INSERT INTO pmm.metrics VALUES ...` statements.

### PMM version detection

The tool auto-detects PMM major version from the first row width:

| Columns | PMM version |
|---------|-------------|
| 228 | PMM 2 |
| 269 | PMM 3 |

Override with `--pmm-version=pmm2` or `--pmm-version=pmm3` if needed.

### Examples

**PMM 3 chunk → SQL (auto-detect):**

```bash
pmm-dump-qan-sql ../source/ch/0.tsv -o ../source/ch/0.sql
```

**PMM 2 chunk:**

```bash
pmm-dump-qan-sql ../source/pmm2/ch/0.tsv -o ../source/pmm2/ch/0.sql
```

**Explicit output path and smaller batches:**

```bash
pmm-dump-qan-sql -o /tmp/qan.sql --batch-size=50 ../source/ch/0.tsv
```

**Verbose logging:**

```bash
pmm-dump-qan-sql -v ../source/ch/0.tsv
```

On success, the tool prints the `clickhouse-client` command for your PMM version, for example:

```text
Detected pmm3 QAN dump (269 columns, 500 rows).
Load into PMM ClickHouse from the PMM server host/container:
clickhouse-client --database=pmm --password=clickhouse --queries-file "../source/ch/0.sql"
```

### Load SQL into PMM ClickHouse

Run the printed command **on the PMM server** (host or container), after copying the `.sql` file if necessary.

**PMM 3** (default ClickHouse password):

```bash
docker cp ../source/ch/0.sql pmm-server:/tmp/0.sql
docker exec -it pmm-server clickhouse-client \
  --database=pmm --password=clickhouse --queries-file /tmp/0.sql
```

**PMM 2** (no password on native protocol by default):

```bash
docker cp ../source/pmm2/ch/0.sql pmm-server:/tmp/0.sql
docker exec -it pmm-server clickhouse-client \
  --database=pmm --queries-file /tmp/0.sql
```

### Flags

| Flag | Description |
|------|-------------|
| `-o`, `--output` | Output `.sql` path (default: same basename as TSV, `.sql`) |
| `--pmm-version` | `auto` (default), `pmm2`, or `pmm3` |
| `--batch-size` | Rows per `INSERT` statement (default: 100) |
| `--database` | Database in SQL (default: `pmm`) |
| `--table` | Table in SQL (default: `metrics`) |
| `-v`, `--verbose` | Debug logging |
