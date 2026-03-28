package widgets

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/attrs/mapattr"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/helpers"
)

var WidgetIDValidator = mapvalidator.KeysAre(stringattr.MachineIDValidator)

func EnsureWidgetIDs(m mapattr.Type[WidgetModel], data map[string]any, key string, h *helpers.Handler) {
	values := data
	if key != helpers.RootKey {
		values, _ = data[key].(map[string]any)
	}

	for widgetID := range mapattr.Iterator(m, h) {
		if v, ok := values[widgetID].(map[string]any); ok {
			if valueID, _ := v["widgetId"].(string); valueID != "" && valueID != widgetID {
				h.Warn("Possible widget mismatch", "The '%s' widget data specifies a different widgetId '%s'. You can update the widget data to use the same widgetId or ignore this warning to use the '%s' widgetId.", widgetID, valueID, widgetID)
			}
			v["widgetId"] = widgetID
		}
	}
}
