package resources

import (
	"context"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/go-sdk/descope/sdk"
	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/outboundapp"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &outboundAppResource{}
	_ resource.ResourceWithConfigure   = &outboundAppResource{}
	_ resource.ResourceWithImportState = &outboundAppResource{}
)

func NewOutboundAppResource() resource.Resource {
	return &outboundAppResource{}
}

type outboundAppResource struct {
	management sdk.Management
}

func (r *outboundAppResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if data, ok := req.ProviderData.(*infra.ProviderData); ok {
		r.management = data.Management
	}
}

func (r *outboundAppResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_outbound_application"
}

func (r *outboundAppResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Descope outbound application for OAuth integrations with external services.",
		Attributes:  outboundapp.Attributes,
	}
}

func (r *outboundAppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating outbound application resource")

	var model outboundapp.Model
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	createReq := outboundapp.ModelToCreateRequest(ctx, &model, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	app, err := infra.RetryOnRateLimit(ctx, func() (*descope.OutboundApp, error) {
		return r.management.OutboundApplication().CreateApplication(ctx, createReq)
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating outbound application", err.Error())
		return
	}

	outboundapp.RefreshModelFromResponse(ctx, &model, app)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
	tflog.Info(ctx, "Outbound application resource created")
}

func (r *outboundAppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading outbound application resource")

	var model outboundapp.Model
	if resp.Diagnostics.Append(req.State.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	id := model.ID.ValueString()
	app, err := infra.RetryOnRateLimit(ctx, func() (*descope.OutboundApp, error) {
		return r.management.OutboundApplication().LoadApplication(ctx, id)
	})
	if err != nil {
		if infra.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading outbound application", err.Error())
		return
	}

	outboundapp.RefreshModelFromResponse(ctx, &model, app)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
	tflog.Info(ctx, "Outbound application resource read")
}

func (r *outboundAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating outbound application resource")

	var plan outboundapp.Model
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	var state outboundapp.Model
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	plan.ID = state.ID
	appReq := outboundapp.ModelToOutboundApp(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	app, err := infra.RetryOnRateLimit(ctx, func() (*descope.OutboundApp, error) {
		return r.management.OutboundApplication().UpdateApplication(ctx, appReq, outboundapp.ClientSecretPtr(&plan))
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating outbound application", err.Error())
		return
	}

	outboundapp.RefreshModelFromResponse(ctx, &plan, app)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, "Outbound application resource updated")
}

func (r *outboundAppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting outbound application resource")

	var model outboundapp.Model
	if resp.Diagnostics.Append(req.State.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	id := model.ID.ValueString()
	err := infra.RetryOnRateLimitNoResult(ctx, func() error {
		return r.management.OutboundApplication().DeleteApplication(ctx, id)
	})
	if err != nil {
		if infra.IsNotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting outbound application", err.Error())
		return
	}

	tflog.Info(ctx, "Outbound application resource deleted")
}

func (r *outboundAppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing outbound application resource")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
