package entity

import (
	"testing"
	"time"
)

func TestService(t *testing.T) *Service {
	return &Service{
		Slug: "NOTE_BOOK",
		Details: `NoteBook is the best word processing app for all your works, \
				from taking down quick notes to writing your books, \
				eBooks and organizing your documents. This app is available for iOS and Mac devices.`,
	}
}

func TestMetric(t *testing.T) *Metric {
	return &Metric{
		Slug:       "READING_TIME_NOTE_1",
		MetricType: "DURATION_IN_NS",
		Details: `The "Read Time" metric allows you to estimate the approximate amount \
				of time it will take the user to read the page from beginning to end, \
				including the content of all snippets, variables, headers and footers, if any.`,
	}
}

func TestEvent(t *testing.T) *Event {
	return &Event{
		TimeStamp: CustomTime{time.Now()}}
}
