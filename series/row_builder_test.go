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

package metric

import (
	"math"
	"strconv"
	"testing"

	"github.com/cespare/xxhash/v2"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/common/pkg/fasttime"
	"github.com/lindb/common/proto/gen/v1/flatMetricsV1"
)

func Test_NewRowBuilder(t *testing.T) {
	var lastData []byte
	for i := 0; i < 20; i++ {
		rb, releaseFunc := NewRowBuilder()

		assert.NoError(t, rb.AddTag([]byte("a"), []byte("b")))
		assert.NoError(t, rb.AddSimpleField([]byte("f1"), flatMetricsV1.SimpleFieldTypeDeltaSum, 1))
		rb.AddMetricName([]byte("namespace"))
		rb.AddTimestamp(111111)
		thisData, err := rb.Build()
		if i > 0 {
			assert.Equal(t, lastData, thisData)
		}
		lastData = append(lastData[:0], thisData...)
		assert.NoError(t, err)
		releaseFunc(rb)
	}
}

func Test_RowBuilder_ErrorCases(t *testing.T) {
	rb := newRowBuilder()
	// tags validation
	assert.Error(t, rb.AddTag(nil, nil))
	assert.Error(t, rb.AddTag([]byte("tag-key"), nil))

	// simple field validation
	assert.Error(t, rb.AddSimpleField([]byte(""), flatMetricsV1.SimpleFieldTypeDeltaSum, 1))
	assert.Error(t, rb.AddSimpleField([]byte("f1"), flatMetricsV1.SimpleFieldTypeUnSpecified, 1))
	assert.Error(t, rb.AddSimpleField([]byte("f1"), flatMetricsV1.SimpleFieldTypeDeltaSum, math.Inf(1)))
	assert.Error(t, rb.AddSimpleField([]byte("f1"), flatMetricsV1.SimpleFieldTypeDeltaSum, math.NaN()))
	assert.Zero(t, rb.SimpleFieldsLen())

	// compound field validation
	assert.Error(t, rb.AddCompoundFieldData([]float64{1, 2}, []float64{1}))
	assert.Error(t, rb.AddCompoundFieldData([]float64{1}, []float64{1}))
	// not increasing
	assert.Error(t, rb.AddCompoundFieldData(
		[]float64{1, 2, 3, 4, 5, 6},
		[]float64{1, 2, 3, 4, 3, math.Inf(1)},
	))
	// last bound not +Inf
	assert.Error(t, rb.AddCompoundFieldData(
		[]float64{1, 2, 3, 4, 5, 6},
		[]float64{1, 2, 3, 4, 5, 6},
	))
	// first bound < 0
	assert.Error(t, rb.AddCompoundFieldData(
		[]float64{1, 2, 3, 4, 5, 6},
		[]float64{-1, 2, 3, 4, 5, math.Inf(1)},
	))
	// value contains Inf
	assert.Error(t, rb.AddCompoundFieldData(
		[]float64{math.Inf(1), 2, 3, 4, 5, 6},
		[]float64{1, 2, 3, 4, 5, math.Inf(1)},
	))
	// values contains negative float
	assert.Error(t, rb.AddCompoundFieldData(
		[]float64{-1, 2, 3, 4, 5, 6},
		[]float64{1, 2, 3, 4, 5, math.Inf(1)},
	))
	// values contains NaN
	assert.Error(t, rb.AddCompoundFieldData(
		[]float64{math.NaN(), 2, 3, 4, 5, 6},
		[]float64{1, 2, 3, 4, 5, math.Inf(1)},
	))

	assert.NoError(t, rb.AddCompoundFieldData(
		[]float64{1, 2, 3, 4, 5, 6},
		[]float64{1, 2, 3, 4, 5, math.Inf(1)},
	))
	// mmsc
	assert.Error(t, rb.AddCompoundFieldMMSC(-1, -1, 0, 0))
}

func Test_RowBuilder_BuildError(t *testing.T) {
	rb := newRowBuilder()
	_, err := rb.Build()
	assert.Error(t, err)

	// fields empty
	rb.AddMetricName([]byte("ab"))
	rb.AddMetricName([]byte("a|b"))
	rb.AddNameSpace([]byte("a|b"))
	assert.Equal(t, []byte("a_b"), rb.nameSpace)
	assert.Equal(t, []byte("a_b"), rb.metricName)
	_, err = rb.Build()
	assert.Error(t, err)
}

