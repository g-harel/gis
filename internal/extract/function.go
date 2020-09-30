package extract

import (
	"go/ast"
	"go/token"

	"github.com/g-harel/gothrough/internal/types"
)

// newFunctionVisitor creates a visitor that collects Funcs into the target array.
func newFunctionVisitor(handler func(Location, types.Function)) visitFunc {
	return func(filepath string, n ast.Node, fset *token.FileSet) bool {
		if n == nil {
			return true
		}

		if funcDeclaration, ok := n.(*ast.FuncDecl); ok {
			if funcDeclaration.Name.IsExported() {
				handler(
					getLocation(filepath),
					types.Function{
						Name:         funcDeclaration.Name.String(),
						Docs:         types.Docs{Text: funcDeclaration.Doc.Text()},
						Arguments:    collectFields(funcDeclaration.Type.Params),
						ReturnValues: collectFields(funcDeclaration.Type.Results),
					},
				)
			}
		}
		return false
	}
}
