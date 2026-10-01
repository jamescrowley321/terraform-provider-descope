package jsonattr

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"maps"

	"github.com/descope/terraform-provider-descope/internal/helpers"
)

func Fingerprint(data map[string]any, versionPath ...string) ([]byte, error) {
	normalized := withoutPath(data, versionPath)
	encoded, err := json.Marshal(normalized)
	if err != nil {
		return nil, err
	}
	return fmt.Appendf(nil, `"%x"`, sha256.Sum256(encoded)), nil
}

func Refresh(s *Type, data map[string]any, versionPath ...string) {
	if len(versionPath) == 0 {
		Set(s, data, helpers.RootKey)
		return
	}
	normalized := withoutPath(data, versionPath)
	old := map[string]any{}
	Get(*s, old, helpers.RootKey)
	previous, current := old, normalized
	for _, key := range versionPath[:len(versionPath)-1] {
		previous, _ = previous[key].(map[string]any)
		current, _ = current[key].(map[string]any)
	}
	key := versionPath[len(versionPath)-1]
	if value, ok := previous[key]; ok && current != nil {
		current[key] = value
	}
	Set(s, normalized, helpers.RootKey)
}

func withoutPath(data map[string]any, path []string) map[string]any {
	result := maps.Clone(data)
	if len(path) == 0 {
		return result
	}
	current := result
	for _, key := range path[:len(path)-1] {
		nested, ok := current[key].(map[string]any)
		if !ok {
			return result
		}
		nested = maps.Clone(nested)
		current[key] = nested
		current = nested
	}
	delete(current, path[len(path)-1])
	return result
}
