package _rand

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// Int Int returns a uniform random value in [min, max]. It panics if max <= 0.
func Int(min, max int) (r int, err error) {
	_r, err := Int64(int64(min), int64(max))
	if err != nil {
		return
	}
	r = int(_r)
	return
}

func Int32(min, max int32) (r int32, err error) {
	_r, err := Int64(int64(min), int64(max))
	if err != nil {
		return
	}
	r = int32(_r)
	return
}

func Int64(min, max int64) (r int64, err error) {
	if min == max {
		r = min
		return
	}
	if min > max {
		max, min = min, max
	}
	result, err := rand.Int(rand.Reader, big.NewInt(max-min))
	if err != nil {
		return
	}
	r = result.Int64() + min
	
	b, err := rand.Int(rand.Reader, big.NewInt(2))
	if err == nil && b.Int64() == 1 {
		r += 1
	}
	
	return
}

func Code(length int) (code string, err error) {
	var r int64
	for i := 0; i < length; i++ {
		r, err = Int64(1, 9)
		if err != nil {
			return
		}
		code += fmt.Sprintf("%d", r)
	}
	return
}
