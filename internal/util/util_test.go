package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundFloat(t *testing.T) {
	t.Run("this test simulate round a value", func(t *testing.T) {
		value := RoundFloat(25.6154, 2)
		assert.Equal(t, value, 25.62)
	})

	t.Run("this test simulate round a value", func(t *testing.T) {
		value := RoundFloat(25.6100002, 2)
		assert.Equal(t, value, 25.61)
	})
}
