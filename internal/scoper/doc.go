// Package scoper provides namespace-scoped views over flat environment variable
// maps. It is useful when multiple services share a single environment and each
// service's variables are distinguished by a common prefix such as "APP_",
// "DB_", or "CACHE_".
//
// Example usage:
//
//	env := map[string]string{
//		"DB_HOST": "localhost",
//		"DB_PORT": "5432",
//		"APP_PORT": "8080",
//	}
//	s := scoper.New(env, "_")
//	dbVars := s.Scope("DB") // {"HOST": "localhost", "PORT": "5432"}
package scoper
