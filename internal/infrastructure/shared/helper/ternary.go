package helper

// Ternary 三元运算符函数，支持泛型
// 用法: result := Ternary(condition, trueValue, falseValue)
func Ternary[T any](condition bool, trueValue, falseValue T) T {
	if condition {
		return trueValue
	}
	return falseValue
}

// TernaryFunc 支持延迟计算的三元运算符
// 用法: result := TernaryFunc(condition, func() T { return trueValue }, func() T { return falseValue })
func TernaryFunc[T any](condition bool, trueFunc, falseFunc func() T) T {
	if condition {
		return trueFunc()
	}
	return falseFunc()
}

// TernaryPtr 针对指针类型的三元运算符，避免nil指针问题
func TernaryPtr[T any](condition bool, trueValue, falseValue *T) *T {
	if condition {
		return trueValue
	}
	return falseValue
}
