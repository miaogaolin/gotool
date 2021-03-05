package stringx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterface(t *testing.T) {
	var data interface{} = 2
	assert.Equal(t, Interface(data), "2")
}
