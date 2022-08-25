package models

type RunMetrics struct {
	Severity   Metric
	Confidence Metric
}

type Metric struct {
	Undefined float32
	Low       float32
	Medium    float32
	High      float32
}
