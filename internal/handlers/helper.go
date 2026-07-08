package handlers

import (
	"github.com/The-Ogulgozel/Banking-system/internal/errors"
	"time"
)

func mustParseDate(s *string) (*time.Time, error) {
	if s == nil || *s == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", *s)
	if err != nil {
		return nil, errors.NewAppError(nil, "invalid date format", "SE-00400")
	}
	return &t, nil
}
