package fga

import (
	"testing"

	"github.com/descope/go-sdk/descope"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestSchemaModelToSDK(t *testing.T) {
	model := &SchemaModel{
		Schema: types.StringValue(`{"types":{"document":{"relations":{"owner":{"this":{}}}}}}`),
	}

	sdk := SchemaModelToSDK(model)
	assert.Equal(t, `{"types":{"document":{"relations":{"owner":{"this":{}}}}}}`, sdk.Schema)
}

func TestSchemaModelToSDKEmpty(t *testing.T) {
	model := &SchemaModel{
		Schema: types.StringValue(""),
	}

	sdk := SchemaModelToSDK(model)
	assert.Equal(t, "", sdk.Schema)
}

func TestSchemaModelFromSDK(t *testing.T) {
	model := &SchemaModel{}
	fgaSchema := &descope.FGASchema{
		Schema: `{"types":{"document":{"relations":{"viewer":{"this":{}}}}}}`,
	}

	SchemaModelFromSDK(model, fgaSchema)
	assert.Equal(t, `{"types":{"document":{"relations":{"viewer":{"this":{}}}}}}`, model.Schema.ValueString())
}

func TestSchemaModelFromSDKEmpty(t *testing.T) {
	model := &SchemaModel{}
	fgaSchema := &descope.FGASchema{Schema: ""}

	SchemaModelFromSDK(model, fgaSchema)
	assert.Equal(t, "", model.Schema.ValueString())
}

func TestSchemaModelFromSDKTrimsTrailingWhitespace(t *testing.T) {
	model := &SchemaModel{}
	fgaSchema := &descope.FGASchema{
		Schema: "model AuthZ 1.0\n\ntype user\n\n",
	}

	SchemaModelFromSDK(model, fgaSchema)
	assert.Equal(t, "model AuthZ 1.0\n\ntype user", model.Schema.ValueString())
}
