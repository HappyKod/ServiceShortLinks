// Package myanalyzers обнаружение os.Exit в main функции
package myanalyzers

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Doc описание при обнаружении
const Doc = "обнаружен os.Exit в main функции"

// ExitAnalyzer запрещено использовать прямой вызов os.Exit в функции main основного пакет
var ExitAnalyzer = &analysis.Analyzer{
	Name: "exit",
	Doc:  Doc,
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		fileName := pass.Fset.Position(file.Pos()).Filename
		if !strings.HasSuffix(fileName, ".go") {
			continue
		}

		ast.Inspect(file, func(node ast.Node) bool {
			f, ok := node.(*ast.FuncDecl)
			if !ok {
				return true
			}
			if f.Name.Name == "main" {
				return true
			}
			for _, s := range f.Body.List {
				expr, exprOk := s.(*ast.ExprStmt)
				if !exprOk {
					return true
				}
				call, callOk := expr.X.(*ast.CallExpr)
				if !callOk {
					return true
				}
				selector, selectorOk := call.Fun.(*ast.SelectorExpr)
				if !selectorOk {
					return true
				}
				i := selector.X.(*ast.Ident)
				if i.Name == "os" && selector.Sel.Name == "Exit" {
					pass.Reportf(selector.Pos(), Doc)
				}
			}
			return true
		})
	}
	return nil, nil
}
