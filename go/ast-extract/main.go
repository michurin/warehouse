package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/tools/go/packages"
)

type outputDto struct { // TODO rename fields according neovim quickfix
	File    string `json:"filename"`
	Line    int    `json:"lnum"`
	Col     int    `json:"col"`
	Message string `json:"text"`
}

func findMethods(root, targetTypeName string) []outputDto {
	result := []outputDto{}
	walkPackage(root, func(node ast.Node, fset *token.FileSet, err error) error {
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
		return nil
	})
	return result
}

func debug(pfx string, n ast.Node) {
	if n == nil {
		return
	}
	ast.Inspect(n, func(n ast.Node) bool {
		fmt.Printf("%[1]s\033[1;32m-->\033[0m [%[2]T][%[2]v]\n", pfx, n)
		ast.Inspect(n, func(nx ast.Node) bool {
			if nx == n { // skip root
				return true
			}
			debug(pfx+"  ", nx)
			return false
		})
		return false
	})
}

func debugFile(path string) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	debug("", file)
}

func findFirstIdent(n ast.Node, name string) *ast.Ident {
	res := (*ast.Ident)(nil)
	ast.Inspect(n, func(n ast.Node) bool {
		if idn, ok := n.(*ast.Ident); ok && idn.Name == name {
			res = idn
		}
		return res == nil
	})
	return res
}

func walkGoFiles(root string, f func(ast.Node, *token.FileSet, error) error) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			panic(path + " " + err.Error())
		}
		return f(file, fset, err)
	})
}

func walkPackage(root string, f func(ast.Node, *token.FileSet, error) error) error {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedSyntax | packages.NeedFiles | packages.NeedTypes,
		Dir:  root,
	}, "./...")
	if err != nil {
		panic(err)
	}
	for _, pkg := range pkgs {
		fset := pkg.Fset
		for _, node := range pkg.Syntax {
			err := f(node, fset, nil)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func findConstructors(root, targetTypeName string) []outputDto { // TODO functions only; consider interface methods? consider type definitions?
	result := []outputDto(nil)
	err := walkGoFiles(root, func(file ast.Node, fset *token.FileSet, err error) error {
		ast.Inspect(file, func(n ast.Node) bool {
			if n != nil {
				if fn, ok := n.(*ast.FuncDecl); ok && fn.Type != nil && fn.Type.Results != nil {
					if e := findFirstIdent(fn.Type.Results, targetTypeName); e != nil {
						pos := fset.Position(fn.Name.Pos())
						result = append(result, outputDto{
							File:    pos.Filename,
							Line:    pos.Line,
							Col:     pos.Column,
							Message: fn.Name.Name + " -> " + e.Name, // TODO add docs here?
						})
						return false
					}
				}
			}
			return true
		})
		return nil
	})
	if err != nil {
		panic(err)
	}
	return result
}

func findInFunctionName() {} // TODO

func findAllVarDefinitions() { // TODO?
	debugFile("example/v.go")
}

func findInString(root string) []outputDto {
	result := []outputDto(nil)
	err := walkGoFiles(root, func(file ast.Node, fset *token.FileSet, err error) error {
		ast.Inspect(file, func(n ast.Node) bool {
			if a, ok := n.(*ast.BasicLit); ok {
				if a.Kind != token.STRING {
					return false
				}
				x, err := strconv.Unquote(a.Value)
				if err != nil {
					panic(err) // impossible for token.STRING
				}
				pos := fset.Position(n.Pos())
				c := pos.Column
				for i, s := range strings.Split(x, "\n") {
					s = strings.TrimSpace(s)
					if len(s) == 0 {
						continue
					}
					result = append(result, outputDto{
						File:    pos.Filename,
						Line:    pos.Line + i,
						Col:     c,
						Message: s,
					})
					// fmt.Println(pos.Filename, pos.Line+i, c, s)
					c = 1
				}
			}
			return true
		})
		return nil
	})
	if err != nil {
		panic(err)
	}
	return result
}

func encode(r []outputDto) {
	err := json.NewEncoder(os.Stdout).Encode(r)
	if err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "USAGE %s methods|strings ...", os.Args[0])
	}
	switch os.Args[1] {
	case "methods":
		dir := os.Args[2]
		targetTypeName := os.Args[3]
		encode(findMethods(dir, targetTypeName))
	case "strings":
		dir := os.Args[2]
		encode(findInString(dir))
	case "constructors":
		dir := os.Args[2]
		targetTypeName := os.Args[3]
		encode(findConstructors(dir, targetTypeName))
	default:
		panic("unknown command")
	}
}
