package prefixer_test

import (
	"testing"

	"github.com/nicholas-eden/envchain/internal/prefixer"
)

func TestAdd_PrependsMissingPrefix(t *testing.T) {
	p := prefixer.New("APP_")
	result := p.Add(map[string]string{"DB_HOST": "localhost", "PORT": "5432"})

	if _, ok := result["APP_DB_HOST"]; !ok {
		t.Error("expected APP_DB_HOST to be present")
	}
	if _, ok := result["APP_PORT"]; !ok {
		t.Error("expected APP_PORT to be present")
	}
}

func TestAdd_SkipsAlreadyPrefixedKeys(t *testing.T) {
	p := prefixer.New("APP_")
	result := p.Add(map[string]string{"APP_DB_HOST": "localhost"})

	if _, dup := result["APP_APP_DB_HOST"]; dup {
		t.Error("expected no double-prefix for already-prefixed key")
	}
	if v := result["APP_DB_HOST"]; v != "localhost" {
		t.Errorf("expected value 'localhost', got %q", v)
	}
}

func TestStrip_RemovesPrefix(t *testing.T) {
	p := prefixer.New("APP_")
	result := p.Strip(map[string]string{"APP_DB_HOST": "localhost", "APP_PORT": "5432"})

	if _, ok := result["DB_HOST"]; !ok {
		t.Error("expected DB_HOST after stripping APP_")
	}
	if _, ok := result["PORT"]; !ok {
		t.Error("expected PORT after stripping APP_")
	}
}

func TestStrip_PassesThroughUnmatchedKeys(t *testing.T) {
	p := prefixer.New("APP_")
	result := p.Strip(map[string]string{"OTHER_KEY": "value"})

	if v := result["OTHER_KEY"]; v != "value" {
		t.Errorf("expected 'value', got %q", v)
	}
}

func TestReplace_SwapsPrefix(t *testing.T) {
	env := map[string]string{
		"OLD_HOST": "localhost",
		"OLD_PORT": "5432",
		"UNRELATED": "yes",
	}
	result := prefixer.Replace(env, "OLD_", "NEW_")

	if _, ok := result["NEW_HOST"]; !ok {
		t.Error("expected NEW_HOST")
	}
	if _, ok := result["NEW_PORT"]; !ok {
		t.Error("expected NEW_PORT")
	}
	if v := result["UNRELATED"]; v != "yes" {
		t.Errorf("expected UNRELATED to pass through, got %q", v)
	}
	if _, old := result["OLD_HOST"]; old {
		t.Error("expected OLD_HOST to be replaced, not retained")
	}
}

func TestAdd_EmptyMap(t *testing.T) {
	p := prefixer.New("SVC_")
	result := p.Add(map[string]string{})
	if len(result) != 0 {
		t.Errorf("expected empty map, got %d entries", len(result))
	}
}
