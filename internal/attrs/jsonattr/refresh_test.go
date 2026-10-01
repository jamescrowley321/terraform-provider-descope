package jsonattr

import (
	"testing"

	"github.com/descope/terraform-provider-descope/internal/helpers"
	"github.com/stretchr/testify/require"
)

func TestRefreshKeepsVersionAndPrecision(t *testing.T) {
	before := Value(`{"metadata":{"componentsVersion":"old","name":"before"},"id":9007199254740993}`)
	data := map[string]any{}
	Get(Value(`{"metadata":{"componentsVersion":"new","name":"after"},"id":9007199254740994}`), data, helpers.RootKey)
	Refresh(&before, data, "metadata", "componentsVersion")
	require.JSONEq(t, `{"metadata":{"componentsVersion":"old","name":"after"},"id":9007199254740994}`, before.ValueString())
	metadata, ok := data["metadata"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "new", metadata["componentsVersion"])
	require.Contains(t, before.ValueString(), "9007199254740994")
}

func TestFingerprintOnlyIgnoresOwnedVersion(t *testing.T) {
	first := map[string]any{"componentsVersion": "1", "styles": map[string]any{"componentsVersion": "theme", "color": "red"}}
	second := map[string]any{"componentsVersion": "2", "styles": map[string]any{"componentsVersion": "theme", "color": "red"}}
	a, err := Fingerprint(first, "componentsVersion")
	require.NoError(t, err)
	b, err := Fingerprint(second, "componentsVersion")
	require.NoError(t, err)
	require.Equal(t, a, b)
	second["styles"] = map[string]any{"componentsVersion": "changed", "color": "red"}
	b, err = Fingerprint(second, "componentsVersion")
	require.NoError(t, err)
	require.NotEqual(t, a, b)
}
