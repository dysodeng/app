package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

const (
	letterBytes   = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

// UUID uuid
func UUID() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

func UUIDv4() string {
	return uuid.New().String()
}

// GeneratePassword 生成密码
// @param string password 明文密码
// @return string
func GeneratePassword(password []byte) string {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	return string(hash)
}

// VerifyPassword 验证密码
// @param string hashedPassword hash密码
// @param string plainPassword 明文密码
// @return bool
func VerifyPassword(hashedPassword string, plainPassword string) bool {
	byteHashByte := []byte(hashedPassword)
	plainPasswordByte := []byte(plainPassword)

	err := bcrypt.CompareHashAndPassword(byteHashByte, plainPasswordByte)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// GenValidateCode 生成指定长度数字字符串
// @param int length 生成字符串长度
// @return string
func GenValidateCode(length int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.New(rand.NewSource(time.Now().UnixMilli()))

	var buf strings.Builder
	for i := 0; i < length; i++ {
		_, _ = fmt.Fprintf(&buf, "%d", numeric[rand.Intn(r)])
	}

	return buf.String()
}

// RandomStringBytesMask 生成随机字符串
// @param int length 生成字符串长度
// @return string
func RandomStringBytesMask(length int) string {

	str := make([]byte, length)

	rand.New(rand.NewSource(time.Now().UnixMilli()))

	for i, cache, reMain := length-1, rand.Int63(), letterIdxMax; i >= 0; {
		if reMain == 0 {
			cache, reMain = rand.Int63(), letterIdxMax
		}

		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			str[i] = letterBytes[idx]
			i--
		}

		cache >>= letterIdxBits
		reMain--
	}

	return string(str)
}

// GbkToUtf8 GBK 转 UTF-8
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := io.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// Utf8ToGbk Utf8 转 gbk
func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := io.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// GetLocalIp 获取本机IP地址
// @return string
func GetLocalIp() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = conn.Close()
	}()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

// CreateOrderNo 生成唯一订单号
// @return string
func CreateOrderNo() string {
	sTime := time.Now().Format("20060102150405")

	t := time.Now().UnixNano()
	s := strconv.FormatInt(t, 10)
	b := string([]byte(s)[len(s)-9:])
	c := string([]byte(b)[:7])

	rand.Seed(t)

	sTime += c + strconv.FormatInt(rand.Int63n(999999-100000)+100000, 10)
	return sTime
}

// ResolveTime 将整数转换为时分秒
// @param int seconds 秒数
// @return int hour 小时数
// @return int minute 分钟数
// @return int second 秒数
func ResolveTime(seconds int) (hour, minute, second int) {
	hour = seconds / 3600
	minute = (seconds - hour*3600) / 60
	second = seconds - hour*3600 - minute*60
	return
}

// FileSizeFormat 格式化文件大小
// @param uint64 fileSize 文件大小(字节)
// @return string
func FileSizeFormat(fileSize uint64) string {
	var byteNum float64 = 1024 // byte
	size := float64(fileSize)

	if size < byteNum { // B
		return fmt.Sprintf("%f", size) + "B"
	} else if size < math.Pow(byteNum, 2) { // KB
		return fmt.Sprintf("%.2f", size/byteNum) + "KB"
	} else if size < math.Pow(byteNum, 3) { // MB
		return fmt.Sprintf("%.2f", size/math.Pow(byteNum, 2)) + "MB"
	} else if size < math.Pow(byteNum, 4) { // GB
		return fmt.Sprintf("%.2f", size/math.Pow(byteNum, 3)) + "GB"
	}
	return fmt.Sprintf("%.2f", size/math.Pow(byteNum, 4)) + "TB"
}

// Find 查找元素在数组中的位置
func Find[T comparable](item T, itemList []T) int {
	for i, l := range itemList {
		if item == l {
			return i
		}
	}
	return -1
}

// Contain 判断元素是否包含在列表中
func Contain[T comparable](item T, itemList []T) bool {
	for _, i := range itemList {
		if item == i {
			return true
		}
	}
	return false
}

// IsContainSlice 判断2个切片是否有包含关系
func IsContainSlice[T comparable](minSlice []T, maxSlice []T) bool {
	var set = make(map[T]struct{})
	for _, value := range maxSlice {
		set[value] = struct{}{}
	}

	for _, minValue := range minSlice {
		if _, ok := set[minValue]; !ok {
			return false
		}
	}

	return true
}

// DiffSlice 切片差集
func DiffSlice[T comparable](firstSlice []T, lastSlice []T) []T {
	var diff []T
	for _, t := range lastSlice {
		if !Contain(t, firstSlice) {
			diff = append(diff, t)
		}
	}
	for _, t := range firstSlice {
		if !Contain(t, lastSlice) {
			diff = append(diff, t)
		}
	}
	return diff
}

