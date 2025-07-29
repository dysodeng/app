package helper

import (
	"fmt"
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
	if fileSize == 0 {
		return "0B"
	}

	units := []string{"B", "KB", "MB", "GB", "TB"}
	size := float64(fileSize)

	for i, unit := range units {
		if size < 1024 || i == len(units)-1 {
			if unit == "B" {
				return fmt.Sprintf("%.0f%s", size, unit)
			}
			return fmt.Sprintf("%.2f%s", size, unit)
		}
		size /= 1024
	}

	return fmt.Sprintf("%.2fTB", size)
}
