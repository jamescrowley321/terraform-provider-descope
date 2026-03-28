package resources

import (
	"context"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/go-sdk/descope/sdk"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jamescrowley321/terraform-provider-descope/internal/infra"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/list"
)

var (
	_ resource.Resource                = &listResource{}
	_ resource.ResourceWithConfigure   = &listResource{}
	_ resource.ResourceWithImportState = &listResource{}
)

func NewListResource() resource.Resource {
	return &listResource{}
}

type listResource struct {
	management sdk.Management
}

func (r *listResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if data, ok := req.ProviderData.(*infra.ProviderData); ok {
		r.management = data.Management
	}
}

func (r *listResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_list"
}

func (r *listResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Descope list for IP allowlisting/denylisting or text-based filtering.",
		Attributes:  list.Attributes,
	}
}

func (r *listResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating list resource")

	var model list.Model
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	listReq := list.ModelToRequest(ctx, &model, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := infra.RetryOnRateLimit(ctx, func() (*descope.List, error) {
		return r.management.List().Create(ctx, listReq)
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating list", err.Error())
		return
	}

	list.RefreshModelFromResponse(ctx, &model, result)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
	tflog.Info(ctx, "List resource created")
}

func (r *listResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading list resource")

	var model list.Model
	if resp.Diagnostics.Append(req.State.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	id := model.ID.ValueString()
	result, err := infra.RetryOnRateLimit(ctx, func() (*descope.List, error) {
		return r.management.List().Load(ctx, id)
	})
	if err != nil {
		if infra.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading list", err.Error())
		return
	}

	list.RefreshModelFromResponse(ctx, &model, result)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
	tflog.Info(ctx, "List resource read")
}

func (r *listResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating list resource")

	var plan list.Model
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	var state list.Model
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	plan.ID = state.ID
	listReq := list.ModelToRequest(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := infra.RetryOnRateLimit(ctx, func() (*descope.List, error) {
		return r.management.List().Update(ctx, plan.ID.ValueString(), listReq)
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating list", err.Error())
		return
	}

	list.RefreshModelFromResponse(ctx, &plan, result)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, "List resource updated")
}

func (r *listResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting list resource")

	var model list.Model
	if resp.Diagnostics.Append(req.State.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	id := model.ID.ValueString()
	err := infra.RetryOnRateLimitNoResult(ctx, func() error {
		return r.management.List().Delete(ctx, id)
	})
	if err != nil {
		if infra.IsNotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting list", err.Error())
		return
	}

	tflog.Info(ctx, "List resource deleted")
}

func (r *listResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing list resource")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
