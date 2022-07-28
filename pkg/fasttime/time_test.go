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

package fasttime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_FastTime(t *testing.T) {
	time.Sleep(time.Millisecond * 100)

	assert.True(t, time.Now().UnixNano()-UnixNano() < 20e6)
	assert.True(t, time.Now().UnixNano()/1e3-UnixMicroseconds() < 20e3)
	assert.True(t, time.Now().UnixNano()/1e6-UnixMilliseconds() < 20)
	assert.True(t, float64(time.Now().UnixNano()/1e9)-float64(UnixTimestamp()) < 0.02)
}
