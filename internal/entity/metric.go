package entity

import (
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
)

var defaultMetricTypes = []string{"INT", "FLOAT", "DURATION", "TIMESTAMP_WITH_TIMEZONE", "BOOL", "STRING"}

type Metric struct {
	MetricID   int    `json:"metric_id"`
	Slug       string `json:"slug"`
	MetricType string `json:"metric_type"`
	Details    string `json:"details"`
}

type AddMetric struct {
	MetricID    int         `json:"metric_id"`
	MetricValue interface{} `json:"metric_value"`
}

type GetMetric struct {
	TimeStamp CustomTime  `json:"time_stamp"`
	Value     interface{} `json:"value"`
}

func (m *Metric) Validate() error {
	m.Slug = strings.Join(strings.Fields(m.Slug), "_")
	m.Slug = strings.ToUpper(m.Slug)

	return validation.ValidateStruct(
		m,
		validation.Field(
			&m.Slug,
			validation.Required,
			validation.Match(regexp.MustCompile(`^[\w]+$`)),
			validation.Length(0, 255),
		),
		validation.Field(
			&m.MetricType,
			validation.Required,
			validation.Match(regexp.MustCompile(strings.Join(defaultMetricTypes[:], "|"))),
			validation.Length(0, 255),
		),
		validation.Field(
			&m.Details,
			validation.Required,
			validation.Length(0, 255),
		),
	)
}
