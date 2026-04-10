#!/usr/bin/env bash
# bench2md.sh — runs all codec benchmarks and updates the Benchmarks section in README.md.
# Usage: ./internal/scripts/bench2md.sh
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
README="$ROOT/README.md"

# Run benchmarks, capture output
BENCH_OUTPUT=$(cd "$ROOT" && go test -tags=safe -bench=. -benchmem -run='^$' -count=1 work 2>&1)

# Parse benchmark lines into associative arrays
declare -A ENC_NS ENC_B ENC_A DEC_NS DEC_B DEC_A
CURRENT_PKG=""

while IFS= read -r line; do
  if [[ "$line" =~ ^pkg:\ (.+)$ ]]; then
    CURRENT_PKG="${BASH_REMATCH[1]}"
    CURRENT_PKG="${CURRENT_PKG#github.com/foomo/goencode/}"
    continue
  fi

  if [[ "$line" =~ ^BenchmarkCodec/(encode|decode)-[0-9]+[[:space:]]+ ]]; then
    direction="${BASH_REMATCH[1]}"
    read -r _ _ ns _ bytes _ allocs _ <<< "$line"
    if [[ "$direction" == "encode" ]]; then
      ENC_NS["$CURRENT_PKG"]="$ns"
      ENC_B["$CURRENT_PKG"]="$bytes"
      ENC_A["$CURRENT_PKG"]="$allocs"
    else
      DEC_NS["$CURRENT_PKG"]="$ns"
      DEC_B["$CURRENT_PKG"]="$bytes"
      DEC_A["$CURRENT_PKG"]="$allocs"
    fi
  fi
done <<< "$BENCH_OUTPUT"

# Sort codec names
CODECS=()
for key in "${!ENC_NS[@]}"; do
  CODECS+=("$key")
done
IFS=$'\n' CODECS=($(sort <<< "${CODECS[*]}")); unset IFS

# Write section content to temp file
SECTION_FILE=$(mktemp)
trap 'rm -f "$SECTION_FILE" "$README.tmp"' EXIT

cat > "$SECTION_FILE" <<EOF
## Benchmarks

> Measured with \`go test -bench=. -benchmem\` on $(go env GOARCH) ($(go env GOOS)).
> Results vary by hardware — use these as relative comparisons between codecs.

| Codec | Encode (ns/op) | Encode (B/op) | Encode (allocs/op) | Decode (ns/op) | Decode (B/op) | Decode (allocs/op) |
|-------|---------------:|---------------:|--------------------:|---------------:|---------------:|--------------------:|
EOF

for codec in "${CODECS[@]}"; do
  echo "| \`$codec\` | ${ENC_NS[$codec]} | ${ENC_B[$codec]} | ${ENC_A[$codec]} | ${DEC_NS[$codec]} | ${DEC_B[$codec]} | ${DEC_A[$codec]} |" >> "$SECTION_FILE"
done

# Verify markers exist
if ! grep -q '<!-- BEGIN BENCHMARKS -->' "$README"; then
  echo "Error: README.md is missing <!-- BEGIN BENCHMARKS --> / <!-- END BENCHMARKS --> markers." >&2
  exit 1
fi

# Replace content between markers using awk + file read
awk '
  /<!-- BEGIN BENCHMARKS -->/ {
    print
    while ((getline line < "'"$SECTION_FILE"'") > 0) print line
    skip = 1
    next
  }
  /<!-- END BENCHMARKS -->/ { skip = 0 }
  skip { next }
  { print }
' "$README" > "$README.tmp"

mv "$README.tmp" "$README"

echo "Benchmarks updated in README.md"
