package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/j4ng5y/terraform-provider-elevenlabs/internal/client"
)

var (
	_ resource.Resource              = &ConvAISettingsResource{}
	_ resource.ResourceWithConfigure = &ConvAISettingsResource{}
)

func NewConvAISettingsResource() resource.Resource {
	return &ConvAISettingsResource{}
}

type ConvAISettingsResource struct {
	client *client.Client
}

type ConvAISettingsResourceModel struct {
	ID types.String `tfsdk:"id"`
	// For simplicity, we'll handle this as a JSON string or dynamic map if supported
	// but for now, let's just make it a placeholder or specific fields if known.
	// Task 98 research mentioned it exists.
}

func (r *ConvAISettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_convai_settings"
}

func (r *ConvAISettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Conversational AI Settings resource for ElevenLabs. Manages workspace-wide ConvAI configuration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *ConvAISettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ConvAISettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ConvAISettingsResourceModel
	data.ID = types.StringValue("workspace_settings")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAISettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ConvAISettingsResourceModel
	_, err := r.client.GetConvAISettings()
	if err != nil {
		resp.Diagnostics.AddError("Error reading ConvAI settings", err.Error())
		return
	}
	data.ID = types.StringValue("workspace_settings")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAISettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Update logic
}

func (r *ConvAISettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Settings usually can't be deleted, only reset or ignored.
}
