package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadJSON(t *testing.T) {
	t.Run("reads valid JSON file", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "test.json")
		if err := os.WriteFile(path, []byte(`{"name":"test","value":42}`), 0600); err != nil {
			t.Fatal(err)
		}

		var target struct {
			Name  string `json:"name"`
			Value int    `json:"value"`
		}
		if err := ReadJSON(path, &target); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if target.Name != "test" || target.Value != 42 {
			t.Fatalf("unexpected values: %+v", target)
		}
	})

	t.Run("returns error for missing file", func(t *testing.T) {
		var target struct{}
		err := ReadJSON("/nonexistent/path.json", &target)
		if err == nil {
			t.Fatal("expected error for missing file")
		}
	})

	t.Run("returns error for invalid JSON", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "bad.json")
		if err := os.WriteFile(path, []byte(`{invalid}`), 0600); err != nil {
			t.Fatal(err)
		}

		var target struct{}
		err := ReadJSON(path, &target)
		if err == nil {
			t.Fatal("expected error for invalid JSON")
		}
	})
}

func TestWriteJSON(t *testing.T) {
	t.Run("writes JSON with indentation and trailing newline", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "out.json")

		data := struct {
			Name string `json:"name"`
		}{Name: "test"}

		if err := WriteJSON(path, &data); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		b, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}

		expected := "{\n  \"name\": \"test\"\n}\n"
		if string(b) != expected {
			t.Fatalf("expected %q, got %q", expected, string(b))
		}
	})

	t.Run("roundtrips data correctly", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "roundtrip.json")

		type Data struct {
			Items []string `json:"items"`
			Count int      `json:"count"`
		}

		original := Data{Items: []string{"a", "b"}, Count: 2}
		if err := WriteJSON(path, &original); err != nil {
			t.Fatal(err)
		}

		var loaded Data
		if err := ReadJSON(path, &loaded); err != nil {
			t.Fatal(err)
		}

		if loaded.Count != original.Count || len(loaded.Items) != len(original.Items) {
			t.Fatalf("roundtrip mismatch: %+v vs %+v", original, loaded)
		}
	})
}
