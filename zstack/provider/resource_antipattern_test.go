// Copyright (c) ZStack.io, Inc.

package provider

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestNoEmptyUUIDStateCorruption scans all resource files for the anti-pattern
// where a Read function rewrites state with an empty UUID/ID on non-not-found
// errors. This pattern causes state corruption: a transient API failure turns
// the resource into a zombie, and the subsequent Delete sees the empty UUID and
// skips remote deletion — orphaning the real resource in ZStack.
//
// The correct pattern is:
//
//	if errors.Is(err, ErrResourceNotFound) {
//	    response.State.RemoveResource(ctx)
//	    return
//	}
//	response.Diagnostics.AddError(...)
//	return
func TestNoEmptyUUIDStateCorruption(t *testing.T) {
	files, err := filepath.Glob("resource_zstack_*.go")
	if err != nil {
		t.Fatalf("failed to glob resource files: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("no resource files found — test may be running from wrong directory")
	}

	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}

		violations := scanForEmptyIDPattern(t, file)
		for _, v := range violations {
			t.Errorf("%s:%d: %s", file, v.line, v.msg)
		}
	}
}

type violation struct {
	line int
	msg  string
}

// scanForEmptyIDPattern looks for lines that assign Uuid or ID to an empty
// StringValue inside a Read function's error branch.
func scanForEmptyIDPattern(t *testing.T, filename string) []violation {
	t.Helper()

	f, err := os.Open(filename)
	if err != nil {
		t.Fatalf("failed to open %s: %v", filename, err)
	}
	defer f.Close()

	var violations []violation
	scanner := bufio.NewScanner(f)
	lineNum := 0
	inReadFunc := false

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if strings.Contains(line, "func (") && strings.Contains(line, ") Read(") {
			inReadFunc = true
		}
		if inReadFunc && lineNum > 1 && strings.HasPrefix(trimmed, "func ") && !strings.Contains(line, ") Read(") {
			inReadFunc = false
		}

		if !inReadFunc {
			continue
		}

		// Detect empty UUID/ID state rewrite:
		//   Uuid: types.StringValue("")
		//   ID: types.StringValue("")
		if strings.Contains(trimmed, `types.StringValue("")`) &&
			(strings.Contains(trimmed, `Uuid:`) || strings.Contains(trimmed, `ID:`)) {
			violations = append(violations, violation{
				line: lineNum,
				msg:  "Read function rewrites state with empty UUID/ID — transient errors will corrupt state and orphan remote resources",
			})
		}
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("error scanning %s: %v", filename, err)
	}

	return violations
}

// TestNoReadRemoveResourceOnTransientError scans Read functions for the
// anti-pattern where RemoveResource is called in a non-not-found error branch.
// This incorrectly removes the resource from state on transient failures.
//
// The pattern looks like:
//
//	if err != nil {
//	    if errors.Is(err, ErrResourceNotFound) {
//	        resp.State.RemoveResource(ctx) // correct
//	        return
//	    }
//	    tflog.Warn(...)
//	    resp.State.RemoveResource(ctx)     // BUG: transient error removes state
//	    return
//	}
func TestNoReadRemoveResourceOnTransientError(t *testing.T) {
	files, err := filepath.Glob("resource_zstack_*.go")
	if err != nil {
		t.Fatalf("failed to glob resource files: %v", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}

		violations := scanForReadRemoveResourceOnError(t, file)
		for _, v := range violations {
			t.Errorf("%s:%d: %s", file, v.line, v.msg)
		}
	}
}

// scanForReadRemoveResourceOnError detects RemoveResource calls in Read
// that follow a tflog.Warn outside of the ErrResourceNotFound / isZStackNotFoundError branch.
func scanForReadRemoveResourceOnError(t *testing.T, filename string) []violation {
	t.Helper()

	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("failed to read %s: %v", filename, err)
	}

	lines := strings.Split(string(data), "\n")
	var violations []violation
	inReadFunc := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.Contains(line, "func (") && strings.Contains(line, ") Read(") {
			inReadFunc = true
		}
		if inReadFunc && i > 0 && strings.HasPrefix(trimmed, "func ") && !strings.Contains(line, ") Read(") {
			inReadFunc = false
		}

		if !inReadFunc {
			continue
		}

		// Look for tflog.Warn followed within 3 lines by RemoveResource,
		// but ONLY if the warn is NOT inside a not-found check.
		if strings.Contains(trimmed, "tflog.Warn(") {
			// Check if this Warn is inside a not-found branch by looking
			// at the preceding 5 lines for ErrResourceNotFound or isZStackNotFoundError
			inNotFoundBranch := false
			for j := max(0, i-5); j < i; j++ {
				prevLine := lines[j]
				if strings.Contains(prevLine, "ErrResourceNotFound") ||
					strings.Contains(prevLine, "isZStackNotFoundError") ||
					strings.Contains(prevLine, "len(") {
					inNotFoundBranch = true
					break
				}
			}
			if inNotFoundBranch {
				continue
			}

			for j := i + 1; j < len(lines) && j <= i+3; j++ {
				nextTrimmed := strings.TrimSpace(lines[j])
				if strings.Contains(nextTrimmed, "RemoveResource(") {
					violations = append(violations, violation{
						line: j + 1,
						msg:  "Read function calls RemoveResource after tflog.Warn on non-not-found error — transient failures will incorrectly remove state",
					})
					break
				}
				if strings.Contains(nextTrimmed, "AddError(") || nextTrimmed == "}" {
					break
				}
			}
		}
	}

	return violations
}

// TestNoDeleteEmptyUUIDGuard scans all resource files for the anti-pattern
// where a Delete function checks for empty UUID and silently skips deletion.
// This guard only exists to handle the corrupted state from the Read
// anti-pattern above. Once Read properly returns errors, this guard becomes
// dead code that masks bugs.
func TestNoDeleteEmptyUUIDGuard(t *testing.T) {
	files, err := filepath.Glob("resource_zstack_*.go")
	if err != nil {
		t.Fatalf("failed to glob resource files: %v", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}

		violations := scanForDeleteEmptyUUIDGuard(t, file)
		for _, v := range violations {
			t.Errorf("%s:%d: %s", file, v.line, v.msg)
		}
	}
}

func scanForDeleteEmptyUUIDGuard(t *testing.T, filename string) []violation {
	t.Helper()

	f, err := os.Open(filename)
	if err != nil {
		t.Fatalf("failed to open %s: %v", filename, err)
	}
	defer f.Close()

	var violations []violation
	scanner := bufio.NewScanner(f)
	lineNum := 0
	inDeleteFunc := false

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if strings.Contains(line, "func (") && strings.Contains(line, ") Delete(") {
			inDeleteFunc = true
		}
		if inDeleteFunc && lineNum > 1 && strings.HasPrefix(trimmed, "func ") && !strings.Contains(line, ") Delete(") {
			inDeleteFunc = false
		}

		if !inDeleteFunc {
			continue
		}

		// Detect both variants:
		//   if state.Uuid == types.StringValue("")
		//   if uuid == ""  /  if state.Uuid.ValueString() == ""
		if strings.Contains(trimmed, `state.Uuid == types.StringValue("")`) {
			violations = append(violations, violation{
				line: lineNum,
				msg:  "Delete function has empty-UUID guard (types.StringValue variant) — masks state corruption bugs",
			})
		}
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("error scanning %s: %v", filename, err)
	}

	return violations
}
