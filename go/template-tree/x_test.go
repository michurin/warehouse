package templatetree_test

import (
	"fmt"
	"text/template"
	"text/template/parse"
)

func RA[T parse.Node](pfx string, a []T) {
	for _, n := range a {
		R(pfx, n)
	}
}

func R(pfx string, nd parse.Node) {
	switch n := nd.(type) {
	case *parse.TextNode:
		fmt.Printf(pfx+"text: %q\n", n.Text)
	case *parse.ActionNode:
		fmt.Println(pfx+"action:", n)
		R("  ", n.Pipe)
	case *parse.CommandNode:
		fmt.Println(pfx+"cmd:", n, n.Position())
		RA(pfx+"  ", n.Args)
	case *parse.FieldNode:
		fmt.Println(pfx+"field:", n)
		for _, d := range n.Ident {
			fmt.Println(pfx+"f: ind:", d)
		}
	case *parse.VariableNode:
		fmt.Println(pfx+"variable:", n)
		for _, d := range n.Ident {
			fmt.Println(pfx+"v: ind:", d)
		}
	case *parse.NumberNode:
		fmt.Println(pfx+"num:", n.Text)
	case *parse.IdentifierNode:
		fmt.Println(pfx+"identifier:", n.Ident)
	case *parse.RangeNode:
		fmt.Println(pfx+"range: type=", n.NodeType)
		fmt.Println(pfx + "range: pipe")
		R(pfx+"  ", n.Pipe)
		fmt.Println(pfx + "range: list")
		R(pfx+"  ", n.List)
		fmt.Println(pfx + "range: else")
		R(pfx+"  ", n.ElseList)
	case *parse.PipeNode:
		if n != nil {
			fmt.Println(pfx + "pipe: decl")
			RA(pfx+"  ", n.Decl)
			fmt.Println(pfx + "pipe: cmds")
			RA(pfx+"  ", n.Cmds)
		}
	case *parse.ListNode:
		fmt.Println(pfx+"list", n)
		if n != nil { // TODO check it everywhere?
			RA(pfx+"  ", n.Nodes)
		}
	case *parse.BoolNode:
		fmt.Println(pfx+"bool:", n.True)
	case *parse.NilNode:
		fmt.Println(pfx + "nil node")
	case *parse.TemplateNode:
		fmt.Println(pfx+"template: type=", n.NodeType, "name=", n.Name)
		fmt.Println(pfx + "template: pipe")
		R(pfx+"  ", n.Pipe)
	case *parse.IfNode:
		fmt.Println(pfx+"if: type=", n.NodeType)
		fmt.Println(pfx + "if: pipe")
		R(pfx+"  ", n.Pipe)
		fmt.Println(pfx + "if: list")
		R(pfx+"  ", n.List)
		fmt.Println(pfx + "if: else")
		R(pfx+"  ", n.ElseList)
	case *parse.WithNode:
		fmt.Println(pfx+"with: type=", n.NodeType)
		fmt.Println(pfx + "with: pipe")
		R(pfx+"  ", n.Pipe)
		fmt.Println(pfx + "with: list")
		R(pfx+"  ", n.List)
		fmt.Println(pfx + "with: else")
		R(pfx+"  ", n.ElseList)
	case *parse.DotNode:
		fmt.Println(pfx + "dot")
	case *parse.ContinueNode:
		fmt.Println(pfx + "continue")
	case *parse.BreakNode:
		fmt.Println(pfx + "break")
	case *parse.StringNode:
		fmt.Printf(pfx+"quoted text: %q (%s)\n", n.Text, n.Quoted)
	case *parse.ChainNode: // TODO
		fmt.Println(pfx+"NOT IMPLEMENTED", n.String())
	case *parse.CommentNode: // never appears?
		fmt.Println(pfx+"NOT IMPLEMENTED", n.String())
	case *parse.BranchNode: // umbrella only? if/with/range
		fmt.Println(pfx+"WON'T BE IMPLEMENTED", n.String())
	default:
		fmt.Printf(pfx+"UNKNOWN NODE: %T\n", n)
	}
}