func Test_RowBuilder_OneSimpleField(t *testing.T) {
	rb := newRowBuilder()
	rb.AddMetricName([]byte("cpu"))
	_ = rb.AddSimpleField([]byte("idle"), flatMetricsV1.SimpleFieldTypeLast, 1)
}

func Test_RowBuilder_BuildTo(t *testing.T) {
	rb := newRowBuilder()
	assert.NoError(t, rb.AddTag([]byte("ip"), []byte("1.1.1.1")))
	assert.NoError(t, rb.AddTag([]byte("host"), []byte("dev-ecs")))
	rb.AddMetricName([]byte("cpu|load"))
	assert.NoError(t, rb.AddSimpleField([]byte("idle"), flatMetricsV1.SimpleFieldTypeLast, 1))
	assert.NoError(t, rb.AddCompoundFieldMMSC(1, 1, 1, 1))

	assert.NoError(t, rb.AddCompoundFieldData(
		[]float64{1, 2, 3, 4, 5, 6},
		[]float64{1, 2, 3, 4, 5, math.Inf(1)},
	))
}

func Test_dedupTagsThenXXHash(t *testing.T) {
	rb := newRowBuilder()
	_ = rb.AddTag([]byte("ccc"), []byte("a"))
	_ = rb.AddTag([]byte("d"), []byte("b"))
	_ = rb.AddTag([]byte("a"), []byte("c"))
	_ = rb.AddTag([]byte("ccc"), []byte("d"))
	_ = rb.AddTag([]byte("ccc"), []byte("e"))
	_ = rb.AddTag([]byte("a"), []byte("f"))
	_ = rb.AddTag([]byte("d"), []byte("g"))

	hash1 := rb.dedupTagsThenXXHash()
	assert.Equal(t, "a=f,ccc=e,d=g", rb.hashBuf.String())
	hash2 := rb.dedupTagsThenXXHash()
	assert.Equal(t, "a=f,ccc=e,d=g", rb.hashBuf.String())
	assert.Equal(t, hash2, hash1)
	assert.NotZero(t, hash2)
}

func Test_dedupTags_EmptyKVs(t *testing.T) {
	rb := newRowBuilder()
	hash1 := rb.dedupTagsThenXXHash()
	assert.Equal(t, "", rb.hashBuf.String())
	assert.Equal(t, hash1, emptyStringHash)
}

func Test_dedupTags_SortedKVs(t *testing.T) {
	rb := newRowBuilder()
	_ = rb.AddTag([]byte("a"), []byte("a"))
	_ = rb.AddTag([]byte("c"), []byte("c"))
	_ = rb.dedupTagsThenXXHash()
	assert.Equal(t, "a=a,c=c", rb.hashBuf.String())
}

func Test_dedupTagsThenXXHash_One(t *testing.T) {
	rb := newRowBuilder()
	_ = rb.AddTag([]byte("ccc"), []byte("a"))
	_ = rb.AddTag([]byte("ccc"), []byte("b"))
	_ = rb.AddTag([]byte("ccc"), []byte("c"))
	_ = rb.AddTag([]byte("ccc"), []byte("d"))
	_ = rb.AddTag([]byte("ccc"), []byte("e"))
	_ = rb.AddTag([]byte("ccc"), []byte("f"))
	_ = rb.AddTag([]byte("ccc"), []byte("g"))

	_ = rb.dedupTagsThenXXHash()
	assert.Equal(t, "ccc=g", rb.hashBuf.String())
}

