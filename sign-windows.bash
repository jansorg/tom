#!/usr/bin/env bash
#
# Signs the Windows binaries of a draft GitHub release with the code signing
# certificate on a YubiKey.
#
# Usage: ./sign-windows.bash <tag> [--dry-run]
#   e.g. ./sign-windows.bash v0.14.0
#
# Workflow:
#   1. Create the draft release: gh workflow run release.yml -f version=0.14.0
#   2. Run this script with the YubiKey plugged in
#   3. Review and publish the draft release on GitHub
#
# The script downloads the Windows archive of the draft release, signs tom.exe
# and tom-pdf.exe (Authenticode, timestamped), repacks the archive, updates
# checksums.txt and uploads both files to the draft release.
# With --dry-run, nothing is uploaded and the files are kept for inspection.
#
# Requirements: gh, jsign, java, yubico-piv-tool (provides libykcs11), osslsigncode (optional, for verification)
#
# Environment variables (all optional):
#   JSIGN_PIN        PIN of the YubiKey, prompted for if not set
#   JSIGN_ALIAS      alias of the certificate, e.g. "X.509 Certificate for Digital Signature" (PIV slot 9c)
#   JSIGN_TSAURL     timestamping server, default: http://timestamp.sectigo.com
#   JSIGN_STORETYPE  keystore type, default: YUBIKEY
#   JSIGN_KEYSTORE   keystore file, only needed for file-based keystores (e.g. for testing with PKCS12)

set -euo pipefail

TAG="${1:-}"
DRY_RUN="${2:-}"
if [[ -z "$TAG" || ( -n "$DRY_RUN" && "$DRY_RUN" != "--dry-run" ) ]]; then
    echo "Usage: $0 <tag> [--dry-run]" >&2
    exit 1
fi

for cmd in gh jsign tar sha256sum; do
    command -v "$cmd" >/dev/null || { echo "Required command not found: $cmd" >&2; exit 1; }
done

# Never modify a published release, its files may already have been downloaded
IS_DRAFT=$(gh release view "$TAG" --json isDraft --jq .isDraft)
if [[ "$IS_DRAFT" != "true" ]]; then
    echo "Release $TAG is already published, only draft releases can be signed." >&2
    exit 1
fi

ARCHIVE=$(gh release view "$TAG" --json assets --jq '.assets[].name | select(endswith("_Windows_amd64.tar.gz"))')
if [[ -z "$ARCHIVE" ]]; then
    echo "No Windows archive found in release $TAG." >&2
    exit 1
fi

WORK_DIR=$(mktemp -d)
if [[ -z "$DRY_RUN" ]]; then
    trap 'rm -rf "$WORK_DIR"' EXIT
fi

echo "Downloading $ARCHIVE and checksums.txt of release $TAG..."
gh release download "$TAG" --dir "$WORK_DIR" --pattern "$ARCHIVE" --pattern checksums.txt

# verify the download before modifying anything
(cd "$WORK_DIR" && grep " $ARCHIVE\$" checksums.txt | sha256sum --check --quiet)

mkdir "$WORK_DIR/content"
tar -xzf "$WORK_DIR/$ARCHIVE" -C "$WORK_DIR/content"
mapfile -t FILES < <(tar -tzf "$WORK_DIR/$ARCHIVE")

if [[ -z "${JSIGN_PIN:-}" ]]; then
    read -r -s -p "YubiKey PIN: " JSIGN_PIN
    echo
fi
export JSIGN_PIN

JSIGN_ARGS=(
    --storetype "${JSIGN_STORETYPE:-YUBIKEY}"
    --storepass env:JSIGN_PIN
    --tsaurl "${JSIGN_TSAURL:-http://timestamp.sectigo.com}"
)
[[ -n "${JSIGN_KEYSTORE:-}" ]] && JSIGN_ARGS+=(--keystore "$JSIGN_KEYSTORE")
[[ -n "${JSIGN_ALIAS:-}" ]] && JSIGN_ARGS+=(--alias "$JSIGN_ALIAS")

for exe in tom.exe tom-pdf.exe; do
    echo "Signing $exe..."
    jsign "${JSIGN_ARGS[@]}" "$WORK_DIR/content/$exe"

    if command -v osslsigncode >/dev/null; then
        # only checks that a valid signature is present, the certificate chain is checked by Windows
        osslsigncode verify -in "$WORK_DIR/content/$exe" 2>&1 | grep -E "^(Signer's certificate|\s+Subject:|Timestamp)" | head -3 || true
    fi
done

echo "Repacking $ARCHIVE..."
rm "$WORK_DIR/$ARCHIVE"
tar -czf "$WORK_DIR/$ARCHIVE" --owner=0 --group=0 --numeric-owner -C "$WORK_DIR/content" "${FILES[@]}"

CHECKSUM=$(cd "$WORK_DIR" && sha256sum "$ARCHIVE")
sed -i "s|^[0-9a-f]\{64\}  $ARCHIVE\$|$CHECKSUM|" "$WORK_DIR/checksums.txt"
(cd "$WORK_DIR" && grep " $ARCHIVE\$" checksums.txt | sha256sum --check --quiet)

if [[ -n "$DRY_RUN" ]]; then
    echo "Dry run, nothing uploaded. Signed files are in $WORK_DIR"
    exit 0
fi

echo "Uploading signed archive and checksums to release $TAG..."
gh release upload "$TAG" --clobber "$WORK_DIR/$ARCHIVE" "$WORK_DIR/checksums.txt"
echo "Done. Review and publish the draft release: $(gh release view "$TAG" --json url --jq .url)"
