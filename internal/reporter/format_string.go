package reporter

import "fmt"

// String returns the string representation of the Format.
func (f Format) String() string {
	return string(f)
}

// ParseFormat parses a string into a Format value.
// Returns an error if the format is not recognized.
func ParseFormat(s string) (Format, error) {
	switch Format(s) {
	case FormatText:
		return FormatText, nil
	case FormatJSON:
		return FormatJSON, nil
	default:
		return "", fmt.Errorf("unknown report format %q: must be one of [text, json]", s)
	}
}

// SupportedFormats returns all supported Format values.
func SupportedFormats() []Format {
	return []Format{FormatText, FormatJSON}
}
