// Package exporter provides utilities for serializing cron schedule data
// into multiple output formats.
//
// Supported formats:
//
//   - JSON  — structured payload via ToJSON / NewPayload
//   - CSV   — tabular export via ToCSV
//   - iCal  — RFC 5545 calendar format via ToICAL
//
// Each exporter accepts a slice of scheduler.ScheduledRun values and optional
// conflict data produced by the scheduler package, enabling downstream tools
// such as calendar applications, dashboards, and data pipelines to consume
// cron schedule information without depending on crontrace internals.
package exporter
