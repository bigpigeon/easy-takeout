package api

import (
	"encoding/json"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getFormData(m map[string]string) string {
	form := url.Values{}
	for k, v := range m {
		form.Add(k, v)
	}
	return form.Encode()
}

func getJsonData(t *testing.T, v interface{}) string {
	b, err := json.Marshal(v)
	assert.Equal(t, err, nil)
	return string(b)
}

func fmtJsonBytes(t *testing.T, bytes []byte) string {
	var v interface{}
	err := json.Unmarshal(bytes, &v)
	assert.Equal(t, err, nil)
	newb, err := json.MarshalIndent(&v, "", "  ")
	assert.Equal(t, err, nil)
	return string(newb)
}

func promiseRespCodeRight(t *testing.T, w *httptest.ResponseRecorder, intervals []int) {
	code := w.Code / 100
	for _, v := range intervals {
		if code == v {
			return
		}
	}
	t.Log(w.Code, w.Body.String())
	t.FailNow()
}

func checkFailErr(t *testing.T, e error) {
	if e != nil {
		t.Error(e)
		t.FailNow()
	}
}
