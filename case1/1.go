package main

import (
	"fmt"
)

type JSONData struct {
	Name string   `json:"name"`
	Code int64    `json:"code"`
	List []string `json:"list"`
}

type Scanner struct {
	src   string
	index int
}

func NewScanner(src string) *Scanner {
	return &Scanner{src: src}
}

func (s *Scanner) scan() string {
	for s.index < len(s.src) && s.src[s.index] <= ' ' {
		s.index++
	}
	return s.src[s.index:]
}

func (s *Scanner) scanString() string {
	if s.index >= len(s.src) || s.src[s.index] != '"' {
		return ""
	}
	s.index++
	start := s.index
	for s.index < len(s.src) && s.src[s.index] != '"' {
		s.index++
	}
	return s.src[start:s.index]
}

func (s *Scanner) scanNumber() int64 {
	start := s.index
	for s.index < len(s.src) && (s.src[s.index] == '-' || (s.src[s.index] >= '0' && s.src[s.index] <= '9')) {
		s.index++
	}
	num := s.src[start:s.index]
	var result int64
	fmt.Sscanf(num, "%d", &result)
	return result
}

func (s *Scanner) scanList() []string {
	if s.index >= len(s.src) || s.src[s.index] != '[' {
		return nil
	}
	s.index++
	var list []string
	for {
		val := s.scanString()
		if val == "" {
			break
		}
		list = append(list, val)
		s.scan()
		if s.index >= len(s.src) || s.src[s.index] != ',' {
			break
		}
		s.index++
	}
	if s.index < len(s.src) && s.src[s.index] == ']' {
		s.index++
	}
	return list
}

func (s *Scanner) scanStruct(data interface{}) {
	s.scan() // remove leading spaces
	if s.index >= len(s.src) || s.src[s.index] != '{' {
		return
	}
	s.index++
	for {
		key := s.scanString()
		if key == "" {
			break
		}
		s.scan()
		if s.index >= len(s.src) || s.src[s.index] != ':' {
			break
		}
		s.index++
		s.scan()
		switch t := data.(type) {
		case *JSONData:
			switch key {
			case "name":
				t.Name = s.scanString()
			case "code":
				t.Code = s.scanNumber()
			case "list":
				t.List = s.scanList()
			}
		}
		s.scan()
		if s.index >= len(s.src) || s.src[s.index] != ',' {
			break
		}
		s.index++
	}
	if s.index < len(s.src) && s.src[s.index] == '}' {
		s.index++
	}
}

func main() {
	s := `{"name":"小芳","code":404,"list":["day","1"]}`
	scan := NewScanner(s)
	data := &JSONData{}
	scan.scanStruct(data)
	fmt.Println(data)
}
