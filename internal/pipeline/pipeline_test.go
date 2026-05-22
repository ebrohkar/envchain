package pipeline_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envchain/internal/chain"
	"github.com/user/envchain/internal/pipeline"
	"github.com/user/envchain/internal/reporter"
	"github.com/user/envchain/internal/validator"
)

func mockEnv(m map[string]string) func(string) (string, bool) {
	return func(key string) (string, bool) {
		v, ok := m[key]
		return v, ok
	}
}

func TestRun_AllResolved(t *testing.T) {
	env := mockEnv(map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	})
	c := chain.New(env)
	c.Add("DB_HOST", "DB_PORT")

	var buf bytes.Buffer
	res, err := pipeline.Run(c, nil, pipeline.Options{
		EnvSource: env,
		Format:    reporter.FormatText,
		Writer:    &buf,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.ChainOK {
		t.Error("expected ChainOK to be true")
	}
	if len(res.Errors) != 0 {
		t.Errorf("expected no errors, got: %v", res.Errors)
	}
}

func TestRun_MissingVariable(t *testing.T) {
	env := mockEnv(map[string]string{
		"DB_HOST": "localhost",
	})
	c := chain.New(env)
	c.Add("DB_HOST", "DB_PORT")

	var buf bytes.Buffer
	res, err := pipeline.Run(c, nil, pipeline.Options{
		EnvSource: env,
		Format:    reporter.FormatText,
		Writer:    &buf,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ChainOK {
		t.Error("expected ChainOK to be false")
	}
	if len(res.Errors) == 0 {
		t.Error("expected errors for missing variable")
	}
}

func TestRun_ValidationFailure(t *testing.T) {
	env := mockEnv(map[string]string{
		"API_KEY": "",
	})
	c := chain.New(env)
	c.Add("API_KEY")

	rules := map[string][]validator.Rule{
		"API_KEY": {validator.RequireNonEmpty},
	}

	var buf bytes.Buffer
	res, err := pipeline.Run(c, rules, pipeline.Options{
		EnvSource: env,
		Format:    reporter.FormatText,
		Writer:    &buf,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ValidatorOK {
		t.Error("expected ValidatorOK to be false")
	}
	if !strings.Contains(buf.String(), "API_KEY") {
		t.Error("expected output to mention API_KEY")
	}
}

func TestRun_NilChain(t *testing.T) {
	_, err := pipeline.Run(nil, nil, pipeline.Options{})
	if err == nil {
		t.Error("expected error for nil chain")
	}
}
