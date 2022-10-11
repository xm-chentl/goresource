package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEmtty(test *testing.T) {
	test.Run("interface.empty.true", func(t *testing.T) {
		a := assert.New(t)
		var v interface{}
		a.Equal(IsEmpty(v), true)
	})

	test.Run("interface.empty.false", func(t *testing.T) {
		a := assert.New(t)
		var v interface{} = 5
		a.Equal(IsEmpty(v), false)
	})

	test.Run("string.empty.true", func(t *testing.T) {
		a := assert.New(t)
		a.Equal(IsEmpty(""), true)
	})

	test.Run("string.empty.false", func(t *testing.T) {
		a := assert.New(t)
		a.Equal(IsEmpty("test"), false)
	})

	test.Run("int.zero.true", func(t *testing.T) {
		a := assert.New(t)
		a.Equal(IsEmpty(0), true)
	})

	test.Run("int.zero.false", func(t *testing.T) {
		a := assert.New(t)
		a.Equal(IsEmpty(12), false)
	})

	test.Run("int64.zero.true", func(t *testing.T) {
		a := assert.New(t)
		v := int64(0)
		a.Equal(IsEmpty(v), true)
	})

	test.Run("int64.zero.false", func(t *testing.T) {
		a := assert.New(t)
		v := int64(23)
		a.Equal(IsEmpty(v), false)
	})

}
