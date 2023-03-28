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

package timeutil

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const date = "20191212 10:11:10"

func Test_ParseTimestamp(t *testing.T) {
	defer func() {
		parseTimeFunc = time.ParseInLocation
	}()
	_, err := ParseTimestamp(date)
	assert.Nil(t, err)

	_, err = ParseTimestamp(date)
	assert.Nil(t, err)

	_, err = ParseTimestamp(date, DataTimeFormat1, DataTimeFormat2)
	assert.Nil(t, err)
	_, err = ParseTimestamp("2019-12-12 10:11:10")
	assert.Nil(t, err)
	_, err = ParseTimestamp("2019/12/12 10:11:10")
	assert.Nil(t, err)
	_, err = ParseTimestamp("20191212101110")
	assert.Nil(t, err)

	parseTimeFunc = func(layout, value string, loc *time.Location) (t time.Time, err error) {
		return time.Now(), fmt.Errorf("err")
	}
	_, err = ParseTimestamp(date)
	assert.Error(t, err)
}

func Test_FormatTimestamp(t *testing.T) {
	now := Now()
	fmt.Println(FormatTimestamp(now, DataTimeFormat2))
}
