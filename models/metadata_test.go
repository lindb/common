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

package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetadata_ToTable(t *testing.T) {
	rows, rs := (&Metadata{}).ToTable()
	assert.Zero(t, rows)
	assert.Empty(t, rs)

	rows, rs = (&Metadata{Type: "namespace", Values: []interface{}{"name"}}).ToTable()
	assert.Equal(t, rows, 1)
	assert.NotEmpty(t, rs)
	rows, rs = (&Metadata{Type: "metric", Values: []interface{}{"name"}}).ToTable()
	assert.Equal(t, rows, 1)
	assert.NotEmpty(t, rs)
	rows, rs = (&Metadata{Type: "tagKey", Values: []interface{}{"name"}}).ToTable()
	assert.Equal(t, rows, 1)
	assert.NotEmpty(t, rs)
	rows, rs = (&Metadata{Type: "tagValue", Values: []interface{}{"name"}}).ToTable()
	assert.Equal(t, rows, 1)
	assert.NotEmpty(t, rs)
	rows, rs = (&Metadata{
		Type:   "field",
		Values: []interface{}{map[string]interface{}{"name": "n", "Type": "sum"}},
	}).ToTable()
	assert.Equal(t, rows, 1)
	assert.NotEmpty(t, rs)
}
