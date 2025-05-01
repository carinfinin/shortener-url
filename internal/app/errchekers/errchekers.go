package errchekers

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

// ErrCheckAnalyzer анализатор проверки в пакете main вызовов функции os.Exit()
var ErrCheckAnalyzer = &analysis.Analyzer{
	Name: "errcheck",
	Doc:  "check for unchecked errors",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {

	if pass.Pkg.Name() != "main" {
		return nil, nil
	}

	for _, file := range pass.Files {

		ast.Inspect(file, func(node ast.Node) bool {

			if c, ok := node.(*ast.FuncDecl); ok && c.Name.Name == "main" {

				ast.Inspect(c, func(n ast.Node) bool {

					if call, ok := n.(*ast.CallExpr); ok {
						if s, ok := call.Fun.(*ast.SelectorExpr); ok {
							if exp, ok := s.X.(*ast.Ident); ok {
								if exp.Name == "os" && s.Sel.Name == "Exit" {
									pass.Reportf(call.Pos(), "прямой вызов os.Exit() в функции main() запрещен")
								}
							}
						}
					}
					return true
				})

			}
			return true
		})

	}
	return nil, nil
}
