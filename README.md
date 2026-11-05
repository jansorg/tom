
[![GoDoc](https://godoc.org/github.com/jansorg/tom?status.svg)](https://godoc.org/github.com/jansorg/tom)
[![Release](https://github.com/jansorg/tom/actions/workflows/release.yml/badge.svg)](https://github.com/jansorg/tom/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/jansorg/tom)](https://goreportcard.com/report/github.com/jansorg/tom)

**This project is still in an early phase. It's usable and working but not yet recommended for production.**

**A simple UI is developed at https://github.com/jansorg/tom-ui (Qt5 tested on Linux and macOS). A preview for mac is available there.**

# tom

tom is a command line application to track time and to simplify office work.
It's able to track time and to create HTML and PDF reports.

`tom` tracks time locally. 
Instead of implementing billing, etc. it will rely on 3rd-party tools, i.e. mostly cloud-based services.

## Get tom
Prebuilt binaries are available on the [releases page](https://github.com/jansorg/tom/releases):

- Linux (amd64): `.tar.gz` archive or `.deb` package, which also installs the bash completion
- macOS (Apple Silicon): `.tar.gz` archive, signed and notarized
- Windows (amd64): `.tar.gz` archive, signed

Each archive contains `tom` and `tom-pdf`.

With Go 1.27 or later installed, you can also install tom with:

```bash
go install github.com/jansorg/tom@latest
```

After that `tom` will be available, assuming that `$(go env GOPATH)/bin` is in your `$PATH`.
See [Building](#building) to build tom from source.

## Basic usage
A typical session looks like this:

```bash
tom create project acme
tom start acme
sleep 10
tom stop
tom report --month 0
```

## Documentation
Documentation about the command line is available at [docs/markdown](./docs/markdown/tom.md)

## Data model
The data is stored in a few JSON files on disk. It's easy to backup and still fast.
Right now it's on disk at `$HOME/.tom/`.
The available commands make it easy to export it into different target formats. Most commands are supporting plain text
and JSON at this time. Other formats may be added in the future.

### Projects
tom supports nested projects. The separator character is the slash '/'.
The simplest form is a project without any subprojects.

If you need to bill a single client for several distinct project then it's better to use subprojects.
For example:
```bash
tom create project client1 client1/web client1/backend
tom start client1/web
tom stop
tom report --split project -p client1
```

This will create a report on all projects which belong to client1 with the tracked time per project.

## Tracking time
### Start
### Stop
### Cancel
### View

## Reporting

### Date and time filters
### Splitting options

The reporting allows to split the tracked time into multiple pieces. Multiple levels of splitting is possible.
Possible values:
- year
- month
- project

For example, `-split year` creates a report where the tracked time is grouped into the years which are covered by tracked time.
`-split year,project` will renders groups of years which are further grouped by the projects which were tracked in this particular year. 
`-split project,month` groups the tracked time by project first and then lists the summaries grouped by month for each project. 

### Plain Text Reports

### HTML Reports

### PDF Reports
#### api2pdf service integration

## Import from other tools

### Import data from Watson
```bash
tom import watson
```

### Import data from Fanurio
This needs a custom CSV export (to be documented).
```bash
tom import fanurio complete-export.csv
```

### Import data from Mac Time Tracker
This imports the CSV format exported by Mac Time Tracker.
```bash
tom import macTimeTracker export.csv
```


## Development

### Building
Requirements: Go 1.27 or later. No other tools are needed, the report templates in [`templates/`](./templates) are embedded into the binary with `go:embed`.

```bash
git clone https://github.com/jansorg/tom.git
cd tom
go build .                        # builds ./tom
go build -o tom-pdf ./go-tom-pdf  # builds ./tom-pdf, the PDF converter
go test ./...
```

To build the release artifacts locally without publishing anything, install [GoReleaser](https://goreleaser.com) v2 and run:

```bash
goreleaser check
goreleaser release --snapshot --clean --skip=publish
```

The archives and the `.deb` package are written to `dist/`.

Generated files, updated automatically for each release:
- `scripts/completions/tom.sh`: bash completion, created by `tom completion`
- `docs/markdown` and `docs/man`: command line documentation, created by [`update_docs.bash`](./update_docs.bash)

### Releasing
Releases are built by the GitHub Actions workflow [`release.yml`](./.github/workflows/release.yml) with GoReleaser.
Every release is created as a **draft** and must be published manually.

1. Start the release workflow on `master`, either from the Actions tab or with:
   ```bash
   gh workflow run release.yml -f version=0.14.0
   ```
   The workflow runs the tests, updates the version in `main.go`, the bash completion and the documentation (as `prep:` commits on `master`),
   creates and pushes the tag `v0.14.0` and builds the draft release.
   The macOS binaries are signed and notarized during the build.
2. Sign the Windows binaries of the draft release with the code signing certificate on the YubiKey:
   ```bash
   ./sign-windows.bash v0.14.0
   ```
   The script asks for the YubiKey PIN, signs `tom.exe` and `tom-pdf.exe`, and replaces the Windows archive and `checksums.txt` of the draft release.
   Use `--dry-run` to sign without uploading. It needs `gh`, `jsign`, `java` and `yubico-piv-tool`.
   If jsign doesn't find the certificate, set `JSIGN_ALIAS`, e.g. `JSIGN_ALIAS="X.509 Certificate for Digital Signature"` for PIV slot 9c.
3. Review the draft on the [releases page](https://github.com/jansorg/tom/releases) and publish it.
   Published releases can't be signed anymore by `sign-windows.bash`.

Pushing a tag (e.g. with [`release.bash`](./release.bash)) also builds a draft release, but doesn't update the version in `main.go`.

macOS signing and notarization needs these repository secrets. Without them, the macOS binaries are not signed.

| Secret | Content |
|---|---|
| `MACOS_SIGN_P12` | "Developer ID Application" certificate with private key as `.p12`, base64-encoded (`base64 -w0 cert.p12`) |
| `MACOS_SIGN_PASSWORD` | Password of the `.p12` file |
| `MACOS_NOTARY_ISSUER_ID` | App Store Connect API issuer ID |
| `MACOS_NOTARY_KEY_ID` | App Store Connect API key ID |
| `MACOS_NOTARY_KEY` | App Store Connect API key (`AuthKey_XXXX.p8`), base64-encoded |

## License
To be decided
