// Copyright (c) ZStack.io, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path/filepath"
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
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.AllErrors)
	if err != nil {
		return nil, nil, err
	}

	var files []*ast.File
	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			files = append(files, f)
		}
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
		badFindings := scanCheck2dFixtures(t, "testdata/antipatterns/check_2d/bad")
		goodFindings := scanCheck2dFixtures(t, "testdata/antipatterns/check_2d/good")

		if len(badFindings) < 2 {
			t.Errorf("expected at least 2 bad findings, got %d", len(badFindings))
		}
		if len(goodFindings) > 0 {
			t.Errorf("expected 0 good findings, got %d: %v", len(goodFindings), goodFindings)
		}

		for _, f := range badFindings {
			t.Logf("BAD: %s:%d - %s", f.File, f.Line, f.Message)
		}

		repoFindings := scanCheck2dInRepo(t, ".")
		t.Logf("Repo sweep found %d check_2d findings in zstack/provider/", len(repoFindings))
		for _, f := range repoFindings {
			t.Logf("REPO: %s:%d - %s", f.File, f.Line, f.Message)
		}
	})

	t.Run("check_2b", func(t *testing.T) {
		badFindings := scanCheck2bFixtures(t, "testdata/antipatterns/check_2b/bad")
		goodFindings := scanCheck2bFixtures(t, "testdata/antipatterns/check_2b/good")

		if len(badFindings) < 2 {
			t.Errorf("expected at least 2 bad findings, got %d", len(badFindings))
		}
		if len(goodFindings) > 0 {
			t.Errorf("expected 0 good findings, got %d: %v", len(goodFindings), goodFindings)
		}

		for _, f := range badFindings {
			t.Logf("BAD: %s:%d - %s", f.File, f.Line, f.Message)
		}

		repoFindings := scanCheck2bInRepo(t, ".")
		t.Logf("Repo sweep found %d check_2b findings in zstack/provider/", len(repoFindings))
		for _, f := range repoFindings {
			t.Logf("REPO: %s:%d - %s", f.File, f.Line, f.Message)
		}
	})

	t.Run("check_2a", func(t *testing.T) {
		badFindings := scanCheck2aFixtures(t, "testdata/antipatterns/check_2a/bad")
		goodFindings := scanCheck2aFixtures(t, "testdata/antipatterns/check_2a/good")

		if len(badFindings) < 1 {
			t.Errorf("expected at least 1 bad finding, got %d", len(badFindings))
		}
		if len(goodFindings) > 0 {
			t.Errorf("expected 0 good findings, got %d: %v", len(goodFindings), goodFindings)
		}

		for _, f := range badFindings {
			t.Logf("BAD: %s:%d - %s", f.File, f.Line, f.Message)
		}

		repoFindings := scanCheck2aInRepo(t, ".")
		t.Logf("Repo sweep found %d check_2a findings in zstack/provider/", len(repoFindings))
		for _, f := range repoFindings {
			t.Logf("REPO: %s:%d - %s", f.File, f.Line, f.Message)
		}
	})
}

func scanCheck2dFixtures(t *testing.T, dir string) []Finding {
	var findings []Finding

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatalf("failed to read directory %s: %v", dir, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go.fixture") {
			filePath := filepath.Join(dir, entry.Name())
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				t.Fatalf("failed to read fixture file %s: %v", filePath, err)
			}

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, entry.Name(), string(content), parser.AllErrors)
			if err != nil {
				t.Fatalf("failed to parse fixture file %s: %v", filePath, err)
			}

			ast.Inspect(file, func(n ast.Node) bool {
				if call, ok := n.(*ast.CallExpr); ok {
					if ident, ok := call.Fun.(*ast.Ident); ok && ident.Name == "append" {
						if !isAppendAssigned(call, file) {
							pos := fset.Position(call.Pos())
							findings = append(findings, Finding{
								File:    entry.Name(),
								Line:    pos.Line,
								Check:   "2d",
								Message: "append result discarded",
							})
						}
					}
				}
				return true
			})
		}
	}

	return findings
}

func isAppendAssigned(appendCall *ast.CallExpr, file *ast.File) bool {
	var isAssigned bool

	ast.Inspect(file, func(n ast.Node) bool {
		if assign, ok := n.(*ast.AssignStmt); ok {
			for _, rhs := range assign.Rhs {
				if rhs == appendCall {
					isAssigned = true
					return false
				}
			}
		}
		return true
	})

	return isAssigned
}

func scanCheck2dInRepo(t *testing.T, dir string) []Finding {
	var findings []Finding

	files, fset, err := parsePackage(dir)
	if err != nil {
		t.Logf("parsePackage failed: %v", err)
		return findings
	}

	for _, file := range files {
		ast.Inspect(file, func(n ast.Node) bool {
			if call, ok := n.(*ast.CallExpr); ok {
				if ident, ok := call.Fun.(*ast.Ident); ok && ident.Name == "append" {
					if !isAppendAssigned(call, file) {
						pos := fset.Position(call.Pos())
						findings = append(findings, Finding{
							File:    pos.Filename,
							Line:    pos.Line,
							Check:   "2d",
							Message: "append result discarded",
						})
					}
				}
			}
			return true
		})
	}

	return findings
}

