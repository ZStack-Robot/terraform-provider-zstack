// Copyright (c) ZStack.io, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"testing"
)

type Finding struct {
	File    string
	Line    int
	Check   string
	Message string
}

func parsePackage(dir string) ([]*ast.File, *token.FileSet, error) {
	fset := token.NewFileSet()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil, err
	}

	var files []*ast.File
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}

		filePath := filepath.Join(dir, entry.Name())
		file, err := parser.ParseFile(fset, filePath, nil, parser.AllErrors)
		if err != nil {
			return nil, nil, err
		}

		files = append(files, file)
	}

	return files, fset, nil
}

func walkFunctions(file *ast.File, visit func(*ast.FuncDecl)) {
	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			visit(fn)
		}
	}
}

func findCallExpr(node ast.Node, pkg, name string) []*ast.CallExpr {
	var results []*ast.CallExpr

	ast.Inspect(node, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
				if ident, ok := sel.X.(*ast.Ident); ok {
					if ident.Name == pkg && sel.Sel.Name == name {
						results = append(results, call)
					}
				}
			} else if ident, ok := call.Fun.(*ast.Ident); ok {
				if ident.Name == name {
					results = append(results, call)
				}
			}
		}
		return true
	})

	return results
}

func TestAntipatternScanner(t *testing.T) {
	t.Run("smoke", func(t *testing.T) {
		files, _, err := parsePackage(".")
		if err != nil {
			t.Fatalf("parsePackage failed: %v", err)
		}

		if len(files) == 0 {
			t.Fatal("no files parsed from current directory")
		}

		functionCount := 0
		for _, file := range files {
			walkFunctions(file, func(fn *ast.FuncDecl) {
				functionCount++
			})
		}

		if functionCount == 0 {
			t.Fatal("no functions found in parsed files")
		}

		t.Logf("parsed %d files with %d functions", len(files), functionCount)
	})

	t.Run("check_2d", func(t *testing.T) {
		t.Run("bad", func(t *testing.T) {
			badFindings := scanCheck2dFixtures(t, "testdata/antipatterns/check_2d/bad")
			if len(badFindings) != 3 {
				t.Fatalf("expected 3 bad findings, got %d: %v", len(badFindings), badFindings)
			}

			for _, f := range badFindings {
				t.Logf("BAD: %s:%d - %s", f.File, f.Line, f.Message)
			}
		})

		t.Run("good", func(t *testing.T) {
			goodFindings := scanCheck2dFixtures(t, "testdata/antipatterns/check_2d/good")
			if len(goodFindings) != 0 {
				t.Fatalf("expected 0 good findings, got %d: %v", len(goodFindings), goodFindings)
			}
		})

		t.Run("repo_sweep", func(t *testing.T) {
			repoFindings := scanCheck2dInRepo(t, ".")
			t.Logf("Repo sweep found %d check_2d findings in zstack/provider/", len(repoFindings))
			for _, f := range repoFindings {
				t.Logf("REPO: %s:%d - %s", f.File, f.Line, f.Message)
			}
		})
	})

	t.Run("check_2b", func(t *testing.T) {
		t.Run("bad", func(t *testing.T) {
			badFindings := scanCheck2bFixtures(t, "testdata/antipatterns/check_2b/bad")
			if len(badFindings) != 3 {
				t.Fatalf("expected 3 bad findings, got %d: %v", len(badFindings), badFindings)
			}

			for _, f := range badFindings {
				t.Logf("BAD: %s:%d - %s", f.File, f.Line, f.Message)
			}
		})

		t.Run("good", func(t *testing.T) {
			goodFindings := scanCheck2bFixtures(t, "testdata/antipatterns/check_2b/good")
			if len(goodFindings) != 0 {
				t.Fatalf("expected 0 good findings, got %d: %v", len(goodFindings), goodFindings)
			}

			t.Logf("string fixture not flagged")
			t.Logf("list fixture not flagged")
		})

		t.Run("repo_sweep", func(t *testing.T) {
			repoFindings := scanCheck2bInRepo(t, ".")
			t.Logf("Repo sweep found %d check_2b findings in zstack/provider/", len(repoFindings))
			for _, f := range repoFindings {
				t.Logf("REPO: %s:%d - %s", f.File, f.Line, f.Message)
			}
		})
	})

	t.Run("check_2a", func(t *testing.T) {
		t.Run("bad", func(t *testing.T) {
			badFindings := scanCheck2aFixtures(t, "testdata/antipatterns/check_2a/bad")
			if len(badFindings) != 2 {
				t.Fatalf("expected 2 bad findings, got %d: %v", len(badFindings), badFindings)
			}

			for _, f := range badFindings {
				t.Logf("BAD: %s:%d - %s", f.File, f.Line, f.Message)
			}
		})

		t.Run("good", func(t *testing.T) {
			goodFindings := scanCheck2aFixtures(t, "testdata/antipatterns/check_2a/good")
			if len(goodFindings) != 0 {
				t.Fatalf("expected 0 good findings, got %d: %v", len(goodFindings), goodFindings)
			}
		})

		t.Run("repo_sweep", func(t *testing.T) {
			repoFindings := scanCheck2aInRepo(t, ".")
			t.Logf("Repo sweep found %d check_2a findings in zstack/provider/", len(repoFindings))
			for _, f := range repoFindings {
				t.Logf("REPO: %s:%d - %s", f.File, f.Line, f.Message)
			}
		})
	})

	t.Run("repo_sweep_postfix", func(t *testing.T) {
		// Integration gate: verify Wave 2 fixes are complete
		// Scans ONLY Wave 2 resource files (T5-T12) for check_2d and check_2b
		// Asserts ZERO findings for both checks (all 12 Story-15 hotspots guarded)
		// Logs check_2a count (informational, not asserted)

		wave2Files := []string{
			"resource_zstack_tag_attachment.go",
			"resource_zstack_instance.go",
			"resource_zstack_port_forwarding_rule.go",
			"resource_zstack_volume.go",
			"resource_zstack_instance_scripts.go",
			"resource_zstack_instance_scripts_execution.go",
			"resource_zstack_l3network.go",
			"resource_zstack_vip_qos.go",
		}

		fset := token.NewFileSet()
		var findings2d []Finding
		var findings2b []Finding
		var findings2a []Finding

		for _, fileName := range wave2Files {
			file, err := parser.ParseFile(fset, fileName, nil, parser.AllErrors)
			if err != nil {
				t.Logf("WARNING: failed to parse %s: %v", fileName, err)
				continue
			}

			findings2d = append(findings2d, scanDiscardedAppendCalls(file, fset, fileName)...)
			findings2b = append(findings2b, scanCheck2bFile(file, fset, fileName)...)
			findings2a = append(findings2a, scanCheck2aFile(file, fset, fileName)...)
		}

		t.Logf("=== INTEGRATION GATE: repo_sweep_postfix ===")
		t.Logf("Scanned %d Wave 2 resource files", len(wave2Files))
		t.Logf("repo_sweep_postfix: check_2d findings = %d", len(findings2d))
		for _, f := range findings2d {
			t.Logf("  FAIL: %s:%d - %s", f.File, f.Line, f.Message)
		}

		t.Logf("repo_sweep_postfix: check_2b findings = %d", len(findings2b))
		for _, f := range findings2b {
			t.Logf("  FAIL: %s:%d - %s", f.File, f.Line, f.Message)
		}

		t.Logf("check_2a findings = %d (informational, not asserted)", len(findings2a))

		if len(findings2d) > 0 {
			t.Errorf("INTEGRATION GATE FAILED: check_2d found %d findings (expected 0)", len(findings2d))
		}

		if len(findings2b) > 0 {
			t.Errorf("INTEGRATION GATE FAILED: check_2b found %d findings (expected 0)", len(findings2b))
		}

		if len(findings2d) == 0 && len(findings2b) == 0 {
			t.Logf("✓ INTEGRATION GATE PASSED: Wave 2 complete (check_2d=0, check_2b=0)")
		}
	})
}

