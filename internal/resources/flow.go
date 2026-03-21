package resources

import (
	"context"
	"encoding/json"

	"github.com/descope/go-sdk/descope/sdk"
	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/flow"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &flowResource{}
	_ resource.ResourceWithConfigure   = &flowResource{}
	_ resource.ResourceWithImportState = &flowResource{}
)

func NewFlowResource() resource.Resource {
	return &flowResource{}
}

type flowResource struct {
	management sdk.Management
}

func (r *flowResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if data, ok := req.ProviderData.(*infra.ProviderData); ok {
		r.management = data.Management
	}
}

func (r *flowResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_flow"
}

func (r *flowResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Descope authentication flow. Flows are imported/exported as JSON definitions.",
		Attributes:  flow.Attributes,
	}
}

func (r *flowResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating flow resource")

	var model flow.Model
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	flowID := model.FlowID.ValueString()
	flowData, err := parseJSON(model.Definition.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid flow definition", "The definition must be valid JSON: "+err.Error())
		return
	}

	err = infra.RetryOnRateLimitNoResult(ctx, func() error {
		return r.management.Flow().ImportFlow(ctx, flowID, flowData)
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating flow", err.Error())
		return
	}

	// Read back the server-normalized definition
	r.refreshModel(ctx, &model, flowID, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
	tflog.Info(ctx, "Flow resource created")
}

func (r *flowResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading flow resource")

	var model flow.Model
	if resp.Diagnostics.Append(req.State.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	flowID := model.FlowID.ValueString()
	if flowID == "" {
		flowID = model.ID.ValueString()
	}

	exported, err := infra.RetryOnRateLimit(ctx, func() (map[string]any, error) {
		return r.management.Flow().ExportFlow(ctx, flowID)
	})
	if err != nil {
		if infra.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading flow", err.Error())
		return
	}

	defJSON, err := json.Marshal(exported)
	if err != nil {
		resp.Diagnostics.AddError("Error serializing flow", err.Error())
		return
	}

	model.ID = types.StringValue(flowID)
	model.FlowID = types.StringValue(flowID)
	model.Definition = types.StringValue(string(defJSON))
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
	tflog.Info(ctx, "Flow resource read")
}

func (r *flowResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating flow resource")

	var plan flow.Model
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	flowID := plan.FlowID.ValueString()
	flowData, err := parseJSON(plan.Definition.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid flow definition", "The definition must be valid JSON: "+err.Error())
		return
	}

	err = infra.RetryOnRateLimitNoResult(ctx, func() error {
		return r.management.Flow().ImportFlow(ctx, flowID, flowData)
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating flow", err.Error())
		return
	}

	r.refreshModel(ctx, &plan, flowID, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, "Flow resource updated")
}

func (r *flowResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting flow resource")

	var model flow.Model
	if resp.Diagnostics.Append(req.State.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	flowID := model.FlowID.ValueString()
	err := infra.RetryOnRateLimitNoResult(ctx, func() error {
		return r.management.Flow().DeleteFlows(ctx, []string{flowID})
	})
	if err != nil {
		if infra.IsNotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting flow", err.Error())
		return
	}

	tflog.Info(ctx, "Flow resource deleted")
}

func (r *flowResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing flow resource")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *flowResource) refreshModel(ctx context.Context, model *flow.Model, flowID string, diags *diag.Diagnostics) {
	exported, err := infra.RetryOnRateLimit(ctx, func() (map[string]any, error) {
		return r.management.Flow().ExportFlow(ctx, flowID)
	})
	if err != nil {
		diags.AddError("Error reading flow after write", err.Error())
		return
	}

	defJSON, err := json.Marshal(exported)
	if err != nil {
		diags.AddError("Error serializing flow", err.Error())
		return
	}

	model.ID = types.StringValue(flowID)
	model.FlowID = types.StringValue(flowID)
	model.Definition = types.StringValue(string(defJSON))
}

func parseJSON(s string) (map[string]any, error) {
	var data map[string]any
	if err := json.Unmarshal([]byte(s), &data); err != nil {
		return nil, err
	}
	return data, nil
}
