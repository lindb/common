package models

// Exemplar represents an exemplar, which is a sample input measurement.
// Exemplars also hold information about the environment when the measurement
// was recorded, for example the span and trace ID of the active span when the
// exemplar was recorded.
type Exemplar struct {
	TraceID  string `json:"traceId"`
	SpanID   string `json:"spanId"`
	Duration int64  `json:"duration"`
}
