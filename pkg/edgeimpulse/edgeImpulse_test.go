package edgeimpulse

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_validateResp(t *testing.T) {
	tests := []struct {
		name       string
		code       int
		body       string
		wantErrStr string
	}{
		{name: "200", code: 200, body: "OK", wantErrStr: ""},
		{name: "299", code: 299, body: "OK", wantErrStr: ""},
		{name: "400", code: 400, body: "error", wantErrStr: "resp code: 400\nerror"},
		{name: "503", code: 503, body: "error", wantErrStr: "resp code: 503\nerror"},
		{name: "503 no err", code: 503, body: "", wantErrStr: "resp code: 503"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tResp := httptest.NewRecorder()
			tResp.Body = bytes.NewBuffer([]byte(tt.body))
			tResp.Code = tt.code
			err := validateResp(tResp.Result())
			if tt.wantErrStr != "" {
				assert.Equal(t, tt.wantErrStr, err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
