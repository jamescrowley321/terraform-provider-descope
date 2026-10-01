package fga

import (
	"github.com/descope/terraform-provider-descope/internal/attrs/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CheckModel struct {
	ProjectID    types.String    `tfsdk:"project_id"`
	ID           stringattr.Type `tfsdk:"id"`
	Resource     stringattr.Type `tfsdk:"resource"`
	ResourceType stringattr.Type `tfsdk:"resource_type"`
	Relation     stringattr.Type `tfsdk:"relation"`
	Target       stringattr.Type `tfsdk:"target"`
	TargetType   stringattr.Type `tfsdk:"target_type"`
	Allowed      types.Bool      `tfsdk:"allowed"`
}
