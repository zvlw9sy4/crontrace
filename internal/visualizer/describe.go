package visualizer

import "strings"

// Describe returns a human-readable description for common cron expressions.
// Falls back to a generic label for unrecognised patterns.
func Describe(expr string) string {
	norm := strings.TrimSpace(expr)

	switch norm {
	case "* * * * *":
		return "Every minute"
	case "0 * * * *":
		return "Every hour"
	case "0 0 * * *":
		return "Daily at midnight"
	case "0 12 * * *":
		return "Daily at noon"
	case "0 0 * * 0":
		return "Weekly on Sunday"
	case "0 0 * * 1":
		return "Weekly on Monday"
	case "0 0 1 * *":
		return "Monthly on 1st"
	case "0 0 1 1 *":
		return "Yearly on Jan 1st"
	case "*/5 * * * *":
		return "Every 5 minutes"
	case "*/10 * * * *":
		return "Every 10 minutes"
	case "*/15 * * * *":
		return "Every 15 minutes"
	case "*/30 * * * *":
		return "Every 30 minutes"
	case "0 */2 * * *":
		return "Every 2 hours"
	case "0 9-17 * * 1-5":
		return "Hourly, business hours Mon-Fri"
	}

	fields := strings.Fields(norm)
	if len(fields) != 5 {
		return "Custom schedule"
	}

	if fields[0] != "*" && strings.HasPrefix(fields[0], "*/") {
		return "Repeating minute interval"
	}
	if fields[1] != "*" && strings.HasPrefix(fields[1], "*/") {
		return "Repeating hour interval"
	}

	return "Custom schedule"
}
