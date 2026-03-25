package list

import (
	"context"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/convert"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ModelToRequest(ctx context.Context, model *Model, diags *diag.Diagnostics) *descope.ListRequest {
	return &descope.ListRequest{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		Type:        descope.ListType(model.Type.ValueString()),
		Data:        convert.StringSetToSlice(ctx, model.Data, diags),
	}
}

func RefreshModelFromResponse(ctx context.Context, model *Model, list *descope.List) {
	model.ID = types.StringValue(list.ID)
	model.Name = types.StringValue(list.Name)
	model.Description = types.StringValue(list.Description)
	model.Type = types.StringValue(string(list.Type))
	model.Data = dataToStringSet(ctx, list.Data)
}

// dataToStringSet converts the SDK's Data field (any) to a Terraform string set.
// For ips/texts lists, Data is []any containing string elements.
func dataToStringSet(ctx context.Context, data any) strsetattr.Type {
	if data == nil {
		return strsetattr.ValueCtx(ctx, []string{})
	}
	switch v := data.(type) {
	case []any:
		strs := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok {
				strs = append(strs, s)
			}
		}
		return strsetattr.ValueCtx(ctx, strs)
	case []string:
		return strsetattr.ValueCtx(ctx, v)
	default:
		return strsetattr.ValueCtx(ctx, []string{})
	}
}
