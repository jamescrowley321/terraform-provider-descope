package list

import (
	"context"
	"testing"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModelToRequest(t *testing.T) {
	ctx := context.Background()
	diags := diag.Diagnostics{}

	model := &Model{
		Name:        types.StringValue("Test List"),
		Description: types.StringValue("A test list"),
		Type:        types.StringValue("ips"),
		Data:        strsetattr.Value([]string{"192.0.2.1", "198.51.100.0/24"}),
	}

	req := ModelToRequest(ctx, model, &diags)

	require.False(t, diags.HasError())
	assert.Equal(t, "Test List", req.Name)
	assert.Equal(t, "A test list", req.Description)
	assert.Equal(t, descope.ListTypeIPs, req.Type)
	assert.ElementsMatch(t, []string{"192.0.2.1", "198.51.100.0/24"}, req.Data)
}

func TestModelToRequestEmpty(t *testing.T) {
	ctx := context.Background()
	diags := diag.Diagnostics{}

	model := &Model{
		Name:        types.StringValue("Empty List"),
		Description: types.StringValue(""),
		Type:        types.StringValue("texts"),
		Data:        strsetattr.Value([]string{}),
	}

	req := ModelToRequest(ctx, model, &diags)

	require.False(t, diags.HasError())
	assert.Equal(t, "Empty List", req.Name)
	assert.Equal(t, descope.ListTypeTexts, req.Type)
	assert.Empty(t, req.Data)
}

func TestRefreshModelFromResponse(t *testing.T) {
	ctx := context.Background()
	model := &Model{}

	list := &descope.List{
		ID:          "list-123",
		Name:        "IP Allowlist",
		Description: "Production IPs",
		Type:        descope.ListTypeIPs,
		Data:        []any{"192.0.2.1", "192.0.2.2"},
	}

	RefreshModelFromResponse(ctx, model, list)

	assert.Equal(t, "list-123", model.ID.ValueString())
	assert.Equal(t, "IP Allowlist", model.Name.ValueString())
	assert.Equal(t, "Production IPs", model.Description.ValueString())
	assert.Equal(t, "ips", model.Type.ValueString())
	assert.False(t, model.Data.IsNull())
	elems := model.Data.Elements()
	assert.Len(t, elems, 2)
}

func TestRefreshModelFromResponseNilData(t *testing.T) {
	ctx := context.Background()
	model := &Model{}

	list := &descope.List{
		ID:   "list-456",
		Name: "Empty",
		Type: descope.ListTypeTexts,
		Data: nil,
	}

	RefreshModelFromResponse(ctx, model, list)

	assert.Equal(t, "list-456", model.ID.ValueString())
	assert.False(t, model.Data.IsNull())
}

func TestDataToStringSetWithStringSlice(t *testing.T) {
	ctx := context.Background()
	result := dataToStringSet(ctx, []string{"a", "b", "c"})
	assert.False(t, result.IsNull())
}
