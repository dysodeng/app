package serializer

import (
	"encoding/json"

	"github.com/dysodeng/app/internal/infrastructure/persistence/cache/contract"
)

type jsonSerializer struct{}

func NewJSONSerializer() contract.Serializer {
	return &jsonSerializer{}
}

func (j *jsonSerializer) Serialize(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j *jsonSerializer) Deserialize(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (j *jsonSerializer) ContentType() string {
	return "application/json"
}
