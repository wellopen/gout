package _http

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
)

func TestFields_Json(t *testing.T) {
	a := readers{
		list: []Reader{
			Param("name", 123),
			Param("ok", 1.1415926),
			Param("body", "good"),
			Param("level", "top"),
		},
	}
	bs, err := ioutil.ReadAll(a.Json())
	require.NoError(t, err)
	fmt.Println(string(bs))
}

func TestFields_UrlEncoded(t *testing.T) {
	a := readers{
		list: []Reader{
			Param("name", 123),
			Param("ok", 1.1415926),
			Param("body", "good"),
			Param("level", "top"),
		},
	}
	bs, err := ioutil.ReadAll(a.UrlEncoded())
	require.NoError(t, err)
	fmt.Println(string(bs))
}
