package types

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type DateOnly time.Time

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
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
