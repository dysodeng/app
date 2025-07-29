package helper

import (
	"strings"

	"github.com/google/uuid"
)

// UUID uuid
func UUID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func UUIDv4() string {
	return uuid.New().String()
}
