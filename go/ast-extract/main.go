package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/token"
	"os"

	"golang.org/x/tools/go/packages"
)

type outputDto struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	Col     int    `json:"col"`
	Message string `json:"message"`
}

func findMethods(node ast.Node, fset *token.FileSet, targetTypeName string) []outputDto {
	result := []outputDto(nil)
	ast.Inspect(node, func(n ast.Node) bool {
		fd, ok := n.(*ast.FuncDecl)
		if ok {
			if fd.Recv == nil {
				return false
			}
			if len(fd.Recv.List) != 1 {
				return false
			}
			e := (*ast.Ident)(nil)
			star := ""
			switch q := fd.Recv.List[0].Type.(type) {
			case *ast.Ident:
				e = q
			case *ast.StarExpr:
				e = q.X.(*ast.Ident)
				star = "*"
			}
			if e != nil && e.Name == targetTypeName {
				pos := fset.Position(fd.Name.Pos())
				result = append(result, outputDto{
					File:    pos.Filename,
					Line:    pos.Line,
					Col:     pos.Column,
					Message: star + e.Name + "." + fd.Name.Name, // TODO add docs here?
				})
			}
			return false
		}
		return true
	})
	return result
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "USAGE %s directory type_name", os.Args[0])
	}
	dir := os.Args[1]
	targetTypeName := os.Args[2]
	cfg := &packages.Config{
		Mode: packages.NeedSyntax | packages.NeedFiles | packages.NeedTypes,
		Dir:  dir,
	}

	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		panic(err)
	}

	result := []outputDto{} // empty, not null
	for _, pkg := range pkgs {
		fs := pkg.Fset
		for _, node := range pkg.Syntax {
			result = append(result, findMethods(node, fs, targetTypeName)...)
		}
	}

	err = json.NewEncoder(os.Stdout).Encode(result)
	if err != nil {
		panic(err)
	}
}
