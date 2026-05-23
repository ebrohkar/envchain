// Package labeler provides a mechanism for attaching arbitrary metadata
// labels to environment variable keys within envchain pipelines.
//
// Labels are key-value pairs (e.g. tier=critical, owner=platform) that
// can be used to annotate, query, and filter variables during validation
// or reporting stages.
//
// Example usage:
//
//	 l := labeler.New()
//	 l.Attach("DB_HOST", "tier", "critical")
//	 l.Attach("DB_HOST", "owner", "platform")
//
//	 keys := l.FindByLabel("tier", "critical")
//	 // keys => ["DB_HOST"]
package labeler
