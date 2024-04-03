package inner

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"unicode"
)

// Token 表示令牌的类型
type Token uint8

const (
	// object部分
	ObjectStart Token = iota // {
	ObjectEnd                // }
	Map                      // 用于标识【映射类】go中是map
	// array部分
	ArrayStart // [
	ArrayEnd   // ]
	Array      // 用于标识【映射类】go中是slice array
	// 基本的数据信息
	Number // 数字类型
	String // 字符串类型
	// 布尔类型和空值null
	Bool // 布尔类型
	Null // null
	// 分割符号
	SeparatorColon // :
	SeparatorComma // ,
	// 特殊标识信息
	EOF    // 文件结束符号
	FAILED // 存在错误标志
	Struct // 结构体类型
)

var ErrUnknow = errors.New("unknown error")

// Tokens 将 Token 映射到其字符串表示形式
var Tokens = map[Token]string{
	ObjectStart: "{",
	ObjectEnd:   "}",
	Map:         "Map类型",

	ArrayStart: "[",
	ArrayEnd:   "]",
	Array:      "Array类型",

	Number: "NUMBER",
	String: "STRING",

	Bool: "BOOL",
	Null: "null",

	SeparatorColon: ":",
	SeparatorComma: ",",

	EOF:    "END OF FILE",
	FAILED: "failed",
}

// Position 表示位置信息
type Position struct {
	Line   int // 行
	Column int // 列
}

// TokenJson 返回数据信息 包括 类型 数据  位置信息
type TokenJson struct {
	Type  Token       // token类型
	Value interface{} // 存储的数值
	Pos   *Position   // 位置信息
}

// 新建tokenJson实例
func tokenJson(token Token, pos *Position) *TokenJson {
	return &TokenJson{token, Tokens[token], pos}
}

// 新建tokenJson实例
func NewTokenJson(token Token, value interface{}, pos *Position) *TokenJson {
	return &TokenJson{token, value, pos}
}

// 开始数据扫描 需要进行 位置信息和阅读指针
type Scanner struct {
	Pos    *Position
	Reader *bufio.Reader
}

// 新建Scanner实例
func NewScanner(reader io.Reader) *Scanner {
	return &Scanner{Pos: &Position{1, 0}, Reader: bufio.NewReader(reader)}
}
func (s *Scanner) Scan() *TokenJson {
	for {
		r, _, err := s.Reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return tokenJson(EOF, s.Pos)
			}
			return tokenJson(FAILED, s.Pos)
		}
		s.Pos.Column++

		switch r {
		case '\n':
			s.updatePos()
		case ':':
			return tokenJson(SeparatorColon, s.Pos)
		case ',':
			return tokenJson(SeparatorComma, s.Pos)
		case '{':
			return tokenJson(ObjectStart, s.Pos)
		case '}':
			return tokenJson(ObjectEnd, s.Pos)
		case '[':
			return tokenJson(ArrayStart, s.Pos)
		case ']':
			return tokenJson(ArrayEnd, s.Pos)
		case 'f', 't', 'n':
			s.backUp()
			return s.BoolOrNone()
		case '"':
			return s.String()
		default:
			if unicode.IsSpace(r) {
				continue
			} else if unicode.IsDigit(r) || r == '-' {
				s.backUp()
				return s.Number()
			} else {
				return tokenJson(FAILED, s.Pos)
			}
		}
	}
}

// 更新pos信息
func (s *Scanner) updatePos() {
	s.Pos.Line++
	s.Pos.Column = 0
}

// 进行pos回退
func (s *Scanner) backUp() error {
	if err := s.Reader.UnreadRune(); err != nil {
		return err
	}
	s.Pos.Column--
	return nil
}

// 只是简单的识别字符串
func (s *Scanner) String() *TokenJson {
	var ans string
	escape := false // 是否处于转义状态的标志
	for {
		r, _, err := s.Reader.ReadRune()
		if err != nil {
			return tokenJson(FAILED, s.Pos)
		}
		if r == '\n' {
			s.updatePos()
		} else {
			s.Pos.Column++
		}
		if escape {
			// 如果处于转义状态，则直接将字符添加到 ans 中，重置转义状态
			ans += string(r)
			escape = false
			continue
		}
		if r == '\\' {
			// 如果遇到转义字符，则将转义状态置为 true
			escape = true
			continue
		}
		if r == '"' {
			return NewTokenJson(String, ans, s.Pos)
		}
		ans += string(r)
	}
}

