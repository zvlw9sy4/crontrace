// Package exporter provides functionality to export cron schedule data
// in various formats including JSON, CSV, and iCalendar.
package exporter

import (
	"fmt"
	"strings"
	"time"

	"github.com/user/crontrace/internal/scheduler"
)

const icalDateFormat = "20060102T150405Z"

// ToICAL converts a list of scheduled next-run times into an iCalendar (.ics)
// formatted string. Each next-run entry becomes a VEVENT block.
func ToICAL(expression string, runs []scheduler.ScheduledRun) (string, error) {
	if expression == "" {
		return "", fmt.Errorf("expression must not be empty")
	}

	var sb strings.Builder

	sb.WriteString("BEGIN:VCALENDAR\r\n")
	sb.WriteString("VERSION:2.0\r\n")
	sb.WriteString("PRODID:-//crontrace//EN\r\n")
	sb.WriteString("CALSCALE:GREGORIAN\r\n")

	for i, run := range runs {
		uid := fmt.Sprintf("crontrace-%s-%d@crontrace", sanitizeUID(expression), i)
		stamp := time.Now().UTC().Format(icalDateFormat)
		start := run.Time.UTC().Format(icalDateFormat)
		// Treat each cron event as a 1-minute block.
		end := run.Time.UTC().Add(time.Minute).Format(icalDateFormat)

		sb.WriteString("BEGIN:VEVENT\r\n")
		fmt.Fprintf(&sb, "UID:%s\r\n", uid)
		fmt.Fprintf(&sb, "DTSTAMP:%s\r\n", stamp)
		fmt.Fprintf(&sb, "DTSTART:%s\r\n", start)
		fmt.Fprintf(&sb, "DTEND:%s\r\n", end)
		fmt.Fprintf(&sb, "SUMMARY:Cron: %s\r\n", expression)
		fmt.Fprintf(&sb, "DESCRIPTION:Scheduled run #%d for expression %s\r\n", i+1, expression)
		sb.WriteString("END:VEVENT\r\n")
	}

	sb.WriteString("END:VCALENDAR\r\n")
	return sb.String(), nil
}

// sanitizeUID replaces characters unsafe for iCalendar UIDs with hyphens.
func sanitizeUID(expr string) string {
	replacer := strings.NewReplacer(
		" ", "-",
		"*", "star",
		"/", "sl",
		",", "c",
	)
	return replacer.Replace(expr)
}
