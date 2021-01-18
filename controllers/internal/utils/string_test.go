package utils

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatWithStruct(t *testing.T) {
	// arrange
	type str struct {
		Name  string
		Value int
	}
	s := str{"Foo", 007}
	expected := `{
	"Name": "Foo",
	"Value": 7
}`
	// act
	result := ToString(s)
	// assert
	assert.Equal(t, expected, result)
}

func TestFormatWithPrimitiveType(t *testing.T) {
	// arrange
	// act
	result := ToString(true)
	// assert
	assert.Equal(t, "true", result)
}

func TestFormatWithNilPointerReference(t *testing.T) {
	// arrange
	type str struct {
		Name  string
		Value int
	}
	var ptr *str = nil
	// act
	result := ToString(ptr)
	// assert
	assert.Equal(t, "null", result)
}

func TestFormatWithCorruptedStructureMetadata(t *testing.T) {
	// arrange
	type str struct {
		Name  string `json:"CorrectName,omitempty"`
		Value int    `json:"Incorrect,OMtEmpt"`
	}
	s := str{"Foo", 007}
	expected := `{
	"CorrectName": "Foo",
	"Incorrect": 7
}`
	// act
	result := ToString(s)
	// assert
	assert.Equal(t, expected, result)
}

func TestFormatWithEmptyStructure(t *testing.T) {
	// arrange
	type str struct {
		Name  string `json:"CorrectName"`
		Value int    `json:"Incorrect"`
	}
	s := str{}
	expected := `{
	"CorrectName": "",
	"Incorrect": 0
}`
	// act
	result := ToString(s)
	// assert
	assert.Equal(t, expected, result)
}

func TestUnsupportedTypeError(t *testing.T) {
	// arrange
	c := make(chan int)
	addr := fmt.Sprintf("%p", c)
	// act
	result := ToString(c)
	// assert
	assert.Equal(t, addr, result)
}

func TestUnsupportedValueError(t *testing.T) {
	// arrange
	v := math.Inf(1)
	// act
	result := ToString(v)
	// assert
	assert.Equal(t, "+Inf", result)
}
