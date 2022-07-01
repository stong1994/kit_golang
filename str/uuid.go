package str

import (
	"github.com/google/uuid"
	"strings"
)

func UUID() string {
	return uuid.New().String()
}

func UUIDHex() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