// 识别布尔类型
func (s *Scanner) BoolOrNone() *TokenJson {
	var ans string
	for {
		r, _, err := s.Reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return s.parseBoolOrNone(ans)
			}
			return tokenJson(FAILED, s.Pos)
		}
		s.Pos.Column++
		if unicode.IsLetter(r) {
			ans += string(r)
		} else {
			s.backUp()
			return s.parseBoolOrNone(ans)
		}
	}
}

// 识别数字信息
func (s *Scanner) Number() *TokenJson {
	var ans string
	var flag bool
	for {
		r, _, err := s.Reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return s.parseNumber(ans, flag)
			}
			return tokenJson(FAILED, s.Pos)
		}
		s.Pos.Column++
		if unicode.IsDigit(r) || r == '.' || r == 'e' || r == 'E' || r == '+' || r == '-' {
			ans += string(r)
			if r == '.' || r == 'e' || r == 'E' {
				flag = true
			}
		} else {
			s.backUp()
			return s.parseNumber(ans, flag)
		}
	}
}
func (s *Scanner) parseNumber(ans string, flag bool) *TokenJson {
	if flag {
		num, err := strconv.ParseFloat(ans, 64)
		if err != nil {
			return tokenJson(FAILED, s.Pos)
		}
		return NewTokenJson(Number, num, s.Pos)
	}
	num, err := strconv.ParseInt(ans, 10, 64)
	if err != nil {
		return tokenJson(FAILED, s.Pos)
	}
	return NewTokenJson(Number, num, s.Pos)
}
func (s *Scanner) parseBoolOrNone(ans string) *TokenJson {
	if ans == "null" {
		return NewTokenJson(Null, nil, s.Pos)
	}
	b, err := strconv.ParseBool(ans)
	if err != nil {
		return tokenJson(FAILED, s.Pos)
	}
	return NewTokenJson(Bool, b, s.Pos)
}

// parse结构体
type Parse struct {
	Scan *Scanner
}

// 实例化 parse 结构体
func NewParse(s *Scanner) *Parse {
	return &Parse{Scan: s}
}

// 对于json信息进行解析处理
func (p *Parse) Json() *TokenJson {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Sprintf("Error parsing JSON: %v", r))
		}
	}()
	tokenJson := p.Scan.Scan()
	// 出现错误的话 直接panic即可

	switch tokenJson.Type {
	case ObjectStart:
		return p.parseObject()
	case ArrayStart:
		return p.parseArray()
	case FAILED:
		panic(fmt.Sprintf("Error parsing JSON: %v,%v,%v", tokenJson, tokenJson.Value, tokenJson.Pos))
	default:
		return tokenJson
	}
}

// parseArray处理数组或切片类型的数据
func (p *Parse) parseObject() *TokenJson {
	ma := make(map[interface{}]interface{}, 0)
	for {
		// 对于key部分进行简单的解析
		keyTokenJson := p.Json()
		if keyTokenJson.Type == ObjectEnd {
			break
		}
		// 解析 : 分割符号
		sepTokenJson := p.Json()
		if sepTokenJson.Type != SeparatorColon {
			return tokenJson(FAILED, p.Scan.Pos)
		}

		// 对于value部分进行简单的解析
		valTokenJson := p.Json()
		if keyTokenJson.Type == ObjectEnd {
			break
		}
		// 数据填充
		ma[keyTokenJson.Value] = valTokenJson.Value
		// 对于后续字段 ，}进行解析
		nextTokenJson := p.Json()
		if nextTokenJson.Type == ObjectEnd {
			break
		}
		if nextTokenJson.Type != SeparatorComma {
			return tokenJson(FAILED, p.Scan.Pos)
		}
	}
	return NewTokenJson(Map, ma, p.Scan.Pos)
}

// parseArray处理数组或切片类型的数据
func (p *Parse) parseArray() *TokenJson {
	arr := make([]interface{}, 0)
	for {
		// 只需要进行简单的value解析即可
		tokenJson := p.Json()
		if tokenJson.Type == SeparatorComma {
			continue
		}
		if tokenJson.Type == ArrayEnd {
			break
		}
		arr = append(arr, tokenJson.Value)
	}
	return NewTokenJson(Array, arr, p.Scan.Pos)
}
