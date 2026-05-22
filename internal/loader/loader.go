package loader

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// EnvMap is a map of environment variable names to their values.
type EnvMap map[string]string

// FromFile reads a .env-style file and returns an EnvMap.
// Lines starting with '#' are treated as comments and ignored.
// Empty lines are also ignored.
// Each valid line must be in KEY=VALUE format.
func FromFile(path string) (EnvMap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("loader: open file %q: %w", path, err)
	}
	defer f.Close()

	env := make(EnvMap)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("loader: %q line %d: invalid format, expected KEY=VALUE", path, lineNum)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if key == "" {
			return nil, fmt.Errorf("loader: %q line %d: empty key", path, lineNum)
		}

		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("loader: scanning %q: %w", path, err)
	}

	return env, nil
}

// FromEnv reads variables from the current process environment for the given keys.
func FromEnv(keys []string) EnvMap {
	env := make(EnvMap, len(keys))
	for _, k := range keys {
		env[k] = os.Getenv(k)
	}
	return env
}
