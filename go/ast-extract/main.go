package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type X struct{}

func (X) One() {}

func (i X) Onex() {}

func (i *X) Onexx() {}

func (_ X) Onexxx() {}

func main() {
	fset := token.NewFileSet()

	node, err := parser.ParseFile(fset, "main.go", nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	ast.Print(fset, node)

	ast.Inspect(node, func(n ast.Node) bool {
		fd, ok := n.(*ast.FuncDecl)
		if ok {
			if fd.Recv == nil {
				return false
			}
			if len(fd.Recv.List) != 1 {
				return false
			}
			switch q := fd.Recv.List[0].Type.(type) {
			case *ast.Ident:
				fmt.Printf("%T %s %d %s\n", n, fd.Name, fd.Recv.NumFields(), q.Name)
			case *ast.StarExpr:
				fmt.Printf("%T %s %d %s\n", n, fd.Name, fd.Recv.NumFields(), q.X.(*ast.Ident).Name)
			}
			return false
		}
		return true
	})
}
