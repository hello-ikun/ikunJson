package inner

import (
	"unsafe"
)

type sliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

// parseArray处理数组或切片类型的数据
func (p *Parse) parseArray() *TokenJson {
	var arr []interface{}

	for {
		// 只需要进行简单的 value 解析即可
		tokenJson := p.Json()
		if tokenJson.Type == SeparatorComma {
			continue
		}
		if tokenJson.Type == ArrayEnd {
			break
		}

		// 这里使用 unsafe.Pointer 进行转换
		// 目前还是没有看懂底层
		arr = append(arr, *(*interface{})(unsafe.Pointer(&tokenJson.Value)))
	}

	return NewTokenJson(Array, arr, p.Scan.Pos)
}
