package fga

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var CheckAttributes = map[string]schema.Attribute{
	"id":            stringattr.Identifier(),
	"resource":      stringattr.Required(),
	"resource_type": stringattr.Required(),
	"relation":      stringattr.Required(),
	"target":        stringattr.Required(),
	"target_type":   stringattr.Required(),
	"allowed": schema.BoolAttribute{
		Computed: true,
	},
}

type CheckModel struct {
	ID           stringattr.Type `tfsdk:"id"`
	Resource     stringattr.Type `tfsdk:"resource"`
	ResourceType stringattr.Type `tfsdk:"resource_type"`
	Relation     stringattr.Type `tfsdk:"relation"`
	Target       stringattr.Type `tfsdk:"target"`
	TargetType   stringattr.Type `tfsdk:"target_type"`
	Allowed      types.Bool      `tfsdk:"allowed"`
}