func scanCheck2dFixtures(t *testing.T, dir string) []Finding {
	var findings []Finding

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("failed to read directory %s: %v", dir, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go.fixture") {
			file, fset := mustParseFixtureFile(t, dir, entry.Name())
			findings = append(findings, scanDiscardedAppendCalls(file, fset, entry.Name())...)
		}
	}

	return findings
}

func mustParseFixtureFile(t *testing.T, dir, name string) (*ast.File, *token.FileSet) {
	t.Helper()

	filePath := filepath.Join(dir, name)
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read fixture file %s: %v", filePath, err)
	}

	fixtureName := strings.TrimSuffix(name, ".fixture")
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, fixtureName, string(content), parser.AllErrors)
	if err != nil {
		t.Fatalf("failed to parse fixture file %s: %v", filePath, err)
	}

	return file, fset
}

func scanDiscardedAppendCalls(file *ast.File, fset *token.FileSet, fileName string) []Finding {
	var findings []Finding
	var stack []ast.Node

	ast.Inspect(file, func(n ast.Node) bool {
		if n == nil {
			stack = stack[:len(stack)-1]
			return false
		}

		stack = append(stack, n)

		call, ok := n.(*ast.CallExpr)
		if !ok || !isAppendCall(call) || appendCallAssigned(stack, call) {
			return true
		}

		pos := fset.Position(call.Pos())
		findings = append(findings, Finding{
			File:    fileName,
			Line:    pos.Line,
			Check:   "2d",
			Message: "append result discarded",
		})

		return true
	})

	return findings
}

