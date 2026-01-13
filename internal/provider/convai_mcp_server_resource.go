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
	_ resource.Resource                = &ConvAIMCPServerResource{}
	_ resource.ResourceWithConfigure   = &ConvAIMCPServerResource{}
	_ resource.ResourceWithImportState = &ConvAIMCPServerResource{}
)

func NewConvAIMCPServerResource() resource.Resource {
	return &ConvAIMCPServerResource{}
}

type ConvAIMCPServerResource struct {
	client *client.Client
}

type ConvAIMCPServerResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	URL  types.String `tfsdk:"url"`
}

func (r *ConvAIMCPServerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_convai_mcp_server"
}

func (r *ConvAIMCPServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Conversational AI MCP Server resource for ElevenLabs.",
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
			"url": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (r *ConvAIMCPServerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ConvAIMCPServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ConvAIMCPServerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	addReq := &models.CreateConvAIMCPServerRequest{
		Name: data.Name.ValueString(),
		URL:  data.URL.ValueString(),
	}

	server, err := r.client.CreateConvAIMCPServer(addReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating ConvAI MCP server", err.Error())
		return
	}

	data.ID = types.StringValue(server.MCPServerID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAIMCPServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// API Read logic
}

func (r *ConvAIMCPServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ConvAIMCPServerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &models.CreateConvAIMCPServerRequest{
		Name: data.Name.ValueString(),
		URL:  data.URL.ValueString(),
	}

	err := r.client.UpdateConvAIMCPServer(data.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating ConvAI MCP server", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAIMCPServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ConvAIMCPServerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteConvAIMCPServer(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting ConvAI MCP server", err.Error())
		return
	}
}

func (r *ConvAIMCPServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
