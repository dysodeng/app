package helper

import (
	"strconv"

	"github.com/bytedance/sonic"
)

// IfaceConvertString 接口转换为字符串
func IfaceConvertString(data interface{}) string {
	var value string
	switch t := data.(type) {
	case string:
		value = t
	case []byte:
		value = BytesToString(t)
	case int8:
		value = strconv.Itoa(int(t))
	case uint8:
		value = strconv.Itoa(int(t))
	case int16:
		value = strconv.Itoa(int(t))
	case uint16:
		value = strconv.Itoa(int(t))
	case int:
		value = strconv.Itoa(t)
	case uint:
		value = strconv.Itoa(int(t))
	case int32:
		value = strconv.Itoa(int(t))
	case uint32:
		value = strconv.Itoa(int(t))
	case int64:
		value = strconv.FormatInt(t, 10)
	case uint64:
		value = strconv.FormatUint(t, 10)
	case float32:
		value = strconv.FormatFloat(float64(t), 'f', -1, 64)
	case float64:
		value = strconv.FormatFloat(t, 'f', -1, 64)
	case nil:
		value = ""
	default:
		jsonByte, _ := sonic.Marshal(data)
		value = string(jsonByte)
	}
	return value
}

// IfaceConvertInt64 接口类型转换为int64
func IfaceConvertInt64(data interface{}) int64 {
	var value int64
	switch t := data.(type) {
	case string:
		value, _ = strconv.ParseInt(t, 10, 64)
	case []byte:
		v := BytesToString(t)
		value, _ = strconv.ParseInt(v, 10, 64)
	case int8:
		value = int64(t)
	case uint8:
		value = int64(t)
	case int16:
		value = int64(t)
	case uint16:
		value = int64(t)
	case int32:
		value = int64(t)
	case uint32:
		value = int64(t)
	case int:
		value = int64(t)
	case uint:
		value = int64(t)
	case int64:
		value = t
	case uint64:
		value = int64(t)
	case float32:
		value = int64(t)
	case float64:
		value = int64(t)
	default:
		value = 0
	}
	return value
}

// IfaceConvertUint64 接口类型转换为uint64
func IfaceConvertUint64(data interface{}) uint64 {
	var value uint64
	switch t := data.(type) {
	case string:
		value, _ = strconv.ParseUint(t, 10, 64)
	case []byte:
		v := BytesToString(t)
		value, _ = strconv.ParseUint(v, 10, 64)
	case int8:
		value = uint64(t)
	case uint8:
		value = uint64(t)
	case int16:
		value = uint64(t)
	case uint16:
		value = uint64(t)
	case int32:
		value = uint64(t)
	case uint32:
		value = uint64(t)
	case int:
		value = uint64(t)
	case uint:
		value = uint64(t)
	case int64:
		value = uint64(t)
	case uint64:
		value = t
	case float32:
		value = uint64(t)
	case float64:
		value = uint64(t)
	default:
		value = 0
	}
	return value
}
