// Package resolver provides variable interpolation for envchain configs.
//
// It expands ${VAR} references within environment variable values using a
// caller-supplied EnvGetter, enabling cross-variable dependencies to be
// resolved before deployment validation.
//
// Basic usage:
//
//	r := resolver.New(os.LookupEnv)
//	resolved, missing := r.Resolve("${DB_HOST}:${DB_PORT}/mydb")
//	if len(missing) > 0 {
//		log.Fatalf("unresolved references: %v", missing)
//	}
//
// To resolve an entire map of variables at once, use ResolveAll:
//
//	results := r.ResolveAll(vars)
//	if resolver.HasMissing(results) {
//		fmt.Println(resolver.MissingSummary(results))
//	}
package resolver
