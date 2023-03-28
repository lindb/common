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
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRecovery(t *testing.T) {
	r := gin.New()
	r.GET("/panic", func(context *gin.Context) {
		panic("panic")
	})
	assert.Panics(t, func() {
		_ = DoRequest(t, r, http.MethodGet, "/panic", "")
	})

	r = gin.New()
	r.Use(Recovery())
	r.GET("/panic", func(context *gin.Context) {
		panic("panic")
	})
	resp := DoRequest(t, r, http.MethodGet, "/panic", "")
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}
