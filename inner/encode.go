package inner

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"unsafe"
)

// Encode data from slice/array + map + struct
func Encode(data interface{}) []byte {
	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		return encodeList(v)
	case reflect.Map:
		return encodeMap(v)
	case reflect.Struct:
		return encodeStruct(v)
	default:
		return nil
	}
}
func SEncode(data interface{}) string {
	return bytesToString(Encode(data))
}

// encode list
func encodeList(v reflect.Value) []byte {
	var buffer bytes.Buffer
	buffer.WriteByte('[')

	for i := 0; i < v.Len(); i++ {
		if i > 0 {
			buffer.WriteString(",")
		}
		buffer.Write(checkDataType(v.Index(i)))
	}
	buffer.WriteByte(']')
	return buffer.Bytes()
}

// encode Map
func encodeMap(v reflect.Value) []byte {
	var buffer bytes.Buffer
	buffer.WriteByte('{')

	keys := v.MapKeys()
	for i, key := range keys {
		if i > 0 {
			buffer.WriteString(",")
		}
		buffer.Write(checkDataType(key))
		buffer.WriteByte(':')
		buffer.Write(checkDataType(v.MapIndex(key)))
	}
	buffer.WriteByte('}')
	return buffer.Bytes()
}

// encode struct
func encodeStruct(v reflect.Value) []byte {
	var buffer bytes.Buffer
	buffer.WriteByte('{')

	for i := 0; i < v.NumField(); i++ {
		if i > 0 {
			buffer.WriteString(",")
		}
		field := v.Field(i)
		fieldType := v.Type().Field(i)
		buffer.WriteString(`"`)
		buffer.WriteString(fieldType.Name)
		buffer.WriteString(`":`)
		buffer.Write(checkDataType(field))
	}
	buffer.WriteByte('}')
	return buffer.Bytes()
}

// check Data Type
func checkDataType(item reflect.Value) []byte {
	var buffer bytes.Buffer

	switch item.Kind() {
	case reflect.Slice, reflect.Array:
		buffer.Write(encodeList(item))
	case reflect.Map:
		buffer.Write(encodeMap(item))
	case reflect.Struct:
		buffer.Write(encodeStruct(item))
	case reflect.String:
		buffer.WriteString(strconv.Quote(item.String()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		buffer.WriteString(strconv.FormatInt(item.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		buffer.WriteString(strconv.FormatUint(item.Uint(), 10))
	case reflect.Float32, reflect.Float64:
		buffer.WriteString(strconv.FormatFloat(item.Float(), 'f', -1, 64))
	case reflect.Bool:
		buffer.WriteString(strconv.FormatBool(item.Bool()))
	case reflect.Interface:
		if item.IsNil() {
			buffer.WriteString("null")
			break
		}
		interfaceValue := item.Elem()
		switch interfaceValue.Kind() {
		case reflect.Slice, reflect.Array:
			buffer.Write(encodeList(interfaceValue))
		case reflect.String:
			buffer.WriteString(strconv.Quote(interfaceValue.String()))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			buffer.WriteString(strconv.FormatInt(interfaceValue.Int(), 10))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			buffer.WriteString(strconv.FormatUint(interfaceValue.Uint(), 10))
		case reflect.Float32, reflect.Float64:
			buffer.WriteString(strconv.FormatFloat(interfaceValue.Float(), 'f', -1, 64))
		case reflect.Bool:
			buffer.WriteString(strconv.FormatBool(interfaceValue.Bool()))
		case reflect.Map:
			buffer.Write(encodeMap(interfaceValue))
		case reflect.Struct:
			buffer.Write(encodeStruct(item))
		default:
			_, _ = fmt.Fprintf(&buffer, "%v", interfaceValue.Interface())

		}
	default:
		panic("unhandled default case")
	}
	return buffer.Bytes()
}
func bytesToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
