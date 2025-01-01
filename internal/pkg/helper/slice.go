package helper

import "github.com/samber/lo"

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
