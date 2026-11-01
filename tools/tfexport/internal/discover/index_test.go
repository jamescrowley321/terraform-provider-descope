package discover

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSupportedTemplateMethodsAreDiscovered(t *testing.T) {
	files := map[string]any{
		"auth/invite.json":        map[string]any{"emailTemplates": []any{map[string]any{"id": "email", "name": "invitation"}}},
		"auth/enchantedlink.json": map[string]any{"textTemplates": []any{map[string]any{"id": "text", "name": "enchanted"}}},
	}
	instances, warnings := Instances(files, nil)
	require.Empty(t, warnings)
	require.ElementsMatch(t, []Instance{{Resource: "descope_email_template", ID: "email", Scope: "invite", Name: "invitation"}, {Resource: "descope_text_template", ID: "text", Scope: "enchantedlink", Name: "enchanted"}}, instances)
}
