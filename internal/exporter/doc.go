// Package exporter provides serialization utilities for crontrace schedule data.
//
// It supports exporting parsed cron schedules, next-run predictions, and
// detected conflicts into structured formats suitable for programmatic
// consumption or integration with external tooling.
//
// # JSON Export
//
// The primary export format is JSON. Use [NewPayload] to construct an
// [ExportPayload] from schedule and conflict data, then call [ToJSON] to
// serialize it:
//
//	payload := exporter.NewPayload(schedules, conflicts)
//	data, err := exporter.ToJSON(payload)
//
// The resulting JSON includes a top-level generated_at timestamp (UTC),
// a schedules array, and a conflicts array.
package exporter
