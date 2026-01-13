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
)

var (
	_ resource.Resource                = &SharedVoiceResource{}
	_ resource.ResourceWithConfigure   = &SharedVoiceResource{}
	_ resource.ResourceWithImportState = &SharedVoiceResource{}
)

func NewSharedVoiceResource() resource.Resource {
	return &SharedVoiceResource{}
}

type SharedVoiceResource struct {
	client *client.Client
}

type SharedVoiceResourceModel struct {
	ID           types.String `tfsdk:"id"`
	PublicUserID types.String `tfsdk:"public_user_id"`
	VoiceID      types.String `tfsdk:"voice_id"`
	Name         types.String `tfsdk:"name"`
}

func (r *SharedVoiceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_shared_voice"
}

func (r *SharedVoiceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Shared Voice resource for ElevenLabs. Allows adding a shared voice to your collection.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"public_user_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"voice_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "The name to give the voice in your collection.",
			},
		},
	}
}

func (r *SharedVoiceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SharedVoiceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SharedVoiceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	voiceID, err := r.client.AddSharedVoice(data.PublicUserID.ValueString(), data.VoiceID.ValueString(), data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error adding shared voice", err.Error())
		return
	}

	data.ID = types.StringValue(voiceID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SharedVoiceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SharedVoiceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	voice, err := r.client.GetVoice(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading shared voice", err.Error())
		return
	}

	data.Name = types.StringValue(voice.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SharedVoiceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// RequiresReplace
}

func (r *SharedVoiceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SharedVoiceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteVoice(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting shared voice from collection", err.Error())
		return
	}
}

func (r *SharedVoiceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
