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

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &VoiceResource{}
	_ resource.ResourceWithConfigure   = &VoiceResource{}
	_ resource.ResourceWithImportState = &VoiceResource{}
)

// VoiceResourceModel describes the resource data model.
type VoiceResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Labels      types.Map    `tfsdk:"labels"`
	Files       types.List   `tfsdk:"files"`
}

func NewVoiceResource() resource.Resource {
	return &VoiceResource{}
}

// VoiceResource defines the resource implementation.
type VoiceResource struct {
	client *client.Client
}

func (r *VoiceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_voice"
}

func (r *VoiceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Voice resource for ElevenLabs. Allows creating and managing custom voices.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for the voice.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the voice.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "A description of the voice.",
			},
			"labels": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Labels associated with the voice.",
			},
			"files": schema.ListAttribute{
				ElementType:         types.StringType,
				Required:            true,
				MarkdownDescription: "List of file paths for voice cloning. Required for creation.",
			},
		},
	}
}

func (r *VoiceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *VoiceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VoiceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var files []string
	resp.Diagnostics.Append(data.Files.ElementsAs(ctx, &files, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	labels := make(map[string]string)
	resp.Diagnostics.Append(data.Labels.ElementsAs(ctx, &labels, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	addReq := &models.AddVoiceRequest{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		Labels:      labels,
		Files:       files,
	}

	voice, err := r.client.AddVoice(addReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating voice", err.Error())
		return
	}

	data.ID = types.StringValue(voice.VoiceID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VoiceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data VoiceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	voice, err := r.client.GetVoice(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading voice", err.Error())
		return
	}

	data.Name = types.StringValue(voice.Name)
	// Note: Description is not returned by the API in the same way, we might need to handle labels
	if voice.Labels != nil {
		labelMap, diag := types.MapValueFrom(ctx, types.StringType, voice.Labels)
		resp.Diagnostics.Append(diag...)
		data.Labels = labelMap
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VoiceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data VoiceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var files []string
	resp.Diagnostics.Append(data.Files.ElementsAs(ctx, &files, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	labels := make(map[string]string)
	resp.Diagnostics.Append(data.Labels.ElementsAs(ctx, &labels, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &models.AddVoiceRequest{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		Labels:      labels,
		Files:       files,
	}

	err := r.client.EditVoice(data.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating voice", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VoiceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VoiceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteVoice(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting voice", err.Error())
		return
	}
}

func (r *VoiceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
