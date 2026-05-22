package chain_test

import (
	"testing"

	"github.com/envchain/envchain/internal/chain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockEnv(vars map[string]string) func(string) (string, bool) {
	return func(key string) (string, bool) {
		v, ok := vars[key]
		return v, ok
	}
}

func TestChain_AllResolved(t *testing.T) {
	c := chain.New()
	c.Add("DB_HOST")
	c.Add("DB_PORT", "DB_HOST")

	c.Resolve(mockEnv(map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}))

	errs := c.Validate()
	assert.Empty(t, errs)
}

func TestChain_MissingVariable(t *testing.T) {
	c := chain.New()
	c.Add("API_KEY")

	c.Resolve(mockEnv(map[string]string{}))

	errs := c.Validate()
	require.Len(t, errs, 1)
	assert.Contains(t, errs[0].Error(), "missing")
}

func TestChain_EmptyVariable(t *testing.T) {
	c := chain.New()
	c.Add("SECRET")

	c.Resolve(mockEnv(map[string]string{"SECRET": ""}))

	errs := c.Validate()
	require.Len(t, errs, 1)
	assert.Contains(t, errs[0].Error(), "empty")
}

func TestChain_UnresolvedDependency(t *testing.T) {
	c := chain.New()
	c.Add("BASE_URL")
	c.Add("SERVICE_URL", "BASE_URL")

	c.Resolve(mockEnv(map[string]string{
		"SERVICE_URL": "http://svc",
		// BASE_URL intentionally missing
	}))

	errs := c.Validate()
	// BASE_URL missing + SERVICE_URL depends on unresolved BASE_URL
	assert.GreaterOrEqual(t, len(errs), 2)
}

func TestChain_UndeclaredDependency(t *testing.T) {
	c := chain.New()
	c.Add("FOO", "BAR") // BAR never declared

	c.Resolve(mockEnv(map[string]string{"FOO": "value"}))

	errs := c.Validate()
	require.Len(t, errs, 1)
	assert.Contains(t, errs[0].Error(), "undeclared")
}
