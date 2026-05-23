package types

import (
	"database/sql/driver"
	"fmt"
	"log"
	"strings"
	"time"
)

type DateOnly time.Time

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	if len(b) == 0 || string(b) == "null" {
		return nil
	}

	// Gin Query Start with string
	if b[0] != '"' {
		return d.UnmarshalText(b)
	}

	// Gin Body string in string
	s := strings.Trim(string(b), "\"")
	if s == "" {
		return nil
	}

	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}

	*d = DateOnly(t)
	return nil
}

func (d DateOnly) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, "\"%s\"", time.Time(d).Format("2006-01-02")), nil
}

func (d DateOnly) Value() (driver.Value, error) {
	return time.Time(d).Format("2006-01-02"), nil
}

func (d *DateOnly) Scan(value any) error {
	if value == nil {
		return nil
	}
	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("failed to scan DateOnly")
	}
	*d = DateOnly(t)
	return nil
}

func (d DateOnly) IsZero() bool {
	return time.Time(d).IsZero()
}

func (d *DateOnly) UnmarshalText(text []byte) error {
	s := string(text)
	log.Println(s)
	if s == "" || s == "null" {
		return nil
	}

	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}

	*d = DateOnly(t)
	return nil
}

func ParseString(s string) (DateOnly, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return DateOnly{}, err
	}

	return DateOnly(t), nil
}

func (d *DateOnly) String() string {
	return time.Time(*d).Format("2006-01-02")
}
