package entity

import (
	"bytes"
	"time"
)

const (
	defaultLayout = time.RFC3339
)

type CustomTime struct {
	time.Time
}

func (t *CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(t.Format(defaultLayout)), nil
}

func (t *CustomTime) UnMarshalJSON(data []byte) error {

	// from json doc: by convention, unmarshalers implement UnmarshalJSON([]byte("null")) as a no-op
	if bytes.Equal(data, []byte("null")) {
		return nil
	}

	time, err := time.Parse(defaultLayout, string(data))
	if err != nil {
		return err
	}

	t = &CustomTime{time}
	return nil
}
