package _http

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
	"time"
)

func testUrl(uri string) string {
	return fmt.Sprintf("http://127.0.0.1:1500%s", uri)
}

func TestNewRequest(t *testing.T) {
	resp := Get(testUrl("/test/params"), Param("method", GET), Param("name", "tom"), Param("age", 18))
	require.NoError(t, resp.Error())
	t.Log(resp.String())
	
	resp = Delete(testUrl("/test/params"), Param("method", DELETE), Param("name", "tom"), Param("age", 18))
	require.NoError(t, resp.Error())
	t.Log(resp.String())
	
	resp = Post(testUrl("/test/json"), Param("name", "tom"), Param("age", 18))
	require.NoError(t, resp.Error())
	t.Log(resp.String())
	
	req := NewRequest(POST, testUrl("/test/form-urlencoded"), Param("name", "tom"), Param("age", 18)).SetContentType(X_WWW_FORM_URLENCODED)
	resp = req.Exec()
	require.NoError(t, resp.Error())
	t.Log(resp.String())
	
	req = NewRequest(POST, testUrl("/test/form-data"), Param("name", "tom"), Param("age", 18), File("file", "./test/test.txt")).SetContentType(FORM_DATA)
	resp = req.Exec()
	require.NoError(t, resp.Error())
	t.Log(strings.TrimSpace(resp.String()))
	j, err := resp.Json()
	require.NoError(t, err)
	path, err := j.Get("file").String()
	require.NoError(t, err)
	_, err = os.Stat(path)
	require.NoError(t, err)
	_ = os.Remove(path)
	
	req = NewRequest(GET, testUrl("/test/timeout")).SetTimeout(1 * time.Second)
	resp = req.Exec()
	require.Error(t, resp.Error())
}
