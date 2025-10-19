package serializer

import (
	"github.com/vmihailenco/msgpack/v5"

	"github.com/dysodeng/app/internal/infrastructure/persistence/cache/contracts"
)

type msgpackSerializer[T any] struct{}

func NewMsgpackSerializer[T any]() contracts.Serializer[T] {
	return &msgpackSerializer[T]{}
}

func (c msgpackSerializer[T]) Encode(v T) ([]byte, error) {
	return msgpack.Marshal(v)
}

func (c msgpackSerializer[T]) Decode(b []byte) (T, error) {
	var t T
	err := msgpack.Unmarshal(b, &t)
	return t, err
}
