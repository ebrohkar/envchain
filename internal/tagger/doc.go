// Package tagger provides a mechanism for attaching arbitrary key-value
// metadata tags to environment variable names.
//
// Tags can be used to annotate variables with contextual information such
// as environment tier ("production", "staging"), owning team, sensitivity
// level, or any other classification relevant to deployment workflows.
//
// Example usage:
//
//	tr := tagger.New()
//	tr.Tag("DB_PASSWORD", "sensitivity", "high")
//	tr.Tag("DB_PASSWORD", "team", "platform")
//
//	tags := tr.Get("DB_PASSWORD")
//	matches := tr.FindByTag("sensitivity", "high")
package tagger