func scanCheck2bFixtures(t *testing.T, dir string) []Finding {
	var findings []Finding

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatalf("failed to read directory %s: %v", dir, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go.fixture") {
			filePath := filepath.Join(dir, entry.Name())
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				t.Fatalf("failed to read fixture file %s: %v", filePath, err)
			}

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, entry.Name(), string(content), parser.AllErrors)
			if err != nil {
				t.Fatalf("failed to parse fixture file %s: %v", filePath, err)
			}

		ast.Inspect(file, func(n ast.Node) bool {
			if call, ok := n.(*ast.CallExpr); ok {
				if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
					if sel.Sel.Name == "IsNull" {
						if !hasIsUnknownCheck(call, file) {
							pos := fset.Position(call.Pos())
							findings = append(findings, Finding{
								File:    entry.Name(),
								Line:    pos.Line,
								Check:   "2b",
								Message: "IsNull() without IsUnknown() check",
							})
						}
					}
				}
			}
			return true
		})
		}
	}

	return findings
}

func hasIsUnknownCheck(isNullCall *ast.CallExpr, file *ast.File) bool {
	var hasCheck bool

	ast.Inspect(file, func(n ast.Node) bool {
		if binExpr, ok := n.(*ast.BinaryExpr); ok {
			if binExpr.Op == token.LOR || binExpr.Op == token.LAND {
				if containsCall(binExpr.X, isNullCall) && containsIsUnknown(binExpr.Y) {
					hasCheck = true
					return false
				}
				if containsCall(binExpr.Y, isNullCall) && containsIsUnknown(binExpr.X) {
					hasCheck = true
					return false
				}
			}
		}
		return true
	})

	return hasCheck
}

func containsCall(expr ast.Expr, target *ast.CallExpr) bool {
	var found bool
	ast.Inspect(expr, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok && call == target {
			found = true
			return false
		}
		return true
	})
	return found
}

func containsIsUnknown(expr ast.Expr) bool {
	var found bool
	ast.Inspect(expr, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
				if sel.Sel.Name == "IsUnknown" {
					found = true
					return false
				}
			}
		}
		return true
	})
	return found
}

func scanCheck2bInRepo(t *testing.T, dir string) []Finding {
	var findings []Finding

	files, fset, err := parsePackage(dir)
	if err != nil {
		t.Logf("parsePackage failed: %v", err)
		return findings
	}

	for _, file := range files {
		ast.Inspect(file, func(n ast.Node) bool {
			if call, ok := n.(*ast.CallExpr); ok {
				if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
					if sel.Sel.Name == "IsNull" {
						if !hasIsUnknownCheck(call, file) {
							pos := fset.Position(call.Pos())
							findings = append(findings, Finding{
								File:    pos.Filename,
								Line:    pos.Line,
								Check:   "2b",
								Message: "IsNull() without IsUnknown() check",
							})
						}
					}
				}
			}
			return true
		})
	}

	return findings
}

func scanCheck2aFixtures(t *testing.T, dir string) []Finding {
	var findings []Finding

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatalf("failed to read directory %s: %v", dir, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go.fixture") {
			filePath := filepath.Join(dir, entry.Name())
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				t.Fatalf("failed to read fixture file %s: %v", filePath, err)
			}

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, entry.Name(), string(content), parser.AllErrors)
			if err != nil {
				t.Fatalf("failed to parse fixture file %s: %v", filePath, err)
			}

			ast.Inspect(file, func(n ast.Node) bool {
				if exprStmt, ok := n.(*ast.ExprStmt); ok {
					if indexExpr, ok := exprStmt.X.(*ast.IndexExpr); ok {
						pos := fset.Position(indexExpr.Pos())
						findings = append(findings, Finding{
							File:    entry.Name(),
							Line:    pos.Line,
							Check:   "2a",
							Message: "field read from API but not assigned to state",
						})
					}
				}
				return true
			})
		}
	}

	return findings
}

func isIndexExprAssigned(indexExpr *ast.IndexExpr, file *ast.File) bool {
	var isAssigned bool

	ast.Inspect(file, func(n ast.Node) bool {
		if assign, ok := n.(*ast.AssignStmt); ok {
			for _, rhs := range assign.Rhs {
				if rhs == indexExpr {
					isAssigned = true
					return false
				}
			}
		}
		return true
	})

	return isAssigned
}

func scanCheck2aInRepo(t *testing.T, dir string) []Finding {
	var findings []Finding

	files, fset, err := parsePackage(dir)
	if err != nil {
		t.Logf("parsePackage failed: %v", err)
		return findings
	}

	for _, file := range files {
		ast.Inspect(file, func(n ast.Node) bool {
			if exprStmt, ok := n.(*ast.ExprStmt); ok {
				if indexExpr, ok := exprStmt.X.(*ast.IndexExpr); ok {
					pos := fset.Position(indexExpr.Pos())
					findings = append(findings, Finding{
						File:    pos.Filename,
						Line:    pos.Line,
						Check:   "2a",
						Message: "field read from API but not assigned to state",
					})
				}
			}
			return true
		})
	}

	return findings
}
