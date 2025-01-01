package helper

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

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