// IfaceConvertString 接口转换为字符串
func IfaceConvertString(data interface{}) string {
	var value string
	switch data.(type) {
	case string:
		value = data.(string)
		break
	case []byte:
		value = string(data.([]byte))
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
		v := string(data.([]byte))
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
		v := string(data.([]byte))
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

// Cartesian 计算矩阵笛卡尔积
func Cartesian(data [][]uint64, delimiter string) []string {
	if delimiter == "" {
		delimiter = ","
	}

	// 保存结果
	var result []string

	count := len(data)
	if count > 1 {
		for i := 0; i < count-1; i++ {
			// 初始化
			if i == 0 {
				for _, u := range data[i] {
					result = append(result, fmt.Sprintf("%d", u))
				}
			}

			// 保存临时数据
			var tmp []string

			// 结果与下一个集合计算笛卡尔积
			for _, res := range result {
				for _, set := range data[i+1] {
					tmp = append(tmp, fmt.Sprintf("%s%s%d", res, delimiter, set))
				}
			}

			// 将笛卡尔积回写入结果
			result = tmp
		}
	} else { // 处理特殊情况
		if count > 0 {
			var resData []string
			if len(data[0]) > 0 {
				for _, u := range data[0] {
					resData = append(resData, fmt.Sprintf("%d", u))
				}
			} else {
				resData = []string{}
			}
			return resData
		}

		return []string{}
	}

	return result
}

// HideCellphone 隐藏手机号码中间四位
func HideCellphone(cellphone string) string {
	if len(cellphone) > 7 {
		return cellphone[:3] + "****" + cellphone[7:]
	} else if len(cellphone) > 3 {
		return cellphone[:3] + "****"
	} else if len(cellphone) > 0 {
		return "****"
	}
	return cellphone
}

func HideEmail(email string) string {
	if email == "" {
		return ""
	}
	temp := strings.Split(email, "@")
	if len(temp) != 2 || len(temp[0]) == 0 || len(temp[1]) == 0 {
		return email[0:1] + "****"
	}
	prefix := temp[0]
	suffix := temp[1]
	return prefix[0:1] + "****@" + suffix
}

// HideRealName 隐藏姓名中间文字
func HideRealName(name string) string {
	n := []rune(name)
	l := len(n)
	var newName string
	if l >= 2 {
		if l > 2 {
			newName = string(n[0]) + "*" + string(n[l-1])
		} else {
			newName = string(n[0]) + "*"
		}
	} else {
		return name
	}
	return newName
}

// RandAreaNum 区间随机数(2端闭合区间)
func RandAreaNum(min, max int) int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return rand.Intn(max-min+1) + min
}

// BigNumberThousandFormat 大数千分位分割格式化
func BigNumberThousandFormat[T ~uint | ~uint16 | ~uint32 | ~uint64 | int | int16 | int32 | int64 | float32 | float64](num T) string {
	printer := message.NewPrinter(language.English)
	return printer.Sprintf("%d", num)
}

// IsToday 是否为今天
func IsToday(t time.Time) bool {
	now := time.Now().Local()
	return t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day()
}

// IsSameYear 是否为同一年
func IsSameYear(t time.Time) bool {
	now := time.Now().Local()
	return t.Year() == now.Year()
}

// LastDayOfMonth 获取月份最后一天日期
func LastDayOfMonth(t time.Time) time.Time {
	t = t.Local()
	year := t.Year()
	month := int(t.Month())
	lastDay := time.Date(year, time.Month(month+1), 0, 23, 59, 59, 0, time.Local).Add(-time.Nanosecond)
	return lastDay
}

// MonthDays 获取月份天数
func MonthDays(t time.Time) int {
	year := t.Year()
	month := int(t.Month())
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 4, 6, 9, 11:
		return 30
	case 2:
		if ((year%4) == 0 && (year%100) != 0) || (year%400) == 0 {
			return 29
		}
		return 28
	default:
		panic("无效的月份")
	}
}

// WeekDay 获取星期几
func WeekDay(t time.Time) uint8 {
	week := t.Weekday()
	switch week {
	case time.Sunday:
		return 7
	case time.Monday:
		return 1
	case time.Tuesday:
		return 2
	case time.Wednesday:
		return 3
	case time.Thursday:
		return 4
	case time.Friday:
		return 5
	case time.Saturday:
		return 6
	}
	return 0
}

// WeekChineseDay 获取星期几中文名称
func WeekChineseDay(t time.Time) string {
	week := WeekDay(t)
	switch week {
	case 1:
		return "周一"
	case 2:
		return "周二"
	case 3:
		return "周三"
	case 4:
		return "周四"
	case 5:
		return "周五"
	case 6:
		return "周六"
	case 7:
		return "周日"
	}
	return ""
}

func WeekChinese(t time.Time) string {
	weekday := t.Weekday()
	var cnWeekday string

	switch weekday {
	case time.Monday:
		cnWeekday = "星期一"
	case time.Tuesday:
		cnWeekday = "星期二"
	case time.Wednesday:
		cnWeekday = "星期三"
	case time.Thursday:
		cnWeekday = "星期四"
	case time.Friday:
		cnWeekday = "星期五"
	case time.Saturday:
		cnWeekday = "星期六"
	case time.Sunday:
		cnWeekday = "星期日"
	}

	return cnWeekday
}

// TraceTime 消耗时间计算
// 用法 defer TraceTime(time.Now()) 可计算函数执行时间
func TraceTime(pre time.Time) time.Duration {
	return time.Since(pre)
}

// ReplaceString 批量替换字符串
func ReplaceString(str string, findSlice, replaceSlice []string) string {
	if len(findSlice) != len(replaceSlice) {
		return str
	}

	for i, find := range findSlice {
		str = strings.Replace(str, find, replaceSlice[i], -1)
	}

	return str
}

// RemoveMarkdownImages 去除markdown文档中的图片
func RemoveMarkdownImages(markdown string) string {
	regex := regexp.MustCompile(`!\[.*?\]\(.*?\)`)
	return regex.ReplaceAllString(markdown, "")
}
