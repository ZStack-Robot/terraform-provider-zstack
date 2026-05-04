#!/usr/bin/env bash
# Copyright (c) ZStack.io, Inc.
# SPDX-License-Identifier: MPL-2.0

set -euo pipefail

# Scan all resource files for IsNull() without IsUnknown() guards
# Categorizes findings by risk level (HIGH for Int64/Bool, MEDIUM for String)
# Exit code: 0 if no findings, 1 if findings detected

cd "$(dirname "${BASH_SOURCE[0]}")/../.." || exit 1

findings=0

echo "=== HIGH RISK: Int64/Bool fields missing IsUnknown guard ==="
for f in zstack/provider/resource_zstack_*.go; do
    [[ "$f" == *_test.go ]] && continue
    [[ ! -f "$f" ]] && continue
    
    name=$(basename "$f" | sed 's/resource_zstack_//;s/\.go//')
    
    grep -n 'IsNull()' "$f" | grep -v 'IsUnknown' | while read -r line; do
        linenum=$(echo "$line" | cut -d: -f1)
        context=$(sed -n "$((linenum-3)),$((linenum+3))p" "$f")
        if echo "$context" | grep -qE 'ValueInt64|ValueBool|Int64Value|BoolValue'; then
            echo "$f:$linenum: HIGH — $(echo "$line" | cut -d: -f2- | xargs)"
            findings=$((findings + 1))
        fi
    done
done

echo ""
echo "=== MEDIUM RISK: String fields missing IsUnknown guard (summary) ==="
for f in zstack/provider/resource_zstack_*.go; do
    [[ "$f" == *_test.go ]] && continue
    [[ ! -f "$f" ]] && continue
    
    name=$(basename "$f" | sed 's/resource_zstack_//;s/\.go//')
    null_count=$(grep -c 'IsNull()' "$f" 2>/dev/null || echo 0)
    unknown_count=$(grep -c 'IsUnknown()' "$f" 2>/dev/null || echo 0)
    gap=$((null_count - unknown_count))
    
    if [ "$gap" -gt 0 ]; then
        echo "$f: MEDIUM — gap=$gap (IsNull=$null_count, IsUnknown=$unknown_count)"
        findings=$((findings + gap))
    fi
done

if [ "$findings" -gt 0 ]; then
    echo ""
    echo "Total findings: $findings"
    exit 1
else
    echo "No findings detected."
    exit 0
fi
