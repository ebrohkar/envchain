// Package pipeline provides a high-level orchestration layer for envchain.
//
// It wires together the loader, resolver, validator, and reporter packages
// into a single Run function that processes an environment variable chain
// end-to-end:
//
//  1. Resolve all variable references using the resolver.
//  2. Apply validation rules via the validator.
//  3. Emit a structured report through the reporter.
//
// Example usage:
//
//	result, err := pipeline.Run(c, rules, pipeline.Options{
//		EnvSource: os.LookupEnv,
//		Format:    reporter.FormatText,
//		Writer:    os.Stdout,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	if !result.ChainOK || !result.ValidatorOK {
//		os.Exit(1)
//	}
package pipeline
