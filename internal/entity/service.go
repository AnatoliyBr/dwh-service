package entity

import (
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Service struct {
	ServiceID int    `json:"service_id"`
	Slug      string `json:"slug"`
	Details   string `json:"details"`
}

func (s *Service) Validate() error {
	s.Slug = strings.Join(strings.Fields(s.Slug), "_")
	s.Slug = strings.ToUpper(s.Slug)

	return validation.ValidateStruct(
		s,
		validation.Field(
			&s.Slug,
			validation.Required,
			validation.Match(regexp.MustCompile(`^[\w]+$`)),
			validation.Length(0, 255),
		),
		validation.Field(
			&s.Details,
			validation.Required,
			validation.Length(0, 255),
		),
	)
}
