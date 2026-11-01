package emit

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/prune"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestExtractedTemplateFilesAreDistinct(t *testing.T) {
	out := t.TempDir()
	seen := map[string]bool{}
	for _, item := range []struct{ kind, label, name, body string }{
		{"descope_email_template", "welcome", "plain_text_body", "email\ncontent"},
		{"descope_text_template", "welcome_plain", "body", "text\ncontent"},
		{"descope_text_template", "welcome", "body", "other\ncontent"},
	} {
		r := &Resource{Type: item.kind, Label: item.label, Attrs: []prune.Attr{{Name: item.name, Value: types.StringValue(item.body)}}}
		files, err := extractFiles(context.Background(), out, "", r)
		require.NoError(t, err)
		path := files[item.name]
		require.NotEmpty(t, path)
		require.False(t, seen[path])
		seen[path] = true
		body, err := os.ReadFile(filepath.Join(out, path))
		require.NoError(t, err)
		require.Equal(t, item.body, string(body))
	}
}
