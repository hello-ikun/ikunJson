package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/hello-ikun/ikunJson/inner"
	"io"
	"os"
	"path/filepath"
)

func main() {
	s := []string{"twitter", "code", "example", "sample"}
	for _, v := range s[:1] {
		fp := fixture(v)
		scan := inner.NewScanner(fp)
		p := inner.NewParse(scan)
		js := p.Json()
		fmt.Println(inner.Tokens[js.Type], js.Value, js.Value)
		fmt.Println()
	}
}

// fuxture returns a *bytes.Reader for the contents of path.
func fixture(path string) *bytes.Reader {
	f, err := os.Open(filepath.Join("testdata", path+".json.gz"))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		panic(err)
	}
	buf, err := io.ReadAll(gz)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(buf)
}
