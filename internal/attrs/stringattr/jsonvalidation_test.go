package stringattr

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestJSONValidation(t *testing.T) {
	for _, value := range []string{"", "  ", "null", "[]", "true", `"text"`, "{} {}", "{} trailing"} {
		t.Run(value, func(t *testing.T) {
			var response validator.StringResponse
			JSONValidator().ValidateString(context.Background(), validator.StringRequest{Path: path.Root("custom_claims"), ConfigValue: types.StringValue(value)}, &response)
			require.True(t, response.Diagnostics.HasError())
		})
	}
	for _, value := range []types.String{types.StringUnknown(), types.StringNull(), types.StringValue(`{"nested":{"list":[true,9007199254740993]}}`), types.StringValue("{}")} {
		var response validator.StringResponse
		JSONValidator().ValidateString(context.Background(), validator.StringRequest{ConfigValue: value}, &response)
		require.False(t, response.Diagnostics.HasError())
	}
}
