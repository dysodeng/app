package serializer

import (
	"encoding/json"

	"github.com/dysodeng/app/internal/infrastructure/persistence/cache/contract"
)

type jsonSerializer[T any] struct{}

func NewJSONSerializer[T any]() contract.Serializer[T] {
	return &jsonSerializer[T]{}
}

func (c jsonSerializer[T]) Encode(v T) ([]byte, error) {
	return json.Marshal(v)
}

func (c jsonSerializer[T]) Decode(b []byte) (T, error) {
	var t T
	err := json.Unmarshal(b, &t)
	return t, err
}
