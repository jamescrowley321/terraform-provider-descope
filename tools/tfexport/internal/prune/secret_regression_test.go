package prune_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/emit"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/prune"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/registry"
	"github.com/stretchr/testify/require"
)

func TestEightByEightExportKeepsKeysSecret(t *testing.T) {
	ctx := context.Background()
	loaded := registry.Load(ctx)
	for _, kind := range []string{"descope_eight_by_eight_viber_connector", "descope_eight_by_eight_whatsapp_connector"} {
		t.Run(kind, func(t *testing.T) {
			sc := loaded[kind].ExportSchema()
			obj := object(ctx, t, sc, map[string]any{"id": "C1", "project_id": "P1", "name": "test", "sub_account_id": "account", "template_id": "template", "api_key": "supplied-secret-must-not-leak"})
			attrs, _ := prune.Object(ctx, sc.Attributes, obj, true)
			dir := t.TempDir()
			require.NoError(t, emit.Write(ctx, dir, &emit.Plan{ProjectID: "P1", Resources: []emit.Resource{{Type: kind, Label: "test", EntityID: "C1", ImportID: "P1/C1", HasProjectID: true, Attrs: attrs}}}))
			entries, err := os.ReadDir(dir)
			require.NoError(t, err)
			content := ""
			for _, entry := range entries {
				data, err := os.ReadFile(dir + "/" + entry.Name())
				require.NoError(t, err)
				content += string(data)
			}
			require.NotContains(t, content, "supplied-secret-must-not-leak")
			require.Regexp(t, `sensitive\s*=\s*true`, content)
			require.True(t, strings.Contains(content, "api_key") && strings.Contains(content, "var."))
		})
	}
}
