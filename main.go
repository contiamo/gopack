package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	dir          = flag.String("dir", ".", "dir to pack")
	packageName  = flag.String("package", "main", "name of the data package")
	variableName = flag.String("variable", "pack", "name of the data variable")
	out          = flag.String("output", "/dev/stdout", "output file")
	exclude      = flag.String("exclude", "", "file prefixes to exclude from packing")
	withHTTP     = flag.Bool("with-http-handler", false, "include convenient http handler")

	tmpl *template.Template
)

// this is the template used to render the resulting go source file
const tmplStr = `package {{.PackageName}}
{{ if .WithHTTPHandler }}
import "github.com/contiamo/gopack/staticserver"
{{ end }}
var {{.VariableName}} = map[string][]byte{
  {{ range $key, $value := .Pack }}"{{$key}}": {
    {{ range $value }}{{ . }},{{ end }}
  },
{{end}}}
{{ if .WithHTTPHandler }}
var {{.VariableName}}Handler = staticserver.New({{.VariableName}})
{{ end }}
`

// we parse the template once at program startup
func init() {
	tmpl = template.Must(template.New("").Parse(tmplStr))
}

// Pack represents a set of files
type Pack map[string][]byte

// NewPack creates a new pack by iterating a directory
func NewPack(dir string, excludeList []string) (Pack, error) {
	pack := make(Pack)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.Mode().IsRegular() {
			bs, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			for _, prefix := range excludeList {
				if strings.HasPrefix(path, prefix) {
					return nil
				}
			}
			pack[path] = bs
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return pack, nil
}

// ToCode generates a go source file containing a map representing the files of the pack
func (pack Pack) ToCode(packageName, variableName string, withHTTPHandler bool, output io.Writer) error {
	type templateData struct {
		PackageName     string
		VariableName    string
		WithHTTPHandler bool
		Pack            Pack
	}
	data := &templateData{
		packageName,
		variableName,
		withHTTPHandler,
		pack,
	}
	if err := tmpl.Execute(output, data); err != nil {
		return err
	}
	return nil
}

// main() kicks of the programs logic
func main() {
	flag.Parse()
	excludeList := strings.Split(*exclude, ",")
	if len(excludeList) == 1 && excludeList[0] == "" {
		excludeList = nil
	}
	pack, err := NewPack(*dir, excludeList)
	if err != nil {
		log.Fatal(err)
	}
	outFile, err := os.Create(*out)
	err = pack.ToCode(*packageName, *variableName, *withHTTP, outFile)
	if err != nil {
		log.Fatal(err)
	}
}
