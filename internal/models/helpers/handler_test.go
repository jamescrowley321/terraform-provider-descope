package helpers

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func TestNewHandler(t *testing.T) {
	var diags diag.Diagnostics
	h := NewHandler(context.Background(), &diags)

	if h.Ctx == nil {
		t.Fatal("expected non-nil context")
	}
	if h.Diagnostics == nil {
		t.Fatal("expected non-nil diagnostics")
	}
	if h.Refs == nil {
		t.Fatal("expected non-nil references map")
	}
}

func TestHandler_Error(t *testing.T) {
	var diags diag.Diagnostics
	h := NewHandler(context.Background(), &diags)

	h.Error("test summary", "detail: %s", "info")
	if !diags.HasError() {
		t.Fatal("expected error in diagnostics")
	}
	errs := diags.Errors()
	if errs[0].Summary() != "test summary" {
		t.Fatalf("expected 'test summary', got %q", errs[0].Summary())
	}
	if errs[0].Detail() != "detail: info" {
		t.Fatalf("expected 'detail: info', got %q", errs[0].Detail())
	}
}

func TestHandler_Warn(t *testing.T) {
	var diags diag.Diagnostics
	h := NewHandler(context.Background(), &diags)

	h.Warn("warn summary", "detail: %d", 42)
	warnings := diags.Warnings()
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(warnings))
	}
	if warnings[0].Detail() != "detail: 42" {
		t.Fatalf("expected 'detail: 42', got %q", warnings[0].Detail())
	}
}

func TestHandler_Invalid(t *testing.T) {
	var diags diag.Diagnostics
	h := NewHandler(context.Background(), &diags)

	h.Invalid("field %s is invalid", "name")
	if !diags.HasError() {
		t.Fatal("expected error")
	}
	if diags.Errors()[0].Summary() != "Invalid Attribute Value" {
		t.Fatalf("unexpected summary: %q", diags.Errors()[0].Summary())
	}

	// Second Invalid is suppressed because diagnostics already has an error
	h.Invalid("another error")
	if len(diags.Errors()) != 1 {
		t.Fatalf("expected 1 error (second suppressed by HasError guard), got %d", len(diags.Errors()))
	}
}

func TestHandler_Missing(t *testing.T) {
	var diags diag.Diagnostics
	h := NewHandler(context.Background(), &diags)

	h.Missing("field %s is required", "name")
	if !diags.HasError() {
		t.Fatal("expected error")
	}
	if diags.Errors()[0].Summary() != "Missing Attribute Value" {
		t.Fatalf("unexpected summary: %q", diags.Errors()[0].Summary())
	}

	// Second Missing is suppressed because diagnostics already has an error
	h.Missing("another missing")
	if len(diags.Errors()) != 1 {
		t.Fatalf("expected 1 error (second suppressed by HasError guard), got %d", len(diags.Errors()))
	}
}

func TestHandler_Conflict(t *testing.T) {
	var diags diag.Diagnostics
	h := NewHandler(context.Background(), &diags)

	h.Conflict("field %s conflicts with %s", "a", "b")
	if !diags.HasError() {
		t.Fatal("expected error")
	}
	if diags.Errors()[0].Summary() != "Conflicting Attribute Values" {
		t.Fatalf("unexpected summary: %q", diags.Errors()[0].Summary())
	}

	// Second call should be suppressed
	h.Conflict("another conflict")
	if len(diags.Errors()) != 1 {
		t.Fatalf("expected 1 error (suppressed duplicate), got %d", len(diags.Errors()))
	}
}

func TestHandler_CrossTypeSuppression(t *testing.T) {
	t.Run("Error suppresses subsequent Invalid", func(t *testing.T) {
		var diags diag.Diagnostics
		h := NewHandler(context.Background(), &diags)

		h.Error("first error", "detail")
		h.Invalid("should be suppressed")
		if len(diags.Errors()) != 1 {
			t.Fatalf("expected 1 error, got %d", len(diags.Errors()))
		}
		if diags.Errors()[0].Summary() != "first error" {
			t.Fatalf("expected first error to remain, got %q", diags.Errors()[0].Summary())
		}
	})

	t.Run("Invalid suppresses subsequent Missing", func(t *testing.T) {
		var diags diag.Diagnostics
		h := NewHandler(context.Background(), &diags)

		h.Invalid("invalid field")
		h.Missing("should be suppressed")
		if len(diags.Errors()) != 1 {
			t.Fatalf("expected 1 error, got %d", len(diags.Errors()))
		}
		if diags.Errors()[0].Summary() != "Invalid Attribute Value" {
			t.Fatalf("expected Invalid error to remain, got %q", diags.Errors()[0].Summary())
		}
	})

	t.Run("Missing suppresses subsequent Conflict", func(t *testing.T) {
		var diags diag.Diagnostics
		h := NewHandler(context.Background(), &diags)

		h.Missing("missing field")
		h.Conflict("should be suppressed")
		if len(diags.Errors()) != 1 {
			t.Fatalf("expected 1 error, got %d", len(diags.Errors()))
		}
		if diags.Errors()[0].Summary() != "Missing Attribute Value" {
			t.Fatalf("expected Missing error to remain, got %q", diags.Errors()[0].Summary())
		}
	})

	t.Run("Error does not suppress warnings", func(t *testing.T) {
		var diags diag.Diagnostics
		h := NewHandler(context.Background(), &diags)

		h.Error("an error", "detail")
		h.Warn("a warning", "detail")
		if len(diags.Errors()) != 1 {
			t.Fatalf("expected 1 error, got %d", len(diags.Errors()))
		}
		if len(diags.Warnings()) != 1 {
			t.Fatalf("expected 1 warning, got %d", len(diags.Warnings()))
		}
	})
}
