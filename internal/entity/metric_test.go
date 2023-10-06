package entity_test

import (
	"testing"

	"github.com/AnatoliyBr/dwh-service/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestMetric_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		m       func() *entity.Metric
		isValid bool
	}{
		{
			name: "valid",
			m: func() *entity.Metric {
				return entity.TestMetric(t)
			},
			isValid: true,
		},
		{
			name: "mixedcase with whitespace in slug",
			m: func() *entity.Metric {
				m := entity.TestMetric(t)
				m.Slug = "reading_TIME note 1 "
				return m
			},
			isValid: true,
		},
		{
			name: "empty slug",
			m: func() *entity.Metric {
				m := entity.TestMetric(t)
				m.Slug = ""
				return m
			},
			isValid: false,
		},
		{
			name: "invalid symbols in slug",
			m: func() *entity.Metric {
				m := entity.TestMetric(t)
				m.Slug = "READING_?#@*&%!"
				return m
			},
			isValid: false,
		},
		{
			name: "long slug",
			m: func() *entity.Metric {
				m := entity.TestMetric(t)
				m.Slug = `FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF \
				FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF \
				FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF \
				FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF \
				FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF \
				FFFFFF`
				return m
			},
			isValid: false,
		},
		{
			name: "invalid type",
			m: func() *entity.Metric {
				m := entity.TestMetric(t)
				m.MetricType = "TIMESTAMP"
				return m
			},
			isValid: false,
		},
		{
			name: "empty type",
			m: func() *entity.Metric {
				m := entity.TestMetric(t)
				m.MetricType = ""
				return m
			},
			isValid: false,
		},
		{
			name: "long type",
			m: func() *entity.Metric {
				m := entity.TestMetric(t)
				m.MetricType = `FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF \
				FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF \
				FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF \
				FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF \
				FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF \
				FFFFFF`
				return m
			},
			isValid: false,
		},
		{
			name: "empty details",
			m: func() *entity.Metric {
				m := entity.TestMetric(t)
				m.Details = ""
				return m
			},
			isValid: false,
		},
		{
			name: "long deatails",
			m: func() *entity.Metric {
				m := entity.TestMetric(t)
				m.Details = `FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF \
				FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF \
				FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF \
				FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF \
				FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF \
				FFFFFF`
				return m
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.m().Validate())
			} else {
				assert.Error(t, tc.m().Validate())
			}
		})
	}
}
