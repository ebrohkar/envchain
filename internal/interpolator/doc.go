// Package interpolator resolves template-style variable references embedded
// within environment variable values.
//
// Supported syntax:
//
//	${VAR}         — replaced with the value of VAR; error if unset
//	${VAR:-default} — replaced with the value of VAR, or "default" if unset
//	$VAR           — shorthand form; error if unset
//
// Example:
//
//	interp := interpolator.New(os.LookupEnv)
//	result, err := interp.Interpolate("http://${HOST:-localhost}:${PORT}")
//
// InterpolateMap applies the same logic to all values in a string map, which
// is the primary use-case when processing a loaded env-file before deployment.
package interpolator
