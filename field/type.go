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

import "math"

// Type represents the aggregation kind of a metric field.
// It is the single canonical definition shared by both the Arrow IPC transport
// layer (github.com/lindb/arrow) and the LinDB storage engine.
// The numeric values are stable and stored on disk — do not reorder.
type Type uint8

// Defines all field types.
// Values 1-5 are used in both LinDB-internal Arrow IPC and the external wire format.
// Values 6-7 are LinDB storage-layer extensions that also appear in LinDB-to-LinDB IPC.
const (
	Unknown   Type = iota // 0
	Sum                   // 1 — sum-aggregated field
	Min                   // 2 — min-aggregated field
	Max                   // 3 — max-aggregated field
	Last                  // 4 — last-value field (gauge)
	First                 // 5 — first-value field
	Histogram             // 6 — histogram bucket (alias for sum, identified by field name convention)
	Exemplar              // 7 — trace exemplar field
)

// IsExemplar reports whether the field type carries trace exemplar data.
func (t Type) IsExemplar() bool {
	return t == Exemplar
}

// String returns the human-readable name of the field type.
func (t Type) String() string {
	switch t {
	case Sum:
		return "sum"
	case Min:
		return "min"
	case Max:
		return "max"
	case Last:
		return "last"
	case First:
		return "first"
	case Histogram:
		return "histogram"
	case Exemplar:
		return "exemplar"
	default:
		return "unknown"
	}
}

// Aggregate combines two float64 field values using this type's aggregation rule.
// Panics for Exemplar, Unknown, and any unrecognized type —
// use field.ExemplarAggregate for exemplar merging.
func (t Type) Aggregate(a, b float64) float64 {
	switch t {
	case Sum, Histogram:
		return a + b
	case Min:
		return math.Min(a, b)
	case Max:
		return math.Max(a, b)
	case Last:
		return b
	case First:
		return a
	default:
		panic("field.Type.Aggregate: unsupported field type")
	}
}
