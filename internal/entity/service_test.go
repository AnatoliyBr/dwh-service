package entity_test

import (
	"testing"

	"github.com/AnatoliyBr/dwh-service/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestService_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		s       func() *entity.Service
		isValid bool
	}{
		{
			name: "valid",
			s: func() *entity.Service {
				return entity.TestService(t)
			},
			isValid: true,
		},
		{
			name: "mixedcase with whitespace in slug",
			s: func() *entity.Service {
				s := entity.TestService(t)
				s.Slug = "note_BOOK  v 1 "
				return s
			},
			isValid: true,
		},
		{
			name: "empty slug",
			s: func() *entity.Service {
				s := entity.TestService(t)
				s.Slug = ""
				return s
			},
			isValid: false,
		},
		{
			name: "invalid symbols in slug",
			s: func() *entity.Service {
				s := entity.TestService(t)
				s.Slug = "NOTE_?#@*&%!"
				return s
			},
			isValid: false,
		},
		{
			name: "long slug",
			s: func() *entity.Service {
				s := entity.TestService(t)
				s.Slug = `FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF \
				FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF \
				FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF \
				FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF \
				FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF \
				FFFFFF`
				return s
			},
			isValid: false,
		},
		{
			name: "empty details",
			s: func() *entity.Service {
				s := entity.TestService(t)
				s.Details = ""
				return s
			},
			isValid: false,
		},
		{
			name: "long deatails",
			s: func() *entity.Service {
				s := entity.TestService(t)
				s.Details = `FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF \
				FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF \
				FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF \
				FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF \
				FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF \
				FFFFFF`
				return s
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.s().Validate())
			} else {
				assert.Error(t, tc.s().Validate())
			}
		})
	}
}
