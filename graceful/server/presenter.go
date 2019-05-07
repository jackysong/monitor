package server

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/kisekivul/utils"
)

func init() {
	t := template.New("dirlist")
	t = t.Funcs(template.FuncMap{
		"tosize": utils.ToString,
		"split":  strings.Split,
		"concat": func(a, b string) string { return a + b },
	})
	var err error
	dirlistHtmlTempl, err = t.Parse(dirlistHtml)
	if utils.ErrorCheck(err) {
		os.Exit(1)
	}
}

type listDir struct {
	Path, Parent      string
	NumFiles, NumDirs int
	TotalSize         int64
	Archive           bool
	Files             []listFile
}

type listFile struct {
	Path, Name string
	Accessible bool
	IsDir      bool
	Size       int64
	Mtime      time.Time
}

type byName []listFile

func (s *Handler) dirlist(w http.ResponseWriter, r *http.Request, dir string) {
	path, _ := filepath.Rel(s.config.Directory, dir)
	parent := ""
	if path != "." {
		parent = "/" + filepath.Join(path, "..")
	}

	list := &listDir{
		Path:    path,
		Parent:  parent,
		Archive: s.config.Archivable,
		Files:   []listFile{},
	}

	//readnames and stat separately so a single failed
	//stat doesn't cause the directory listing to fail
	d, err := os.Open(dir)
	if utils.ErrorCheck(err) {
		utils.Return_Content_500(w, err.Error())
		return
	}
	names, err := d.Readdirnames(-1)
	if utils.ErrorCheck(err) {
		utils.Return_Content_500(w, err.Error())
		return
	}

	for _, n := range names {
		if n == ".DS_Store" {
			continue //Nope.
		}
		lf := listFile{
			Name: n,
			Path: "/" + filepath.Join(path, n),
		}
		//attempt to stat
		if f, err := os.Stat(filepath.Join(dir, n)); err == nil {
			lf.Accessible = true
			var size int64
			if f.IsDir() {
				n += "/"
				list.NumDirs++
			} else {
				list.NumFiles++
				size = f.Size()
				list.TotalSize += size
			}
			lf.IsDir = f.IsDir()
			lf.Size = size
			lf.Mtime = f.ModTime()
		}

		list.Files = append(list.Files, lf)
	}

	sort.Slice(list.Files, func(p, q int) bool {
		return list.Files[p].Name < list.Files[q].Name
	})

	accepts := strings.Split(r.Header.Get("Accept"), ",")
	buff := &bytes.Buffer{}
	contype := ""
	for _, accept := range accepts {
		typeencoding := strings.SplitN(accept, "/", 2)
		if len(typeencoding) != 2 {
			continue
		}
		switch typeencoding[1] {
		case "json":
			b, _ := json.MarshalIndent(list, "", "  ")
			buff.Write(b)
		case "xml":
			b, _ := xml.MarshalIndent(list, "", "  ")
			buff.Write(b)
		case "html":
			dirlistHtmlTempl.Execute(buff, list)
		default:
			continue
		}
		contype = accept
		break
	}

	if contype == "" {
		for _, f := range list.Files {
			buff.WriteString(f.Name + "\n")
		}
		contype = "text/plain"
	}

	w.Header().Set("Content-Type", contype)
	w.WriteHeader(200)
	w.Write(buff.Bytes())
}

var dirlistHtmlTempl *template.Template

var dirlistHtml = `
<html>
	<head>
		<title>{{ .Path }}</title>
		<style>
			html,body {
				height:100%;
				width:100%;
				font-family: Courier, monospace;
			}
			a {
				text-decoration: none;
			}
			table {
				margin: 5%;
			}
			.path {
				text-style: underline;
			}
			.name {
				text-align: right;
				padding-right: 30px;
			}
			.name a {
				word-wrap:break-word;
				display: inline-block;
				width: 300px;
			}
			.size {
				text-align: left;
			}
			.archive {
				font-size: 0.8em;
			}

		</style>
	</head>
	<body>
		<table>
			<tr>
				<th class="name">Name</th>
				<th class="size">Size</th>
			</tr>
			<tr class="file item">
				<td class="name"><a href="/{{ .Path }}">.</a></td>
				<td class="size">-</td>
			</tr>
			{{if ne .Parent ""}}<tr class="file item">
				<td class="name"><a href="{{ .Parent }}">..</a></td>
				<td class="size">-</td>
			</tr>{{end}}
			{{range .Files}}<tr class="file item">
				<td class="name">
					{{if .Accessible}}
						<a href="{{ .Path }}{{if .IsDir}}/{{end}}">{{ .Name }}</a>
					{{else}}
						{{ .Name }}
					{{end}}
				</td>
				<td class="size" alt="{{ .Size }} bytes">
					{{if .IsDir}}-{{else if not .Accessible}}-{{else}}{{ tosize .Size }}{{end}}
				</td>
			</tr>{{end}}
			{{if .NumFiles}}<tr class="files">
				<th class="name">
					{{.NumFiles}} file{{if ne .NumFiles 1}}s{{end}}
				</th>
				<th class="size" alt="{{ .TotalSize }} bytes">
					{{ tosize .TotalSize }}
				</th>
			</tr>{{end}}
			{{if .NumDirs}}<tr class="files">
				<th class="name">
					{{.NumDirs}} dir{{if ne .NumDirs 1}}s{{end}}
				</th>
				<th>
				</th>
			</tr>{{end}}
			{{if .Archive}}<tr class="archive">
				<th class="name">
					download all as
				</th>
				<th>
					<a href="/{{ .Path }}.zip">zip</a>,
					<a href="/{{ .Path }}.tar">tar</a>,
					<a href="/{{ .Path }}.tar.gz">tar.gz</a>
				</th>
			</tr>{{end}}
		</table>
	</body>
</html>
`
