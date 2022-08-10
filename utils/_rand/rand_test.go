package _rand

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRandInt(t *testing.T) {
	for i := 0; i < 100000; i++ {
		r, err := Int(5, 99)
		require.NoError(t, err)
		require.GreaterOrEqual(t, r, 5)
		require.LessOrEqual(t, r, 99)
	}
}

func TestRandInt2(t *testing.T) {
	for i := 0; i < 100; i++ {
		r, err := Int(-100, 1)
		require.NoError(t, err)
		fmt.Println(r)
	}
}

func TestRandCode(t *testing.T) {
	for i := 0; i < 100; i++ {
		r, err := Code(6)
		require.NoError(t, err)
		t.Log(r)
	}
}
