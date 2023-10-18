package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"text/template"
)

type tdata struct {
	Title string
	Dirs  []tdir
	Files []tfile
}

type tdir struct {
	Name string
	Path string
}

type tfile struct {
	Name string
	Path string
}

func response(pth, show string) ([]byte, error) {
	if show != "" {
		return ioutil.ReadFile(show) // vulnerable!
	}

	tmpl, err := template.New("x").Parse(`
<style>
body { background-color: #fff; color: #777; }
* { font-family: "Gill Sans", sans-serif; }
a { text-decoration: none; }
.box {
    display: flex;
	flex-wrap: wrap;
}
.box > div:hover {
	background-color: #eee;
}
.title {
	text-align: center;
	font-size: xx-small;
}
img {
	width: 300px;
}
</style>
<h1>{{.Title}}</h1>
<ul>
	{{range $val := .Dirs}}
		<li>
			<a href="?path={{$val.Path}}">{{$val.Name}}</a>
		</li>
	{{end}}
</ul>
<hr>
<div class="box">
{{range $val := .Files}}
	<div>
		<div><a href="?show={{$val.Path}}" target="_blank"><img src="?show={{$val.Path}}" alt="{{$val.Name}}"></a></div>
		<div class="title">{{$val.Name}}</div>
	</div>
{{end}}
</div>
`)
	if err != nil {
		return nil, err
	}

	tdata := tdata{Title: pth}

	dir, err := ioutil.ReadDir(pth)
	if err != nil {
		return nil, err
	}
	for _, e := range dir {
		n := e.Name()
		p := path.Join(pth, n)
		if e.IsDir() {
			tdata.Dirs = append(tdata.Dirs, tdir{Name: n, Path: p})
		} else {
			if strings.HasSuffix(n, ".png") {
				tdata.Files = append(tdata.Files, tfile{Name: n, Path: p})
			}
		}
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, tdata)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func defStr(x, def string) string {
	if x == "" {
		return def
	}
	return x
}

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	err = os.Chdir(home)
	if err != nil {
		panic(err)
	}
	http.ListenAndServe(":9099", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		qv := r.URL.Query()
		body, err := response(defStr(qv.Get("path"), "."), qv.Get("show"))
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(body)
	}))
}
