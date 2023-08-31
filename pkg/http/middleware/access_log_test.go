// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package middleware

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/common/pkg/logger"
)

func TestAccessLogMiddleware(t *testing.T) {
	defer func() {
		pathUnescapeFunc = url.PathUnescape
	}()

	r := gin.New()
	r.Use(AccessLog(logger.GetLogger(logger.AccessLogModule, "HTTP")))
	r.GET("/home", func(c *gin.Context) {
		_ = c.Error(fmt.Errorf("err"))
		c.JSON(http.StatusOK, "ok")
	})

	pathUnescapeFunc = func(s string) (string, error) {
		return "err-path", fmt.Errorf("err")
	}
	_ = DoRequest(t, r, http.MethodPut, "/test", `{"username": "admin", "password": "admin123"}`)

	pathUnescapeFunc = url.PathUnescape
	_ = DoRequest(t, r, http.MethodGet, "/home", `{"username": "admin", "password": "admin123"}`)
}

func Test_real_ip(t *testing.T) {
	req, _ := http.NewRequestWithContext(context.TODO(), "GET", "/health-check", bytes.NewReader([]byte("test")))
	req.Header.Add("X-Real-Ip", "real-ip")
	assert.Equal(t, "real-ip", realIP(req))

	req, _ = http.NewRequestWithContext(context.TODO(), "GET", "/health-check", bytes.NewReader([]byte("test")))
	req.Header.Add("X-Forwarded-For", "forward-ip")
	assert.Equal(t, "forward-ip", realIP(req))
	req, _ = http.NewRequestWithContext(context.TODO(), "GET", "/health-check", bytes.NewReader([]byte("test")))
	req.RemoteAddr = "1.1.1.1:1023"
	assert.Equal(t, "1.1.1.1", realIP(req))
}

// DoRequest does http request for test.
func DoRequest(t *testing.T, r *gin.Engine, method, path, reqBody string, headers ...http.Header) *httptest.ResponseRecorder {
	t.Helper()

	var body io.Reader
	if reqBody != "" {
		body = bytes.NewBufferString(reqBody)
	}
	req, _ := http.NewRequestWithContext(context.TODO(), method, path, body)
	if len(headers) == 0 {
		req.Header.Set("content-type", "application/json")
	} else {
		for _, header := range headers {
			req.Header = header
		}
	}
	resp := newCloseNotifyingRecorder()
	r.ServeHTTP(resp, req)
	return resp.ResponseRecorder
}

type closeNotifyingRecorder struct {
	*httptest.ResponseRecorder
	closed chan bool
}

func newCloseNotifyingRecorder() *closeNotifyingRecorder {
	return &closeNotifyingRecorder{
		httptest.NewRecorder(),
		make(chan bool, 1),
	}
}

func (c *closeNotifyingRecorder) CloseNotify() <-chan bool {
	return c.closed
}
