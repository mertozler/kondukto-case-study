package models

import "time"

type ScanData struct {
	ScanID   string
	ScanData ScanDatas
}

type ScanDatas struct {
	Errors      []interface{}          `json:"errors"`
	GeneratedAt time.Time              `json:"generated_at"`
	Metrics     map[string]MetricsData `json:"metrics"`
	Results     []Results              `json:"results"`
}

type MetricsData struct {
	CONFIDENCEHIGH      float64 `json:"CONFIDENCE.HIGH"`
	CONFIDENCELOW       float64 `json:"CONFIDENCE.LOW"`
	CONFIDENCEMEDIUM    float64 `json:"CONFIDENCE.MEDIUM"`
	CONFIDENCEUNDEFINED float64 `json:"CONFIDENCE.UNDEFINED"`
	SEVERITYHIGH        float64 `json:"SEVERITY.HIGH"`
	SEVERITYLOW         float64 `json:"SEVERITY.LOW"`
	SEVERITYMEDIUM      float64 `json:"SEVERITY.MEDIUM"`
	SEVERITYUNDEFINED   float64 `json:"SEVERITY.UNDEFINED"`
	Loc                 int     `json:"loc"`
	Nosec               int     `json:"nosec"`
}

type Results struct {
	Code            string `json:"code"`
	Filename        string `json:"filename"`
	IssueConfidence string `json:"issue_confidence"`
	IssueSeverity   string `json:"issue_severity"`
	IssueText       string `json:"issue_text"`
	LineNumber      int    `json:"line_number"`
	LineRange       []int  `json:"line_range"`
	MoreInfo        string `json:"more_info"`
	TestID          string `json:"test_id"`
	TestName        string `json:"test_name"`
}
