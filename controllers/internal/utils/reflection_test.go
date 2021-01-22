package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTypeStruct(t *testing.T) {
	// arrange
	type ab struct{}
	x := ab{}
	// act
	r := GetType(x)
	// assert
	assert.Equal(t, "ab", r)
}

func TestGetTypePrimitive(t *testing.T) {
	// arrange
	x := 5
	// act
	r := GetType(x)
	// assert
	assert.Equal(t, "int", r)
}

func TestGetTypePointer(t *testing.T) {
	// arrange
	type ab struct{}
	x := &ab{}
	// act
	r := GetType(x)
	// assert
	assert.Equal(t, "*ab", r)
}
