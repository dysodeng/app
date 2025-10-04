package helper

import (
	"math/rand"

	"github.com/samber/lo"
)

// IndexOf 获取元素在数组中的位置
func IndexOf[T comparable](itemList []T, item T) int {
	return lo.IndexOf(itemList, item)
}

// Contain 判断元素是否包含在列表中
func Contain[T comparable](itemList []T, item T) bool {
	return lo.Contains(itemList, item)
}

// IsContainSlice 判断2个切片是否有包含关系
func IsContainSlice[T comparable](maxSlice []T, minSlice []T) bool {
	return lo.Every(maxSlice, minSlice)
}

// DiffSlice 切片差集
func DiffSlice[T comparable](firstSlice []T, lastSlice []T) []T {
	var diff []T
	for _, t := range lastSlice {
		if !Contain(firstSlice, t) {
			diff = append(diff, t)
		}
	}
	for _, t := range firstSlice {
		if !Contain(lastSlice, t) {
			diff = append(diff, t)
		}
	}
	return diff
}

// RandomSliceUnique 从切片中随机抽取 n 条不重复元素
func RandomSliceUnique[T any](src []T, n int) []T {
	l := len(src)
	if n <= 0 {
		return []T{}
	}
	if n > l {
		n = l // 如果 n 大于切片长度，取全部
	}

	indices := rand.Perm(l)[:n] // 生成不重复的随机索引
	result := make([]T, n)
	for i, idx := range indices {
		result[i] = src[idx]
	}
	return result
}
