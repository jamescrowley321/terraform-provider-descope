package fga

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/attrs/stringattr"
)

type CheckModel struct {
	ID           stringattr.Type `tfsdk:"id"`
	Resource     stringattr.Type `tfsdk:"resource"`
	ResourceType stringattr.Type `tfsdk:"resource_type"`
	Relation     stringattr.Type `tfsdk:"relation"`
	Target       stringattr.Type `tfsdk:"target"`
	TargetType   stringattr.Type `tfsdk:"target_type"`
	Allowed      types.Bool      `tfsdk:"allowed"`
}
