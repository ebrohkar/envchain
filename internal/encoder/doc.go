// Package encoder provides utilities for encoding environment variable values
// into common formats such as base64, hex, and URL-encoded strings.
//
// It supports selective encoding — only a specified subset of keys can be
// targeted, leaving all other variables untouched. This is useful when
// preparing configs for systems that require encoded secrets or tokens.
//
// Supported formats:
//   - FormatBase64: standard base64 encoding
//   - FormatHex:    lowercase hexadecimal representation
//   - FormatURL:    percent-encoded URL-safe string
//
// Example:
//
//	enc := encoder.New(encoder.FormatBase64, "DB_PASSWORD", "API_SECRET")
//	result, err := enc.Apply(env)
package encoder
