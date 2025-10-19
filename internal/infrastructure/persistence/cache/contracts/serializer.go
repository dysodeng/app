package contracts

// Serializer 序列化器接口
type Serializer[T any] interface {
	Encode(v T) ([]byte, error)
	Decode(b []byte) (T, error)
}
