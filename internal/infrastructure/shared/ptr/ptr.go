package ptr

// Of 将值转换为指向该值的指针
func Of[T any](value T) *T {
	return &value
}

// Value 从指针中获取值，如果指针为nil则返回零值
// 注意：如果指针为nil，返回的零值可能不是你期望的零值，
// 例如：*int 为 0，*string 为 ""，*bool 为 false 等。
func Value[T any](value *T) T {
	if value == nil {
		var zero T
		return zero
	}
	return *value
}

// SliceOf 将值切片转换为指针切片
func SliceOf[T any](value []T) []*T {
	if value == nil {
		return nil
	}
	result := make([]*T, len(value))
	for i, v := range value {
		result[i] = &v
	}
	return result
}

// SliceValue 将指针切片转换为值切片
func SliceValue[T any](value []*T) []T {
	if value == nil {
		return nil
	}
	result := make([]T, len(value))
	for i, v := range value {
		if v == nil {
			var zero T
			result[i] = zero
		} else {
			result[i] = *v
		}
	}
	return result
}

// MapOf 将值映射转换为指针映射
func MapOf[K comparable, V any](value map[K]V) map[K]*V {
	if value == nil {
		return nil
	}
	result := make(map[K]*V, len(value))
	for k, v := range value {
		result[k] = &v
	}
	return result
}

// MapValue 将指针映射转换为值映射
func MapValue[K comparable, V any](value map[K]*V) map[K]V {
	if value == nil {
		return nil
	}
	result := make(map[K]V, len(value))
	for k, v := range value {
		if v == nil {
			var zero V
			result[k] = zero
		} else {
			result[k] = *v
		}
	}
	return result
}