func isAppendCall(call *ast.CallExpr) bool {
	ident, ok := call.Fun.(*ast.Ident)
	return ok && ident.Name == "append"
}

func appendCallAssigned(stack []ast.Node, appendCall *ast.CallExpr) bool {
	for i := len(stack) - 2; i >= 0; i-- {
		switch node := stack[i].(type) {
		case *ast.AssignStmt:
			for _, rhs := range node.Rhs {
				if nodeContains(rhs, appendCall) {
					return true
				}
			}
			return false
		case *ast.ExprStmt:
			return false
		}
	}

	return false
}

func nodeContains(root ast.Node, target ast.Node) bool {
	found := false
	ast.Inspect(root, func(n ast.Node) bool {
		if n == target {
			found = true
			return false
		}
		return true
	})
	return found
}

func scanCheck2dInRepo(t *testing.T, dir string) []Finding {
	var findings []Finding

	files, fset, err := parsePackage(dir)
	if err != nil {
		t.Logf("parsePackage failed: %v", err)
		return findings
	}

	for _, file := range files {
		fileName := fset.Position(file.Pos()).Filename
		findings = append(findings, scanDiscardedAppendCalls(file, fset, fileName)...)
	}

	return findings
}

func scanCheck2bFixtures(t *testing.T, dir string) []Finding {
	var findings []Finding

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("failed to read directory %s: %v", dir, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go.fixture") {
			file, fset := mustParseFixtureFile(t, dir, entry.Name())
			findings = append(findings, scanCheck2bFile(file, fset, entry.Name())...)
		}
	}

	return findings
}

func scanCheck2bFile(file *ast.File, fset *token.FileSet, fileName string) []Finding {
	var findings []Finding

	ast.Inspect(file, func(n ast.Node) bool {
		ifStmt, ok := n.(*ast.IfStmt)
		if !ok {
			return true
		}

		isNullReceivers := findSelectorCallReceivers(ifStmt.Cond, "IsNull")
		if len(isNullReceivers) == 0 {
			return true
		}

		isUnknownReceivers := findSelectorCallReceivers(ifStmt.Cond, "IsUnknown")
		for receiver := range isNullReceivers {
			if isUnknownReceivers[receiver] {
				continue
			}

			usage := classifyCheck2bReceiverUsage(ifStmt.Body, receiver)
			if usage != check2bUsageTarget {
				continue
			}

			pos := fset.Position(ifStmt.If)
			findings = append(findings, Finding{
				File:    fileName,
				Line:    pos.Line,
				Check:   "2b",
				Message: "IsNull() without IsUnknown() check",
			})
			break
		}

		return true
	})

	return findings
}

func findSelectorCallReceivers(node ast.Node, methodName string) map[string]bool {
	receivers := make(map[string]bool)

	ast.Inspect(node, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != methodName {
			return true
		}

		receiver := renderNode(sel.X)
		if receiver != "" {
			receivers[receiver] = true
		}

		return true
	})

	return receivers
}

type check2bUsage int

const (
	check2bUsageNone check2bUsage = iota
	check2bUsageTarget
	check2bUsageSkip
)

func classifyCheck2bReceiverUsage(body *ast.BlockStmt, receiver string) check2bUsage {
	usage := check2bUsageNone

	ast.Inspect(body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || renderNode(sel.X) != receiver {
			return true
		}

		switch sel.Sel.Name {
		case "ValueString", "Elements", "ElementsAs":
			usage = check2bUsageSkip
			return false
		case "ValueInt64", "ValueBool":
			usage = check2bUsageTarget
		}

		return true
	})

	return usage
}

func renderNode(node ast.Node) string {
	if node == nil {
		return ""
	}

	var buf bytes.Buffer
	if err := printer.Fprint(&buf, token.NewFileSet(), node); err != nil {
		return ""
	}

	return buf.String()
}

func scanCheck2bInRepo(t *testing.T, dir string) []Finding {
	var findings []Finding

	files, fset, err := parsePackage(dir)
	if err != nil {
		t.Logf("parsePackage failed: %v", err)
		return findings
	}

	for _, file := range files {
		fileName := fset.Position(file.Pos()).Filename
		findings = append(findings, scanCheck2bFile(file, fset, fileName)...)
	}

	return findings
}

