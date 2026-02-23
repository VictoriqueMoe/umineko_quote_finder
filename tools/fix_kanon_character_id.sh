#!/bin/bash

# Fixes Kanon's Episode 2 lines (quote IDs 20600528-20600543) that are
# incorrectly tagged with character ID "46" (Erika) instead of "16" (Kanon).

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
DATA_DIR="$SCRIPT_DIR/../internal/quote/data"

for file in "$DATA_DIR"/*.txt; do
    if [ ! -f "$file" ]; then
        echo "Skipping $file (not found)"
        continue
    fi

    count=$(grep -c '"46"\*"2060' "$file" 2>/dev/null || true)
    if [ "$count" -eq 0 ]; then
        echo "No changes needed in $file"
        continue
    fi

    sed -i '' 's/"46"\*"2060/"16"\*"2060/g' "$file"
    echo "Fixed $count occurrence(s) in $file"
done
