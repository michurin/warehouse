package templatetree_test

import (
	"fmt"
	"text/template"
	"text/template/parse"
)

func R(pfx string, nodes ...parse.Node) {
	for _, it := range nodes {
		switch n := it.(type) {
		case *parse.TextNode:
			fmt.Printf(pfx+"text: %q\n", n.Text)
		case *parse.ActionNode:
			fmt.Println(pfx+"action:", n)
			fmt.Println(pfx + "a: decls:")
			for _, d := range n.Pipe.Decl {
				fmt.Println(pfx+"a: decls:", d)
				for _, i := range d.Ident {
					fmt.Println(pfx + i)
				}
			}
			fmt.Println(pfx + "a: cmds:")
			for _, d := range n.Pipe.Cmds {
				R(pfx+"  ", d)
			}
		case *parse.CommandNode:
			fmt.Println(pfx+"cmd:", n, n.Position())
			for _, a := range n.Args {
				R(pfx+"  ", a)
			}
		case *parse.FieldNode:
			fmt.Println(pfx+"field:", n)
			for _, d := range n.Ident {
				fmt.Println(pfx+"f: ind:", d)
			}
		case *parse.NumberNode:
			fmt.Println(pfx+"num:", n.Text)
		case *parse.IdentifierNode:
			fmt.Println(pfx+"identifier:", n.Ident)
		default:
			fmt.Printf(pfx+"UNKNOWN NODE: %T\n", n)
		}
	}
}

func Example() {
	tt := template.Must(template.New("x").Parse(`a {{ .X.Y 99 | js }} b`))
	R("", tt.Root.Nodes...)
	// output:
	// text: "a "
	// action: {{.X.Y 99 | js}}
	// a: decls:
	// a: cmds:
	//   cmd: .X.Y 99 5
	//     field: .X.Y
	//     f: ind: X
	//     f: ind: Y
	//     num: 99
	//   cmd: js 15
	//     identifier: js
	// text: " b"
}
