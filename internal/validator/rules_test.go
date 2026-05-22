package validator_test

import (
	"testing"

	"github.com/envchain/envchain/internal/validator"
)

func TestRequireNonEmpty_Pass(t *testing.T) {
	rule := validator.RequireNonEmpty()
	if err := rule("KEY", "value"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRequireNonEmpty_Fail(t *testing.T) {
	rule := validator.RequireNonEmpty()
	if err := rule("KEY", ""); err == nil {
		t.Error("expected error for empty value")
	}
}

func TestRequireMatch_Pass(t *testing.T) {
	rule := validator.RequireMatch(`^\d+$`)
	if err := rule("PORT", "8080"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRequireMatch_Fail(t *testing.T) {
	rule := validator.RequireMatch(`^\d+$`)
	if err := rule("PORT", "abc"); err == nil {
		t.Error("expected error for non-numeric value")
	}
}

func TestRequireOneOf_Pass(t *testing.T) {
	rule := validator.RequireOneOf("debug", "info", "warn", "error")
	if err := rule("LOG_LEVEL", "info"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRequireOneOf_Fail(t *testing.T) {
	rule := validator.RequireOneOf("debug", "info", "warn", "error")
	if err := rule("LOG_LEVEL", "verbose"); err == nil {
		t.Error("expected error for disallowed value")
	}
}

func TestApplyRules_MultipleErrors(t *testing.T) {
	rules := []validator.Rule{
		validator.RequireNonEmpty(),
		validator.RequireMatch(`^\d+$`),
	}
	errs := validator.ApplyRules("PORT", "", rules)
	if len(errs) != 2 {
		t.Errorf("expected 2 errors, got %d", len(errs))
	}
}

func TestApplyRules_NoErrors(t *testing.T) {
	rules := []validator.Rule{
		validator.RequireNonEmpty(),
		validator.RequireMatch(`^\d+$`),
	}
	errs := validator.ApplyRules("PORT", "3000", rules)
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %d: %v", len(errs), errs)
	}
}
