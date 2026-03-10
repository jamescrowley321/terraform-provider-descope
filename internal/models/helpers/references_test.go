package helpers

import (
	"testing"
)

func TestModelReference_ReferenceValue(t *testing.T) {
	t.Run("returns ID when set", func(t *testing.T) {
		ref := &ModelReference{ID: "id-123", Key: "key-456"}
		if ref.ReferenceValue() != "id-123" {
			t.Fatalf("expected 'id-123', got %q", ref.ReferenceValue())
		}
	})

	t.Run("returns Key when ID is empty", func(t *testing.T) {
		ref := &ModelReference{Key: "key-456"}
		if ref.ReferenceValue() != "key-456" {
			t.Fatalf("expected 'key-456', got %q", ref.ReferenceValue())
		}
	})

	t.Run("returns empty when both empty", func(t *testing.T) {
		ref := &ModelReference{}
		if ref.ReferenceValue() != "" {
			t.Fatalf("expected empty string, got %q", ref.ReferenceValue())
		}
	})
}

func TestModelReference_ProviderValue(t *testing.T) {
	t.Run("returns type:id when both set", func(t *testing.T) {
		ref := &ModelReference{Type: "connector", ID: "id-123"}
		if ref.ProviderValue() != "connector:id-123" {
			t.Fatalf("expected 'connector:id-123', got %q", ref.ProviderValue())
		}
	})

	t.Run("falls back to ReferenceValue when type is empty", func(t *testing.T) {
		ref := &ModelReference{ID: "id-123"}
		if ref.ProviderValue() != "id-123" {
			t.Fatalf("expected 'id-123', got %q", ref.ProviderValue())
		}
	})

	t.Run("falls back to ReferenceValue when id is empty", func(t *testing.T) {
		ref := &ModelReference{Type: "connector", Key: "key-456"}
		if ref.ProviderValue() != "key-456" {
			t.Fatalf("expected 'key-456', got %q", ref.ProviderValue())
		}
	})
}

func TestReferencesMap_Add_Get(t *testing.T) {
	t.Run("stores reference and retrieves by key and name", func(t *testing.T) {
		refs := ReferencesMap{}
		refs.Add(ConnectorReferenceKey, "http", "id-123", "MyConnector")

		ref := refs.Get(ConnectorReferenceKey, "MyConnector")
		if ref == nil {
			t.Fatal("expected reference, got nil")
		}
		if ref.ID != "id-123" {
			t.Fatalf("expected ID 'id-123', got %q", ref.ID)
		}
		if ref.Type != "http" {
			t.Fatalf("expected Type 'http', got %q", ref.Type)
		}
	})

	t.Run("returns nil for missing reference", func(t *testing.T) {
		refs := ReferencesMap{}
		ref := refs.Get(ConnectorReferenceKey, "Missing")
		if ref != nil {
			t.Fatalf("expected nil, got %v", ref)
		}
	})

	t.Run("returns synthetic reference for DescopeConnector without Add", func(t *testing.T) {
		refs := ReferencesMap{}
		ref := refs.Get(ConnectorReferenceKey, DescopeConnector)
		if ref == nil {
			t.Fatal("expected Descope connector reference, got nil")
		}
		if ref.Key != DescopeConnector {
			t.Fatalf("expected Key %q, got %q", DescopeConnector, ref.Key)
		}
	})

	t.Run("auto-generates key when added with empty ID", func(t *testing.T) {
		refs := ReferencesMap{}
		refs.Add(RoleReferenceKey, "role", "", "MyRole")

		ref := refs.Get(RoleReferenceKey, "MyRole")
		if ref == nil {
			t.Fatal("expected reference, got nil")
		}
		if ref.Key == "" {
			t.Fatal("expected generated key, got empty string")
		}
		if ref.ID != "" {
			t.Fatalf("expected empty ID, got %q", ref.ID)
		}
	})
}

func TestReferencesMap_Name(t *testing.T) {
	t.Run("reverse-looks up name from reference ID", func(t *testing.T) {
		refs := ReferencesMap{}
		refs.Add(ConnectorReferenceKey, "http", "id-123", "MyConnector")

		name := refs.Name("id-123")
		if name != "MyConnector" {
			t.Fatalf("expected 'MyConnector', got %q", name)
		}
	})

	t.Run("returns empty for unknown ID", func(t *testing.T) {
		refs := ReferencesMap{}
		name := refs.Name("unknown-id")
		if name != "" {
			t.Fatalf("expected empty string, got %q", name)
		}
	})

	t.Run("returns DescopeConnector for DescopeConnector ID", func(t *testing.T) {
		refs := ReferencesMap{}
		name := refs.Name(DescopeConnector)
		if name != DescopeConnector {
			t.Fatalf("expected %q, got %q", DescopeConnector, name)
		}
	})
}
