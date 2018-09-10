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

var dir = flag.String("dir", ".", "dir to pack")
var packageName = flag.String("package", "main", "name of the data package")
var variableName = flag.String("variable", "pack", "name of the data variable")
var out = flag.String("output", "/dev/stdout", "output file")
var exclude = flag.String("exclude", "", "file prefixes to exclude from packing")

const tmpl = `package {{.PackageName}}
var {{.VariableName}} = map[string][]byte{
  {{ range $key, $value := .Pack.Entries }}"{{$key}}": {
    {{ range $value }}{{ . }},{{ end }}
  },
{{end}}}
`

type Pack struct {
	Entries map[string][]byte
}

func NewPack(dir string, excludeList []string) (*Pack, error) {
	pack := &Pack{make(map[string][]byte)}
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
			pack.Entries[path] = bs
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return pack, nil
}

func (pack *Pack) ToCode(packageName, variableName string, output io.Writer) error {

	t := template.Must(template.New("").Parse(tmpl))
	data := struct {
		PackageName  string
		VariableName string
		Pack         *Pack
	}{
		packageName,
		variableName,
		pack,
	}
	if err := t.Execute(output, data); err != nil {
		return err
	}
	return nil
}

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
	err = pack.ToCode(*packageName, *variableName, outFile)
	if err != nil {
		log.Fatal(err)
	}
}
