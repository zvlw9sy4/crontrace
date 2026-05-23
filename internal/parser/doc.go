// Package parser provides functionality for parsing and validating
// standard 5-field cron expressions used by crontrace.
//
// A cron expression has the following structure:
//
//	┌───────────── minute        (0 - 59)
//	│ ┌─────────── hour          (0 - 23)
//	│ │ ┌───────── day of month  (1 - 31)
//	│ │ │ ┌─────── month         (1 - 12)
//	│ │ │ │ ┌───── day of week   (0 - 7, Sunday = 0 or 7)
//	│ │ │ │ │
//	* * * * *
//
// Supported syntax per field:
//   - *        : any value
//   - n        : exact value
//   - n,m,...  : list of values
//   - n-m      : range of values
//   - */step   : every step units
//   - n-m/step : every step units within range
//
// Example usage:
//
//	expr, err := parser.Parse("*/15 9-17 * * 1-5")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(expr.Hour) // "9-17"
package parser
