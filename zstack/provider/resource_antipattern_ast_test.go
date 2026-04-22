// Copyright (c) ZStack.io, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// Finding represents a single antipattern finding from AST scanning.
type Finding struct {
	File    string
	Line    int
	Check   string
	Message string
}

// parsePackage parses all Go files in a directory and returns the AST files and token set.
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

// walkFunctions walks all function declarations in a file and calls visit for each.
func walkFunctions(file *ast.File, visit func(*ast.FuncDecl)) {
	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			visit(fn)
		}
	}
}

// findCallExpr finds all call expressions matching a package and function name.
func findCallExpr(node ast.Node, pkg, name string) []*ast.CallExpr {
	var results []*ast.CallExpr

	ast.Inspect(node, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			// Check if this is a call to pkg.name or just name
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

// TestAntipatternScanner is the top-level test harness for AST-based antipattern scanning.
func TestAntipatternScanner(t *testing.T) {
	t.Run("smoke", func(t *testing.T) {
		// Smoke test: parse resource_zstack_volume.go without error
		files, _, err := parsePackage(".")
		if err != nil {
			t.Fatalf("parsePackage failed: %v", err)
		}

		if len(files) == 0 {
			t.Fatal("no files parsed from current directory")
		}

		// Verify we can walk functions
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
}
