package encoder

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// Format represents the encoding format to apply to environment variable values.
type Format int

const (
	FormatBase64 Format = iota
	FormatHex
	FormatURL
)

// Encoder applies an encoding transformation to environment variable values.
type Encoder struct {
	format Format
	keys   map[string]struct{}
}

// New returns an encoder that applies the given format to all matching keys.
// If no keys are provided, all variables are encoded.
func New(format Format, keys ...string) *Encoder {
	km := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		km[strings.ToUpper(k)] = struct{}{}
	}
	return &Encoder{format: format, keys: km}
}

// Apply encodes values in the provided env map and returns a new map.
func (e *Encoder) Apply(env map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if e.shouldEncode(k) {
			encoded, err := e.encode(v)
			if err != nil {
				return nil, fmt.Errorf("encoder: failed to encode %q: %w", k, err)
			}
			out[k] = encoded
		} else {
			out[k] = v
		}
	}
	return out, nil
}

func (e *Encoder) shouldEncode(key string) bool {
	if len(e.keys) == 0 {
		return true
	}
	_, ok := e.keys[strings.ToUpper(key)]
	return ok
}

func (e *Encoder) encode(value string) (string, error) {
	switch e.format {
	case FormatBase64:
		return base64.StdEncoding.EncodeToString([]byte(value)), nil
	case FormatHex:
		var sb strings.Builder
		for _, b := range []byte(value) {
			fmt.Fprintf(&sb, "%02x", b)
		}
		return sb.String(), nil
	case FormatURL:
		return urlEscape(value), nil
	default:
		return "", fmt.Errorf("unsupported format: %d", e.format)
	}
}

func urlEscape(s string) string {
	var sb strings.Builder
	for _, b := range []byte(s) {
		if isURLSafe(b) {
			sb.WriteByte(b)
		} else {
			fmt.Fprintf(&sb, "%%%02X", b)
		}
	}
	return sb.String()
}

func isURLSafe(b byte) bool {
	return (b >= 'A' && b <= 'Z') ||
		(b >= 'a' && b <= 'z') ||
		(b >= '0' && b <= '9') ||
		b == '-' || b == '_' || b == '.' || b == '~'
}
