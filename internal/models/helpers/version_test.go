package helpers

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func TestEnsureModelVersion(t *testing.T) {
	t.Run("no warning for current version", func(t *testing.T) {
		var diags diag.Diagnostics
		EnsureModelVersion(1.0, &diags)
		if diags.HasError() {
			t.Fatal("unexpected error")
		}
		if len(diags.Warnings()) != 0 {
			t.Fatalf("expected no warnings, got %d", len(diags.Warnings()))
		}
	})

	t.Run("no warning for older version", func(t *testing.T) {
		var diags diag.Diagnostics
		EnsureModelVersion(0.5, &diags)
		if len(diags.Warnings()) != 0 {
			t.Fatalf("expected no warnings, got %d", len(diags.Warnings()))
		}
	})

	t.Run("warns for newer version", func(t *testing.T) {
		var diags diag.Diagnostics
		EnsureModelVersion(2.0, &diags)
		if len(diags.Warnings()) != 1 {
			t.Fatalf("expected 1 warning, got %d", len(diags.Warnings()))
		}
		w := diags.Warnings()[0]
		if w.Summary() != "Update the Descope terraform provider" {
			t.Fatalf("unexpected warning summary: %q", w.Summary())
		}
	})
}
