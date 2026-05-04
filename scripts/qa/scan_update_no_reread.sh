#!/usr/bin/env bash
# Copyright (c) ZStack.io, Inc.
# SPDX-License-Identifier: MPL-2.0

set -euo pipefail

# Scan all resource Update methods to detect if they trust SDK return values
# without performing read-after-write verification
# Exit code: 0 if no dangerous patterns found, 1 if dangerous patterns detected

cd "$(dirname "${BASH_SOURCE[0]}")/../.." || exit 1

dangerous=0
safe=0

for f in zstack/provider/resource_zstack_*.go; do
    [[ "$f" == *_test.go ]] && continue
    [[ ! -f "$f" ]] && continue
    
    name=$(basename "$f" | sed 's/resource_zstack_//;s/\.go//')
    
    update_line=$(grep -n 'r\.client\.Update' "$f" 2>/dev/null | head -1)
    [ -z "$update_line" ] && continue
    
    linenum=$(echo "$update_line" | cut -d: -f1)
    update_end=$((linenum + 40))
    
    has_reread=$(sed -n "${linenum},${update_end}p" "$f" 2>/dev/null | grep -c 'findResourceByQuery\|r\.client\.Query\|r\.client\.Get' 2>/dev/null || true)
    
    if [ "$has_reread" -gt 0 ]; then
        safe=$((safe + 1))
        echo "$f:$linenum: SAFE"
    else
        dangerous=$((dangerous + 1))
        echo "$f:$linenum: DANGER — no read-after-write"
    fi
done

echo ""
echo "Summary: safe=$safe, dangerous=$dangerous"

if [ "$dangerous" -gt 0 ]; then
    exit 1
else
    exit 0
fi