func Example() {
	tt := template.Must(template.New("x").Parse(
		`a {{ true | js }} {{ nil }} {{ .X.Y 99 | js -}}
		{{ range $v := .A }}{{ $v | js }}{{ continue }}{{ break }}{{ end -}}
		{{ if true }}TRUE{{ else if false }}FALSE{{ else }}ELSE{{ end -}}
		{{/* comment, won't show up in AST */ -}}
		{{ define "T1" }}T1{{ end -}}
		{{ block "T2" .X.T2 }}T2{{ end -}}
		{{ template "T1" -}}
		{{ template "T2" -}}
		{{ with true }}{{ . }}{{ end -}}
		{{ "x" }}`,
	))
	R("", tt.Root)
	// output:
	// list a {{true | js}} {{nil}} {{.X.Y 99 | js}}{{range $v := .A}}{{$v | js}}{{continue}}{{break}}{{end}}{{if true}}TRUE{{else}}{{if false}}FALSE{{else}}ELSE{{end}}{{end}}{{template "T2" .X.T2}}{{template "T1"}}{{template "T2"}}{{with true}}{{.}}{{end}}{{"x"}}
	//   text: "a "
	//   action: {{true | js}}
	//   pipe: decl
	//   pipe: cmds
	//     cmd: true 5
	//       bool: true
	//     cmd: js 12
	//       identifier: js
	//   text: " "
	//   action: {{nil}}
	//   pipe: decl
	//   pipe: cmds
	//     cmd: nil 21
	//       nil node
	//   text: " "
	//   action: {{.X.Y 99 | js}}
	//   pipe: decl
	//   pipe: cmds
	//     cmd: .X.Y 99 31
	//       field: .X.Y
	//       f: ind: X
	//       f: ind: Y
	//       num: 99
	//     cmd: js 41
	//       identifier: js
	//   range: type= 15
	//   range: pipe
	//     pipe: decl
	//       variable: $v
	//       v: ind: $v
	//     pipe: cmds
	//       cmd: .A 65
	//         field: .A
	//         f: ind: A
	//   range: list
	//     list {{$v | js}}{{continue}}{{break}}
	//       action: {{$v | js}}
	//   pipe: decl
	//   pipe: cmds
	//     cmd: $v 73
	//       variable: $v
	//       v: ind: $v
	//     cmd: js 78
	//       identifier: js
	//       continue
	//       break
	//   range: else
	//     list <nil>
	//   if: type= 10
	//   if: pipe
	//     pipe: decl
	//     pipe: cmds
	//       cmd: true 127
	//         bool: true
	//   if: list
	//     list TRUE
	//       text: "TRUE"
	//   if: else
	//     list {{if false}}FALSE{{else}}ELSE{{end}}
	//       if: type= 10
	//       if: pipe
	//         pipe: decl
	//         pipe: cmds
	//           cmd: false 149
	//             bool: false
	//       if: list
	//         list FALSE
	//           text: "FALSE"
	//       if: else
	//         list ELSE
	//           text: "ELSE"
	//   template: type= 17 name= T2
	//   template: pipe
	//     pipe: decl
	//     pipe: cmds
	//       cmd: .X.T2 279
	//         field: .X.T2
	//         f: ind: X
	//         f: ind: T2
	//   template: type= 17 name= T1
	//   template: pipe
	//   template: type= 17 name= T2
	//   template: pipe
	//   with: type= 19
	//   with: pipe
	//     pipe: decl
	//     pipe: cmds
	//       cmd: true 356
	//         bool: true
	//   with: list
	//     list {{.}}
	//       action: {{.}}
	//   pipe: decl
	//   pipe: cmds
	//     cmd: . 366
	//       dot
	//   with: else
	//     list <nil>
	//   action: {{"x"}}
	//   pipe: decl
	//   pipe: cmds
	//     cmd: "x" 386
	//       quoted text: "x" ("x")
}
