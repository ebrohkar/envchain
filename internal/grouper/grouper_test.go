package grouper_test

import (
	"testing"

	"github.com/envchain/envchain/internal/grouper"
)

func baseEnv() map[string]string {
	return map[string]string{
		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
		"APP_NAME":    "envchain",
		"APP_VERSION": "1.0.0",
		"STANDALONE":  "yes",
	}
}

func TestGroup_ByPrefix_SplitsCorrectly(t *testing.T) {
	g := grouper.New(grouper.ByPrefix("_"))
	groups := g.Group(baseEnv())

	if len(groups["DB"]) != 2 {
		t.Errorf("expected 2 DB keys, got %d", len(groups["DB"]))
	}
	if len(groups["APP"]) != 2 {
		t.Errorf("expected 2 APP keys, got %d", len(groups["APP"]))
	}
	if _, ok := groups[""]["STANDALONE"]; !ok {
		t.Error("expected STANDALONE in ungrouped bucket")
	}
}

func TestGroup_CustomFn_GroupsCorrectly(t *testing.T) {
	fn := func(key string) string {
		if key == "DB_HOST" || key == "DB_PORT" {
			return "database"
		}
		return "other"
	}
	g := grouper.New(fn)
	groups := g.Group(baseEnv())

	if len(groups["database"]) != 2 {
		t.Errorf("expected 2 keys in database group, got %d", len(groups["database"]))
	}
	if len(groups["other"]) != 3 {
		t.Errorf("expected 3 keys in other group, got %d", len(groups["other"]))
	}
}

func TestGroupNames_ReturnsSorted(t *testing.T) {
	g := grouper.New(grouper.ByPrefix("_"))
	groups := g.Group(baseEnv())
	names := grouper.GroupNames(groups)

	if len(names) == 0 {
		t.Fatal("expected at least one group name")
	}
	for i := 1; i < len(names); i++ {
		if names[i] < names[i-1] {
			t.Errorf("group names not sorted: %v", names)
		}
	}
}

func TestGroup_EmptyEnv_ReturnsEmptyMap(t *testing.T) {
	g := grouper.New(grouper.ByPrefix("_"))
	groups := g.Group(map[string]string{})

	if len(groups) != 0 {
		t.Errorf("expected empty groups, got %d", len(groups))
	}
}

func TestByPrefix_NoSep_PlacesInUngrouped(t *testing.T) {
	fn := grouper.ByPrefix("_")
	if fn("NOSEP") != "" {
		t.Error("expected empty group for key without separator")
	}
	if fn("WITH_SEP") != "WITH" {
		t.Error("expected WITH as group name")
	}
}
