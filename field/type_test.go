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

package field

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestType_String(t *testing.T) {
	assert.Equal(t, "unknown", Unknown.String())
	assert.Equal(t, "sum", Sum.String())
	assert.Equal(t, "min", Min.String())
	assert.Equal(t, "max", Max.String())
	assert.Equal(t, "last", Last.String())
	assert.Equal(t, "first", First.String())
	assert.Equal(t, "histogram", Histogram.String())
	assert.Equal(t, "exemplar", Exemplar.String())
	// unrecognized value falls back to "unknown"
	assert.Equal(t, "unknown", Type(99).String())
}

func TestType_IsExemplar(t *testing.T) {
	assert.True(t, Exemplar.IsExemplar())
	assert.False(t, Sum.IsExemplar())
	assert.False(t, Unknown.IsExemplar())
}

// TestType_NumericValues verifies the stable numeric values that are persisted on disk
// and carried over the Arrow IPC wire. Changing these would be a breaking change.
func TestType_NumericValues(t *testing.T) {
	assert.Equal(t, Type(0), Unknown)
	assert.Equal(t, Type(1), Sum)
	assert.Equal(t, Type(2), Min)
	assert.Equal(t, Type(3), Max)
	assert.Equal(t, Type(4), Last)
	assert.Equal(t, Type(5), First)
	assert.Equal(t, Type(6), Histogram)
	assert.Equal(t, Type(7), Exemplar)
}

func TestType_Aggregate(t *testing.T) {
	assert.Equal(t, 100.0, Sum.Aggregate(1, 99.0))
	assert.Equal(t, 100.0, Histogram.Aggregate(1, 99.0), "Histogram uses sum aggregation")

	assert.Equal(t, 1.0, Min.Aggregate(1, 99.0))
	assert.Equal(t, 1.0, Min.Aggregate(99.0, 1))

	assert.Equal(t, 99.0, Max.Aggregate(1, 99.0))
	assert.Equal(t, 99.0, Max.Aggregate(99.0, 1))

	assert.Equal(t, 99.0, Last.Aggregate(1, 99.0), "last returns b")
	assert.Equal(t, 1.0, Last.Aggregate(99.0, 1), "last returns b")

	assert.Equal(t, 1.0, First.Aggregate(1, 99.0), "first returns a")
	assert.Equal(t, 99.0, First.Aggregate(99.0, 1), "first returns a")
}

func TestType_Aggregate_Panics(t *testing.T) {
	assert.Panics(t, func() { Exemplar.Aggregate(1, 2) })
	assert.Panics(t, func() { Unknown.Aggregate(1, 2) })
	assert.Panics(t, func() { Type(99).Aggregate(1, 2) })
}
