package registry_test

import (
	"context"
	"slices"
	"testing"

	"github.com/descope/terraform-provider-descope/internal/provider"
	"github.com/descope/terraform-provider-descope/internal/resources"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/registry"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestAllResourcesExportable(t *testing.T) {
	ctx := context.Background()

	total := 0
	notExportable := []string{}
	for _, constructor := range provider.NewDescopeProvider("test")().Resources(ctx) {
		total++
		if _, ok := constructor().(resources.ExportableResource); !ok {
			metadata := resource.MetadataResponse{}
			constructor().Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "descope"}, &metadata)
			notExportable = append(notExportable, metadata.TypeName)
		}
	}
	slices.Sort(notExportable)
	if !slices.Equal(notExportable, []string{"descope_sso", "descope_tenant"}) {
		t.Errorf("unexpected resources outside project-snapshot export: %v", notExportable)
	}

	loaded := registry.Load(ctx)
	if len(loaded) != total-len(notExportable) {
		t.Errorf("registry has %d resources, expected %d", len(loaded), total-len(notExportable))
	}

	for _, name := range []string{
		"descope_project", "descope_flow", "descope_widget", "descope_styles", "descope_magiclink_settings",
		"descope_email_template", "descope_text_template", "descope_smtp_connector",
		"descope_role", "descope_oidc_app", "descope_app_role", "descope_jwt_template",
	} {
		exportable, ok := loaded[name]
		if !ok {
			t.Errorf("expected %s in the registry", name)
			continue
		}
		if len(exportable.ExportSchema().Attributes) == 0 {
			t.Errorf("expected %s to have a non-empty schema", name)
		}
	}

	for name, exportable := range loaded {
		if name != "descope_"+exportable.ExportName() {
			t.Errorf("registry key %q doesn't match resource name %q", name, exportable.ExportName())
		}
	}
}
