package inner

import (
	"fmt"
	"reflect"
)

// 解析结构体
type ParseStruct struct {
	Scan *Scanner
}

// 实例化 parse 结构体
func NewParseStruct(s *Scanner) *ParseStruct {
	return &ParseStruct{Scan: s}
}

// 对于json信息进行解析处理
func (p *ParseStruct) Json(data interface{}) *TokenJson {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Sprintf("Error parsing JSON: %v", r))
		}
	}()
	tokenJson := p.Scan.Scan()
	// 出现错误的话 直接panic即可

	switch tokenJson.Type {
	case ObjectStart:
		return p.parseStruct(data)
	case ArrayStart:
		return p.parseArray(data)
	case FAILED:
		panic(fmt.Sprintf("Error parsing JSON: %v,%v,%v", tokenJson, tokenJson.Value, tokenJson.Pos))
	default:
		return tokenJson
	}
}

func (p *ParseStruct) parseArray(data interface{}) *TokenJson {
	dataValue := reflect.ValueOf(data)
	if dataValue.Kind() != reflect.Ptr || dataValue.IsNil() {
		return tokenJson(FAILED, p.Scan.Pos)
	}

	// 创建一个与 dataValue 类型相同的切片
	arrValue := reflect.MakeSlice(dataValue.Type().Elem(), 0, 0)
	// fmt.Println(arrValue, arrValue.Type())
	for {
		// 对于key部分进行简单的解析
		valueTokenJson := p.Json(data)
		if valueTokenJson.Type == ArrayEnd {
			break
		}
		// fmt.Println("data", data, valueTokenJson.Value)
		// 将值追加到切片中
		arrValue = reflect.Append(arrValue, reflect.ValueOf(valueTokenJson.Value))
		sepTokenJson := p.Json(data)
		if sepTokenJson.Type == ArrayEnd {
			break
		}
		if sepTokenJson.Type == EOF {
			return NewTokenJson(Array, data, p.Scan.Pos)
		}
		if sepTokenJson.Type != SeparatorComma {
			return tokenJson(FAILED, p.Scan.Pos)
		}
	}
	// 将切片转换为接口类型返回
	arr := arrValue.Interface()
	return NewTokenJson(Array, arr, p.Scan.Pos)
}

// parseObject 处理结构体类型的数据
func (p *ParseStruct) parseStruct(data interface{}) *TokenJson {
	dataValue := reflect.ValueOf(data)
	if dataValue.IsNil() {
		return NewTokenJson(Struct, nil, p.Scan.Pos)
	}
	if dataValue.Kind() != reflect.Ptr {
		return tokenJson(FAILED, p.Scan.Pos)
	}

	structValue := dataValue.Elem()
	if structValue.Kind() != reflect.Struct {
		return tokenJson(FAILED, p.Scan.Pos)
	}
	for i := 0; i < structValue.NumField(); i++ {
		fieldName := structValue.Type().Field(i).Tag.Get("json")
		if fieldName == "" {
			continue
		}
		// 对于key部分进行简单的解析
		keyTokenJson := p.Json(data)
		if keyTokenJson.Type == ObjectEnd {
			break
		}

		sepTokenJson := p.Json(data)

		if sepTokenJson.Type != SeparatorColon {
			return tokenJson(FAILED, p.Scan.Pos)
		}
		// 注意错误 此处传入的data是val的地址
		valTokenJson := p.Json(structValue.Field(i).Addr().Interface())
		if keyTokenJson.Type == ObjectEnd {
			break
		}
		// 使用反射动态填充结构体字段
		field := structValue.Field(i)
		if field.CanSet() {
			switch field.Kind() {
			case reflect.Struct:
				// 如果是 struct类的话 我们直接实现迭代即可
				field.Set(reflect.ValueOf(valTokenJson.Value).Elem())
			case reflect.Array, reflect.Slice:
				fmt.Println(valTokenJson.Value)
				field.Set(reflect.ValueOf(valTokenJson.Value))
			case reflect.Interface: //此处传入的是 null类型 不做处理即可
				break
			default:
				fmt.Println(valTokenJson.Value)
				field.Set(reflect.ValueOf(valTokenJson.Value).Convert(field.Type()))
			}
		}
		nextTokenJson := p.Json(data)
		if nextTokenJson.Type == ObjectEnd {
			break
		}
		if nextTokenJson.Type == EOF {
			return NewTokenJson(Struct, data, p.Scan.Pos)
		}
		if nextTokenJson.Type != SeparatorComma {
			return tokenJson(FAILED, p.Scan.Pos)
		}
	}
	return NewTokenJson(Struct, data, p.Scan.Pos)
}
