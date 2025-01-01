package helper

import (
	"fmt"
	"log"
	"net"
	"regexp"
)

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

// RemoveMarkdownLink 去除markdown文档中的超链接
func RemoveMarkdownLink(markdown string) string {
	regex := regexp.MustCompile(`\[.*?\(.*?\)`)
	return regex.ReplaceAllString(markdown, "")
}
