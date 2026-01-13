package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/j4ng5y/terraform-provider-elevenlabs/internal/client"
	"github.com/j4ng5y/terraform-provider-elevenlabs/internal/models"
)

var (
	_ resource.Resource                = &ConvAIToolResource{}
	_ resource.ResourceWithConfigure   = &ConvAIToolResource{}
	_ resource.ResourceWithImportState = &ConvAIToolResource{}
)

func NewConvAIToolResource() resource.Resource {
	return &ConvAIToolResource{}
}

type ConvAIToolResource struct {
	client *client.Client
}

type ConvAIToolResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func (r *ConvAIToolResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_convai_tool"
}

func (r *ConvAIToolResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Conversational AI Tool resource for ElevenLabs.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (r *ConvAIToolResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *ConvAIToolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ConvAIToolResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	addReq := &models.CreateConvAIToolRequest{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
	}

	tool, err := r.client.CreateConvAITool(addReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating ConvAI tool", err.Error())
		return
	}

	data.ID = types.StringValue(tool.ToolID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAIToolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ConvAIToolResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tool, err := r.client.GetConvAITool(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading ConvAI tool", err.Error())
		return
	}

	data.Name = types.StringValue(tool.Name)
	data.Description = types.StringValue(tool.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAIToolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ConvAIToolResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &models.CreateConvAIToolRequest{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
	}

	err := r.client.UpdateConvAITool(data.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating ConvAI tool", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAIToolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ConvAIToolResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteConvAITool(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting ConvAI tool", err.Error())
		return
	}
}

func (r *ConvAIToolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
