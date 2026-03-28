package flows

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/helpers"
)

var StylesAttributes = map[string]schema.Attribute{
	"data": stringattr.Required(stringattr.JSONValidator("styles")),
}

type StylesModel struct {
	Data stringattr.Type `tfsdk:"data"`
}

func (m *StylesModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	data["data"] = getStylesData(m.Data, h)
	return data
}

func (m *StylesModel) SetValues(h *helpers.Handler, data map[string]any) {
	if m.Data.ValueString() != "" {
		return // We do not currently update the styles data if it's already set because it might be different after apply
	}

	if v, ok := data["data"].(map[string]any); ok {
		b, err := json.Marshal(v)
		if err != nil {
			h.Error("Invalid style data", "Failed to parse JSON: %s", err.Error())
			return
		}
		m.Data = stringattr.Value(string(b))
	}
}

// Computed Mapping

func getStylesData(data stringattr.Type, _ *helpers.Handler) map[string]any {
	m := map[string]any{}
	if err := json.Unmarshal([]byte(data.ValueString()), &m); err != nil {
		panic("Invalid styles data after validation: " + err.Error())
	}
	return m
}
