package sstr

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUUID(t *testing.T) {
	str := UUID()
	fmt.Println(str)
	assert.Equal(t, 36, len(str))
}

func TestUUIDHex(t *testing.T) {
	str := UUIDHex()
	fmt.Println(str)
	assert.Equal(t, 32, len(str))
}
