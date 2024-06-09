package random

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRandomSTR(t *testing.T) {
	testdata := []int{1, 2, 3, 6, 10}
	for _, size := range testdata {
		res := NewRandomSTR(size)
		require.Equal(t, size, len(res))
	}
}
