// Copyright (c) ZStack.io, Inc.

package provider

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestNoEmptyUUIDStateCorruption scans all resource files for the anti-pattern
// where a Read function rewrites state with an empty UUID/ID on non-not-found
// errors. This pattern causes state corruption: a transient API failure turns
// the resource into a zombie, and the subsequent Delete sees the empty UUID and
// skips remote deletion — orphaning the real resource in ZStack.
func TestNoEmptyUUIDStateCorruption(t *testing.T) {
	for _, file := range resourceFiles(t) {
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

func resourceFiles(t *testing.T) []string {
	t.Helper()
	_, thisFile, _, _ := runtime.Caller(0)
	dir := filepath.Dir(thisFile)
	files, err := filepath.Glob(filepath.Join(dir, "resource_zstack_*.go"))
	if err != nil {
		t.Fatalf("failed to glob resource files: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("no resource files found — test may be running from wrong directory")
	}

	filtered := files[:0]
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		filtered = append(filtered, file)
	}
	return filtered
}

func parseGoFile(t *testing.T, filename string) (*token.FileSet, *ast.File) {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		t.Fatalf("failed to parse %s: %v", filename, err)
	}
	return fset, file
}

func scanForEmptyIDPattern(t *testing.T, filename string) []violation {
	t.Helper()
	fset, file := parseGoFile(t, filename)

	var violations []violation
	ast.Inspect(file, func(node ast.Node) bool {
		fn, ok := node.(*ast.FuncDecl)
		if !ok || fn.Name == nil || fn.Name.Name != "Read" || fn.Body == nil {
			return true
		}

		ast.Inspect(fn.Body, func(inner ast.Node) bool {
			lit, ok := inner.(*ast.CompositeLit)
			if !ok {
				return true
			}

			for _, elt := range lit.Elts {
				kv, ok := elt.(*ast.KeyValueExpr)
				if !ok {
					continue
				}
				key, ok := kv.Key.(*ast.Ident)
				if !ok || (key.Name != "Uuid" && key.Name != "ID") {
					continue
				}
				if isTypesStringValueEmpty(kv.Value) {
					violations = append(violations, violation{
						line: fset.Position(kv.Pos()).Line,
						msg:  "Read function rewrites state with empty UUID/ID — transient errors will corrupt state and orphan remote resources",
					})
				}
			}
			return true
		})

		return false
	})

	return violations
}

func isTypesStringValueEmpty(expr ast.Expr) bool {
	call, ok := expr.(*ast.CallExpr)
	if !ok || len(call.Args) != 1 {
		return false
	}
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || selector.Sel == nil || selector.Sel.Name != "StringValue" {
		return false
	}
	x, ok := selector.X.(*ast.Ident)
	if !ok || x.Name != "types" {
		return false
	}
	arg, ok := call.Args[0].(*ast.BasicLit)
	return ok && arg.Kind == token.STRING && arg.Value == `""`
}

// TestNoReadRemoveResourceOnTransientError scans Read functions for the
// anti-pattern where RemoveResource is called in a non-not-found error branch.
func TestNoReadRemoveResourceOnTransientError(t *testing.T) {
	for _, file := range resourceFiles(t) {
		violations := scanForReadRemoveResourceOnError(t, file)
		for _, v := range violations {
			t.Errorf("%s:%d: %s", file, v.line, v.msg)
		}
	}
}

func scanForReadRemoveResourceOnError(t *testing.T, filename string) []violation {
	t.Helper()
	fset, file := parseGoFile(t, filename)

	var violations []violation
	ast.Inspect(file, func(node ast.Node) bool {
		fn, ok := node.(*ast.FuncDecl)
		if !ok || fn.Name == nil || fn.Name.Name != "Read" || fn.Body == nil {
			return true
		}

		ast.Inspect(fn.Body, func(inner ast.Node) bool {
			ifStmt, ok := inner.(*ast.IfStmt)
			if !ok || ifStmt.Body == nil {
				return true
			}
			if !blockDirectlyContainsCall(ifStmt.Body, "tflog", "Warn") {
				return true
			}
			if conditionMentionsNotFound(ifStmt.Cond) {
				return true
			}
			if isLenZeroNotFoundCheck(ifStmt.Cond) {
				return true
			}
			if blockDirectlyContainsSelectorName(ifStmt.Body, "RemoveResource") {
				violations = append(violations, violation{
					line: fset.Position(ifStmt.Pos()).Line,
					msg:  "Read function calls RemoveResource after tflog.Warn on non-not-found error — transient failures will incorrectly remove state",
				})
			}
			return true
		})

		return false
	})

	return violations
}

func conditionMentionsNotFound(expr ast.Expr) bool {
	found := false
	ast.Inspect(expr, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.Ident:
			if n.Name == "ErrResourceNotFound" {
				found = true
				return false
			}
		case *ast.SelectorExpr:
			if n.Sel != nil && n.Sel.Name == "ErrResourceNotFound" {
				found = true
				return false
			}
		case *ast.CallExpr:
			if isSelectorCall(node, "", "isZStackNotFoundError") || isIdentCall(node, "isZStackNotFoundError") {
				found = true
				return false
			}
		}
		return true
	})
	return found
}

