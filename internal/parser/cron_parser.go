package parser

import (
	"fmt"
	"strconv"
	"strings"
)

// CronExpression holds the parsed fields of a cron expression.
type CronExpression struct {
	Raw     string
	Minute  string
	Hour    string
	Day     string
	Month   string
	Weekday string
}

// Parse parses a standard 5-field cron expression string.
func Parse(expr string) (*CronExpression, error) {
	expr = strings.TrimSpace(expr)
	fields := strings.Fields(expr)
	if len(fields) != 5 {
		return nil, fmt.Errorf("invalid cron expression %q: expected 5 fields, got %d", expr, len(fields))
	}

	c := &CronExpression{
		Raw:     expr,
		Minute:  fields[0],
		Hour:    fields[1],
		Day:     fields[2],
		Month:   fields[3],
		Weekday: fields[4],
	}

	if err := validateField(c.Minute, 0, 59, "minute"); err != nil {
		return nil, err
	}
	if err := validateField(c.Hour, 0, 23, "hour"); err != nil {
		return nil, err
	}
	if err := validateField(c.Day, 1, 31, "day"); err != nil {
		return nil, err
	}
	if err := validateField(c.Month, 1, 12, "month"); err != nil {
		return nil, err
	}
	if err := validateField(c.Weekday, 0, 7, "weekday"); err != nil {
		return nil, err
	}

	return c, nil
}

// validateField checks that a single cron field is syntactically valid.
func validateField(field string, min, max int, name string) error {
	if field == "*" {
		return nil
	}
	// Handle step values like */5 or 1-5/2
	parts := strings.Split(field, "/")
	if len(parts) > 2 {
		return fmt.Errorf("invalid %s field %q: too many slashes", name, field)
	}
	if len(parts) == 2 {
		step, err := strconv.Atoi(parts[1])
		if err != nil || step < 1 {
			return fmt.Errorf("invalid %s field %q: bad step value", name, field)
		}
	}
	base := parts[0]
	if base == "*" {
		return nil
	}
	// Handle ranges like 1-5
	rangeParts := strings.Split(base, "-")
	if len(rangeParts) == 2 {
		lo, err1 := strconv.Atoi(rangeParts[0])
		hi, err2 := strconv.Atoi(rangeParts[1])
		if err1 != nil || err2 != nil || lo < min || hi > max || lo > hi {
			return fmt.Errorf("invalid %s field %q: range out of bounds [%d-%d]", name, field, min, max)
		}
		return nil
	}
	// Handle lists like 1,2,3
	for _, v := range strings.Split(base, ",") {
		n, err := strconv.Atoi(v)
		if err != nil || n < min || n > max {
			return fmt.Errorf("invalid %s field %q: value %q out of range [%d-%d]", name, field, v, min, max)
		}
	}
	return nil
}
