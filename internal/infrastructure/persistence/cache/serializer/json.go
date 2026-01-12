package serializer

import (
	"github.com/bytedance/sonic"

	"github.com/dysodeng/app/internal/infrastructure/persistence/cache/contracts"
)

type jsonSerializer[T any] struct{}

func NewJSONSerializer[T any]() contracts.Serializer[T] {
	return &jsonSerializer[T]{}
}

func (c jsonSerializer[T]) Encode(v T) ([]byte, error) {
	return sonic.Marshal(v)
}

func (c jsonSerializer[T]) Decode(b []byte) (T, error) {
	var t T
	err := sonic.Unmarshal(b, &t)
	return t, err
}