func buildFlatMetric(builder *flatbuffers.Builder) {
	builder.Reset()

	var (
		keys       [10]flatbuffers.UOffsetT
		values     [10]flatbuffers.UOffsetT
		fieldNames [10]flatbuffers.UOffsetT
		kvs        [10]flatbuffers.UOffsetT
		fields     [10]flatbuffers.UOffsetT
	)
	for i := 0; i < 10; i++ {
		keys[i] = builder.CreateString("key" + strconv.Itoa(i))
		values[i] = builder.CreateString("value" + strconv.Itoa(i))
		fieldNames[i] = builder.CreateString("counter" + strconv.Itoa(i))
	}
	for i := 9; i >= 0; i-- {
		flatMetricsV1.KeyValueStart(builder)
		flatMetricsV1.KeyValueAddKey(builder, keys[i])
		flatMetricsV1.KeyValueAddValue(builder, values[i])
		kvs[i] = flatMetricsV1.KeyValueEnd(builder)
	}

	// serialize field names
	for i := 0; i < 10; i++ {
		flatMetricsV1.SimpleFieldStart(builder)
		flatMetricsV1.SimpleFieldAddName(builder, fieldNames[i])
		switch i {
		case 0:
			flatMetricsV1.SimpleFieldAddType(builder, flatMetricsV1.SimpleFieldTypeLast)
		case 1:
			flatMetricsV1.SimpleFieldAddType(builder, flatMetricsV1.SimpleFieldTypeFirst)
		case 2:
			flatMetricsV1.SimpleFieldAddType(builder, flatMetricsV1.SimpleFieldTypeMin)
		case 3:
			flatMetricsV1.SimpleFieldAddType(builder, flatMetricsV1.SimpleFieldTypeMax)
		case 4:
			flatMetricsV1.SimpleFieldAddType(builder, flatMetricsV1.SimpleFieldTypeUnSpecified)
		default:
			flatMetricsV1.SimpleFieldAddType(builder, flatMetricsV1.SimpleFieldTypeDeltaSum)
		}
		flatMetricsV1.SimpleFieldAddValue(builder, float64(i))
		fields[i] = flatMetricsV1.SimpleFieldEnd(builder)
	}

	flatMetricsV1.MetricStartKeyValuesVector(builder, 10)
	for i := 9; i >= 0; i-- {
		builder.PrependUOffsetT(kvs[i])
	}
	kvsAt := builder.EndVector(10)

	flatMetricsV1.MetricStartSimpleFieldsVector(builder, 10)
	for i := 9; i >= 0; i-- {
		builder.PrependUOffsetT(fields[i])
	}
	fieldsAt := builder.EndVector(10)

	// add compound buckets
	flatMetricsV1.CompoundFieldStartValuesVector(builder, 10)
	for i := 9; i >= 0; i-- {
		builder.PrependFloat64(float64(i))
	}
	compoundFieldValues := builder.EndVector(10)
	// add explicit bounds
	flatMetricsV1.CompoundFieldStartExplicitBoundsVector(builder, 10)
	for i := 9; i >= 0; i-- {
		builder.PrependFloat64(float64(i))
	}
	compoundFieldBounds := builder.EndVector(10)
	flatMetricsV1.CompoundFieldStart(builder)
	flatMetricsV1.CompoundFieldAddCount(builder, 1024)
	flatMetricsV1.CompoundFieldAddSum(builder, 1024*1024)
	flatMetricsV1.CompoundFieldAddMin(builder, 10)
	flatMetricsV1.CompoundFieldAddMax(builder, 2048)
	flatMetricsV1.CompoundFieldAddValues(builder, compoundFieldValues)
	flatMetricsV1.CompoundFieldAddExplicitBounds(builder, compoundFieldBounds)
	compoundField := flatMetricsV1.CompoundFieldEnd(builder)

	// serialize metric
	metricName := builder.CreateString("hello")
	namespace := builder.CreateString("default-ns")
	flatMetricsV1.MetricStart(builder)
	flatMetricsV1.MetricAddNamespace(builder, namespace)
	flatMetricsV1.MetricAddName(builder, metricName)
	flatMetricsV1.MetricAddTimestamp(builder, fasttime.UnixMilliseconds())
	flatMetricsV1.MetricAddKeyValues(builder, kvsAt)
	flatMetricsV1.MetricAddHash(builder, xxhash.Sum64String("hello"))
	flatMetricsV1.MetricAddSimpleFields(builder, fieldsAt)
	flatMetricsV1.MetricAddCompoundField(builder, compoundField)

	end := flatMetricsV1.MetricEnd(builder)
	builder.Finish(end)
}
