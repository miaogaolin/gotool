package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testUser struct {
	Height int
}

func TestSortStructSlice(t *testing.T) {
	data := []testUser{
		{1},
		{2},
	}
	err := SortStructSlice(data, false, "Height")
	assert.Nil(t, err)
	assert.Equal(t, []testUser{
		{2},
		{1},
	}, data)
}
