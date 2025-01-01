package helper

import (
	"strings"

	"github.com/google/uuid"
)

// UUID uuid
func UUID() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

func UUIDv4() string {
	return uuid.New().String()
}
