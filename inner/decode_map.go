package inner

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