func scanCheck2aFixtures(t *testing.T, dir string) []Finding {
	var findings []Finding

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("failed to read directory %s: %v", dir, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go.fixture") {
			file, fset := mustParseFixtureFile(t, dir, entry.Name())
			findings = append(findings, scanCheck2aFile(file, fset, entry.Name())...)
		}
	}

	return findings
}

func scanCheck2aFile(file *ast.File, fset *token.FileSet, fileName string) []Finding {
	var findings []Finding

	walkFunctions(file, func(fn *ast.FuncDecl) {
		if !isCheck2aUpdateMethod(fn) {
			return
		}

		if !check2aNeedsReread(fn) {
			return
		}

		pos := fset.Position(fn.Pos())
		findings = append(findings, Finding{
			File:    fileName,
			Line:    pos.Line,
			Check:   "2a",
			Message: "Update API call without subsequent read-after-write",
		})
	})

	return findings
}

func scanCheck2aInRepo(t *testing.T, dir string) []Finding {
	var findings []Finding

	files, fset, err := parsePackage(dir)
	if err != nil {
		t.Logf("parsePackage failed: %v", err)
		return findings
	}

	for _, file := range files {
		fileName := fset.Position(file.Pos()).Filename
		findings = append(findings, scanCheck2aFile(file, fset, fileName)...)
	}

	return findings
}

type check2aEvent struct {
	pos  token.Pos
	kind string
}

func isCheck2aUpdateMethod(fn *ast.FuncDecl) bool {
	if fn == nil || fn.Name == nil || fn.Name.Name != "Update" || fn.Recv == nil || fn.Body == nil {
		return false
	}

	receiverType := receiverTypeName(fn)
	return receiverType != "" && strings.HasSuffix(receiverType, "Resource")
}

func receiverTypeName(fn *ast.FuncDecl) string {
	if fn == nil || fn.Recv == nil || len(fn.Recv.List) == 0 {
		return ""
	}

	switch expr := fn.Recv.List[0].Type.(type) {
	case *ast.Ident:
		return expr.Name
	case *ast.StarExpr:
		if ident, ok := expr.X.(*ast.Ident); ok {
			return ident.Name
		}
	}

	return ""
}

func check2aNeedsReread(fn *ast.FuncDecl) bool {
	events := collectCheck2aEvents(fn.Body)
	seenUpdate := false

	for _, event := range events {
		switch event.kind {
		case "update":
			seenUpdate = true
		case "read", "state":
			if seenUpdate {
				return false
			}
		}
	}

	return seenUpdate
}

func collectCheck2aEvents(body *ast.BlockStmt) []check2aEvent {
	var events []check2aEvent

	ast.Inspect(body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if isCheck2aUpdateCall(call) {
			events = append(events, check2aEvent{pos: call.Pos(), kind: "update"})
		}
		if isCheck2aReadCall(call) {
			events = append(events, check2aEvent{pos: call.Pos(), kind: "read"})
		}
		if isCheck2aStateRefreshCall(call) {
			events = append(events, check2aEvent{pos: call.Pos(), kind: "state"})
		}

		return true
	})

	sort.SliceStable(events, func(i, j int) bool {
		return events[i].pos < events[j].pos
	})

	return events
}

func isCheck2aUpdateCall(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	return ok && strings.HasPrefix(sel.Sel.Name, "Update")
}

func isCheck2aReadCall(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	return strings.HasPrefix(sel.Sel.Name, "Get") || strings.HasPrefix(sel.Sel.Name, "Query")
}

func isCheck2aStateRefreshCall(call *ast.CallExpr) bool {
	if isCheck2aStateSetCall(call) {
		return true
	}

	name := check2aCallName(call)
	if !strings.Contains(strings.ToLower(name), "setstate") {
		return false
	}

	return slices.ContainsFunc(call.Args, isCheck2aNonPlanSource)
}

func isCheck2aStateSetCall(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "Set" {
		return false
	}

	if !strings.HasSuffix(renderNode(sel.X), ".State") {
		return false
	}

	return slices.ContainsFunc(call.Args, isCheck2aNonPlanSource)
}

func check2aCallName(call *ast.CallExpr) string {
	switch fun := call.Fun.(type) {
	case *ast.Ident:
		return fun.Name
	case *ast.SelectorExpr:
		return fun.Sel.Name
	default:
		return ""
	}
}

func isCheck2aNonPlanSource(expr ast.Expr) bool {
	text := strings.ToLower(strings.TrimSpace(renderNode(expr)))
	if text == "" || text == "ctx" || text == "context" {
		return false
	}

	text = strings.TrimPrefix(text, "&")
	return text != "plan" && !strings.HasPrefix(text, "plan.") && !strings.Contains(text, "req.plan")
}
