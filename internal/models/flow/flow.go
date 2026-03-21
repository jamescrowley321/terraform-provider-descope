package flow

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var Attributes = map[string]schema.Attribute{
	"id": stringattr.Identifier(),
	"flow_id": schema.StringAttribute{
		Required:      true,
		PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
	},
	"definition": schema.StringAttribute{
		Required:    true,
		Description: "The flow definition as a JSON string. Use file() or jsonencode() to provide the value.",
	},
}

type Model struct {
	ID         stringattr.Type `tfsdk:"id"`
	FlowID     stringattr.Type `tfsdk:"flow_id"`
	Definition stringattr.Type `tfsdk:"definition"`
}
