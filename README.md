# gank

A local, privacy-focused CLI tool to convert bank statement PDFs into [hledger](https://hledger.org/) journal format. No data leaves your machine. Single binary, zero runtime dependencies.

## Requirements

- Go 1.21+ (for building)
- No runtime dependencies

## Build

```bash
make build
```

Or directly:
```bash
go build -o extractor ./cmd/extractor
```

### Cross-Compilation

```bash
make cross-build
```

This produces binaries for:
- `extractor-darwin-amd64` (macOS Intel)
- `extractor-darwin-arm64` (macOS Apple Silicon)
- `extractor-linux-amd64` (Linux)
- `extractor-windows-amd64.exe` (Windows)

## Usage

### 1. Inspect the PDF

Dump the raw extracted content to see exactly what the PDF reader extracts from your statement:

```bash
./extractor inspect statement.pdf
```

Use this output to identify the transaction line format, then tune the matching bank config in `banks/`.

### 2. Configure your bank

Bank configs live in `banks/<name>/` with one YAML file per account type. Edit the one matching your bank and account (e.g. `banks/klar/credit.yaml`):

| Field | Description |
|---|---|
| `transaction_pattern` | Regex with named groups for date, description, and amount |
| `group_date` / `group_description` / `group_amount` | Names of those regex groups (default: `date`, `description`, `amount`) |
| `date_input_format` | Go time format of dates in the PDF (e.g. `02 01 2006`) |
| `date_output_format` | Output date format for hledger (default: `2006-01-02`) |
| `month_map` | Map of month names → numbers for non-English statements (e.g. `enero: "01"`) |
| `debit_is_positive` | `true` if the statement shows expenses as positive amounts (e.g. credit cards) |
| `account_assets` | hledger account for this bank account (use `liabilities:...` for credit cards) |
| `account_expenses` | Default account for outgoing transactions |
| `account_income` | Default account for incoming transactions |
| `currency` | Currency symbol (e.g. `$`, `€`) |
| `exclude_patterns` | List of regex patterns; transactions matching any are skipped |
| `section_start` / `section_end` | Optional regex to scope parsing to a page section |

### Adding a new bank

Create `banks/<name>/<account>.yaml`:

```yaml
transaction_pattern: '(?P<date>\d{2}/\d{2}/\d{4})\s+(?P<description>.+?)\s+(?P<amount>\d+\.\d{2})'
date_input_format: "02/01/2006"
debit_is_positive: false
account_assets: assets:bank:mybank
currency: "$"
```

For banks with non-English month names:

```yaml
month_map:
  enero: "01"
  febrero: "02"
date_input_format: "02 01 2006"
```

### 3. Convert

```bash
# Print to stdout
./extractor convert --bank klar --account credit statement.pdf

# Write to file
./extractor convert --bank klar --account credit statement.pdf -o transactions.journal
```

### 4. Import into hledger

```bash
hledger -f transactions.journal stats
hledger -f transactions.journal register
```

Or include the journal in your main ledger file:

```
include transactions.journal
```

## Output format

**Credit card** (`debit_is_positive: true`):
```
2026-01-24 * Taxi
    expenses:unknown          $100.60
    liabilities:credit:klar   $-100.60
```

**Checking / savings** (`debit_is_positive: false`):
```
2026-01-24 * Salary
    assets:bank:klar:checking   $5000.00
    income:unknown              $-5000.00
```

After importing, reclassify `expenses:unknown` and `income:unknown` into proper categories using hledger's account tagging.

## Project structure

```
gank/
├── cmd/extractor/main.go          # CLI entry point
├── internal/
│   ├── config/
│   │   ├── config.go             # BankConfig struct
│   │   └── loader.go             # YAML config loader
│   └── extractor/
│       ├── extractor.go          # PDF text extraction
│       ├── parser.go             # Transaction parsing
│       └── formatter.go          # hledger output formatting
├── banks/
│   └── klar/
│       ├── credit.yaml
│       ├── checking.yaml
│       └── investing.yaml
├── go.mod
├── go.sum
└── Makefile
```

## Privacy

All processing is done locally. The PDF library reads files from disk only and makes no network requests.

## Testing

```bash
make test
```

## License

MIT