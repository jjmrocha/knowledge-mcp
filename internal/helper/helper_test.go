package helper_test

import (
	"testing"

	"github.com/jjmrocha/knowledge-mcp/internal/helper"
	"github.com/stretchr/testify/assert"
)

func TestToPointer_Int(t *testing.T) {
	// given
	value := 42
	// when
	ptr := helper.ToPointer(value)
	// then
	assert.NotNil(t, ptr)
	assert.Equal(t, value, *ptr)
}

func TestToPointer_String(t *testing.T) {
	// given
	value := "hello"
	// when
	ptr := helper.ToPointer(value)
	// then
	assert.NotNil(t, ptr)
	assert.Equal(t, value, *ptr)
}

func TestToPointer_Struct(t *testing.T) {
	// given
	type point struct{ X, Y int }
	value := point{X: 1, Y: 2}
	// when
	ptr := helper.ToPointer(value)
	// then
	assert.NotNil(t, ptr)
	assert.Equal(t, value, *ptr)
}

func TestToPointer_ZeroValue(t *testing.T) {
	// when
	ptr := helper.ToPointer(0)
	// then
	assert.NotNil(t, ptr)
	assert.Equal(t, 0, *ptr)
}

func TestToPointer_IndependentCopy(t *testing.T) {
	// given
	value := 10
	// when
	ptr := helper.ToPointer(value)
	value = 20
	// then — ptr must hold the original value (10), unaffected by the later change to value (20)
	assert.Equal(t, 10, *ptr)
	assert.Equal(t, 20, value)
}
