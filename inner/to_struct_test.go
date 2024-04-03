package inner

import (
	"fmt"
	"strings"
	"testing"
)

type JSONData struct {
	Name string `json:"name"`
	Code int    `json:"code"`
	List []struct {
		Day int `json:"day"`
	} `json:"list"`
}

func TestStruct(t *testing.T) {
	s := `{"name":"小芳","code":404,"list":[{"day":1},{"day":2}]}`
	scan := NewScanner(strings.NewReader(s))
	p := NewParseStruct(scan)
	// js := p.Json()
	// fmt.Println(js.Value)
	p1 := &JSONData{}
	x := p.Json(p1)
	fmt.Println(x, p1)
}
