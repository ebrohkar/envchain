// Package auditor provides a lightweight audit trail for envchain pipeline runs.
//
// It records discrete events — such as variable resolution, missing keys, and
// validation outcomes — with timestamps and structured metadata. The audit log
// can be inspected programmatically or written as a human-readable summary.
//
// Usage:
//
//	a := auditor.New(os.Stderr)
//	a.Record(auditor.EventResolved, "DB_HOST", "resolved from environment")
//	a.Record(auditor.EventMissing, "API_KEY", "not set")
//	a.WriteSummary()
//	if a.HasFailures() {
//		os.Exit(1)
//	}
package auditor
