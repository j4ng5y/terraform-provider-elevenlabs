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
	_ resource.Resource                = &PVCVoiceResource{}
	_ resource.ResourceWithConfigure   = &PVCVoiceResource{}
	_ resource.ResourceWithImportState = &PVCVoiceResource{}
)

// PVCVoiceResourceModel describes the resource data model.
type PVCVoiceResourceModel struct {
	ID           types.String   `tfsdk:"id"`
	Name         types.String   `tfsdk:"name"`
	Language     types.String   `tfsdk:"language"`
	Description  types.String   `tfsdk:"description"`
	Labels       types.Map      `tfsdk:"labels"`
	State        types.String   `tfsdk:"state"`
	Verification types.String   `tfsdk:"verification"`
	Settings     *VoiceSettings `tfsdk:"settings"`
}

func NewPVCVoiceResource() resource.Resource {
	return &PVCVoiceResource{}
}

// PVCVoiceResource defines the resource implementation.
type PVCVoiceResource struct {
	client *client.Client
}

func (r *PVCVoiceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pvc_voice"
}

func (r *PVCVoiceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "PVC Voice resource for ElevenLabs. Allows creating and managing Professional Voice Cloning voices.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for the PVC voice.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the PVC voice.",
			},
			"language": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The language code for the PVC voice (e.g., 'en', 'es', 'fr').",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "A description of the PVC voice.",
			},
			"labels": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Labels associated with the PVC voice.",
			},
			"state": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The current training state of the PVC voice.",
			},
			"verification": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The verification status of the PVC voice.",
			},
			"settings": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"stability": schema.Float64Attribute{
						Optional: true,
						Computed: true,
					},
					"similarity_boost": schema.Float64Attribute{
						Optional: true,
						Computed: true,
					},
					"style": schema.Float64Attribute{
						Optional: true,
						Computed: true,
					},
					"use_speaker_boost": schema.BoolAttribute{
						Optional: true,
						Computed: true,
					},
				},
			},
		},
	}
}

func (r *PVCVoiceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PVCVoiceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PVCVoiceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &models.CreatePVCVoiceRequest{
		Name:        data.Name.ValueString(),
		Language:    data.Language.ValueString(),
		Description: data.Description.ValueString(),
	}

	if !data.Labels.IsNull() {
		labels := make(map[string]string)
		resp.Diagnostics.Append(data.Labels.ElementsAs(ctx, &labels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		createReq.Labels = labels
	}

	voice, err := r.client.CreatePVCVoice(createReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create PVC voice, got error: %s", err))
		return
	}

	data.ID = types.StringValue(voice.VoiceID)
	data.Name = types.StringValue(voice.Name)
	data.Language = types.StringValue(voice.Language)
	data.Description = types.StringValue(voice.Description)
	data.State = types.StringValue(voice.State)
	data.Verification = types.StringValue(voice.Verification)

	if voice.Labels != nil {
		labels, diags := types.MapValueFrom(ctx, types.StringType, voice.Labels)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.Labels = labels
	} else if data.Labels.IsUnknown() {
		data.Labels = types.MapNull(types.StringType)
	}

	if voice.Settings != nil {
		data.Settings = &VoiceSettings{
			Stability:       types.Float64Value(voice.Settings.Stability),
			SimilarityBoost: types.Float64Value(voice.Settings.SimilarityBoost),
			Style:           types.Float64Value(voice.Settings.Style),
			UseSpeakerBoost: types.BoolValue(voice.Settings.UseSpeakerBoost),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PVCVoiceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PVCVoiceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	voice, err := r.client.GetPVCVoice(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read PVC voice, got error: %s", err))
		return
	}

	data.Name = types.StringValue(voice.Name)
	data.Language = types.StringValue(voice.Language)
	data.Description = types.StringValue(voice.Description)
	data.State = types.StringValue(voice.State)
	data.Verification = types.StringValue(voice.Verification)

	if voice.Labels != nil {
		labels, diags := types.MapValueFrom(ctx, types.StringType, voice.Labels)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.Labels = labels
	} else if data.Labels.IsUnknown() {
		data.Labels = types.MapNull(types.StringType)
	}

	if voice.Settings != nil {
		data.Settings = &VoiceSettings{
			Stability:       types.Float64Value(voice.Settings.Stability),
			SimilarityBoost: types.Float64Value(voice.Settings.SimilarityBoost),
			Style:           types.Float64Value(voice.Settings.Style),
			UseSpeakerBoost: types.BoolValue(voice.Settings.UseSpeakerBoost),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PVCVoiceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data PVCVoiceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &models.UpdatePVCVoiceRequest{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
	}

	if !data.Labels.IsNull() {
		labels := make(map[string]string)
		resp.Diagnostics.Append(data.Labels.ElementsAs(ctx, &labels, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		updateReq.Labels = labels
	}

	err := r.client.UpdatePVCVoice(data.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update PVC voice, got error: %s", err))
		return
	}

	// Read the updated voice to get current state
	voice, err := r.client.GetPVCVoice(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read updated PVC voice, got error: %s", err))
		return
	}

	data.State = types.StringValue(voice.State)
	data.Verification = types.StringValue(voice.Verification)

	if voice.Labels != nil {
		labels, diags := types.MapValueFrom(ctx, types.StringType, voice.Labels)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.Labels = labels
	} else if data.Labels.IsUnknown() {
		data.Labels = types.MapNull(types.StringType)
	}

	if voice.Settings != nil {
		data.Settings = &VoiceSettings{
			Stability:       types.Float64Value(voice.Settings.Stability),
			SimilarityBoost: types.Float64Value(voice.Settings.SimilarityBoost),
			Style:           types.Float64Value(voice.Settings.Style),
			UseSpeakerBoost: types.BoolValue(voice.Settings.UseSpeakerBoost),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PVCVoiceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PVCVoiceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeletePVCVoice(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete PVC voice, got error: %s", err))
		return
	}
}

func (r *PVCVoiceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
