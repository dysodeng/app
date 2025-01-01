package helper

import (
	"encoding/json"
	"strconv"
)

// IfaceConvertString 接口转换为字符串
func IfaceConvertString(data interface{}) string {
	var value string
	switch data.(type) {
	case string:
		value = data.(string)
		break
	case []byte:
		value = BytesToString(data.([]byte))
		break
	case int8:
		it := data.(int8)
		value = strconv.Itoa(int(it))
		break
	case uint8:
		it := data.(uint8)
		value = strconv.Itoa(int(it))
		break
	case int16:
		it := data.(int16)
		value = strconv.Itoa(int(it))
		break
	case uint16:
		it := data.(uint16)
		value = strconv.Itoa(int(it))
		break
	case int:
		it := data.(int)
		value = strconv.Itoa(it)
		break
	case uint:
		it := data.(uint)
		value = strconv.Itoa(int(it))
		break
	case int32:
		it := data.(int32)
		value = strconv.Itoa(int(it))
		break
	case uint32:
		it := data.(uint32)
		value = strconv.Itoa(int(it))
		break
	case int64:
		it := data.(int64)
		value = strconv.FormatInt(it, 10)
		break
	case uint64:
		it := data.(uint64)
		value = strconv.FormatUint(it, 10)
		break
	case float32:
		ft := data.(float32)
		value = strconv.FormatFloat(float64(ft), 'f', -1, 64)
		break
	case float64:
		ft := data.(float64)
		value = strconv.FormatFloat(ft, 'f', -1, 64)
		break
	case nil:
		value = ""
		break
	default:
		jsonByte, _ := json.Marshal(data)
		value = string(jsonByte)
	}
	return value
}

// IfaceConvertInt64 接口类型转换为int64
func IfaceConvertInt64(data interface{}) int64 {
	var value int64
	switch data.(type) {
	case string:
		v := data.(string)
		value, _ = strconv.ParseInt(v, 10, 64)
		break
	case []byte:
		v := BytesToString(data.([]byte))
		value, _ = strconv.ParseInt(v, 10, 64)
		break
	case int8:
		value = int64(data.(int8))
		break
	case uint8:
		value = int64(data.(uint8))
		break
	case int16:
		value = int64(data.(int16))
		break
	case uint16:
		value = int64(data.(uint16))
		break
	case int32:
		value = int64(data.(int32))
		break
	case uint32:
		value = int64(data.(uint32))
		break
	case int:
		value = int64(data.(int))
		break
	case uint:
		value = int64(data.(uint))
		break
	case int64:
		value = data.(int64)
		break
	case uint64:
		value = int64(data.(uint64))
		break
	case float32:
		value = int64(data.(float32))
		break
	case float64:
		value = int64(data.(float64))
		break
	default:
		value = 0
		break
	}
	return value
}

// IfaceConvertUint64 接口类型转换为uint64
func IfaceConvertUint64(data interface{}) uint64 {
	var value uint64
	switch data.(type) {
	case string:
		v := data.(string)
		value, _ = strconv.ParseUint(v, 10, 64)
		break
	case []byte:
		v := BytesToString(data.([]byte))
		value, _ = strconv.ParseUint(v, 10, 64)
		break
	case int8:
		value = uint64(data.(int8))
		break
	case uint8:
		value = uint64(data.(uint8))
		break
	case int16:
		value = uint64(data.(int16))
		break
	case uint16:
		value = uint64(data.(uint16))
		break
	case int32:
		value = uint64(data.(int32))
		break
	case uint32:
		value = uint64(data.(uint32))
		break
	case int:
		value = uint64(data.(int))
		break
	case uint:
		value = uint64(data.(uint))
		break
	case int64:
		value = uint64(data.(int64))
		break
	case uint64:
		value = data.(uint64)
		break
	case float32:
		value = uint64(data.(float32))
		break
	case float64:
		value = uint64(data.(float64))
		break
	default:
		value = 0
		break
	}
	return value
}