func isLenZeroNotFoundCheck(expr ast.Expr) bool {
	binaryExpr, ok := expr.(*ast.BinaryExpr)
	if !ok || binaryExpr.Op != token.EQL {
		return false
	}
	call, ok := binaryExpr.X.(*ast.CallExpr)
	if !ok {
		return false
	}
	ident, ok := call.Fun.(*ast.Ident)
	if !ok || ident.Name != "len" || len(call.Args) != 1 {
		return false
	}
	right, ok := binaryExpr.Y.(*ast.BasicLit)
	return ok && right.Kind == token.INT && right.Value == "0"
}

func blockDirectlyContainsCall(block *ast.BlockStmt, receiver, method string) bool {
	for _, stmt := range block.List {
		if exprStmt, ok := stmt.(*ast.ExprStmt); ok && isSelectorCall(exprStmt.X, receiver, method) {
			return true
		}
	}
	return false
}

func isSelectorCall(node ast.Node, receiver, method string) bool {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return false
	}
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || selector.Sel == nil || selector.Sel.Name != method {
		return false
	}
	if receiver == "" {
		return true
	}
	x, ok := selector.X.(*ast.Ident)
	return ok && x.Name == receiver
}

func isIdentCall(node ast.Node, name string) bool {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return false
	}
	ident, ok := call.Fun.(*ast.Ident)
	return ok && ident.Name == name
}

func blockDirectlyContainsSelectorName(block *ast.BlockStmt, method string) bool {
	for _, stmt := range block.List {
		exprStmt, ok := stmt.(*ast.ExprStmt)
		if !ok {
			continue
		}
		call, ok := exprStmt.X.(*ast.CallExpr)
		if !ok {
			continue
		}
		selector, ok := call.Fun.(*ast.SelectorExpr)
		if ok && selector.Sel != nil && selector.Sel.Name == method {
			return true
		}
	}
	return false
}

// TestNoDeleteEmptyUUIDGuard scans all resource files for the anti-pattern
// where a Delete function checks for empty UUID and silently skips deletion.
func TestNoDeleteEmptyUUIDGuard(t *testing.T) {
	for _, file := range resourceFiles(t) {
		violations := scanForDeleteEmptyUUIDGuard(t, file)
		for _, v := range violations {
			t.Errorf("%s:%d: %s", file, v.line, v.msg)
		}
	}
}

func scanForDeleteEmptyUUIDGuard(t *testing.T, filename string) []violation {
	t.Helper()
	fset, file := parseGoFile(t, filename)

	var violations []violation
	ast.Inspect(file, func(node ast.Node) bool {
		fn, ok := node.(*ast.FuncDecl)
		if !ok || fn.Name == nil || fn.Name.Name != "Delete" || fn.Body == nil {
			return true
		}

		ast.Inspect(fn.Body, func(inner ast.Node) bool {
			ifStmt, ok := inner.(*ast.IfStmt)
			if !ok {
				return true
			}
			if containsTypesStringValueEmpty(ifStmt.Cond) {
				violations = append(violations, violation{
					line: fset.Position(ifStmt.Pos()).Line,
					msg:  "Delete function has empty-UUID guard (types.StringValue variant) — masks state corruption bugs",
				})
			}
			return true
		})

		return false
	})

	return violations
}

func containsTypesStringValueEmpty(expr ast.Expr) bool {
	found := false
	ast.Inspect(expr, func(node ast.Node) bool {
		if isTypesStringValueEmptyExpr(node) {
			found = true
			return false
		}
		return true
	})
	return found
}

func isTypesStringValueEmptyExpr(node ast.Node) bool {
	expr, ok := node.(ast.Expr)
	if !ok {
		return false
	}
	return isTypesStringValueEmpty(expr)
}
