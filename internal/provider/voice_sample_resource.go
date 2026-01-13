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
	_ resource.Resource                = &VoiceSampleResource{}
	_ resource.ResourceWithConfigure   = &VoiceSampleResource{}
	_ resource.ResourceWithImportState = &VoiceSampleResource{}
)

func NewVoiceSampleResource() resource.Resource {
	return &VoiceSampleResource{}
}

type VoiceSampleResource struct {
	client *client.Client
}

type VoiceSampleResourceModel struct {
	ID       types.String `tfsdk:"id"`
	VoiceID  types.String `tfsdk:"voice_id"`
	FilePath types.String `tfsdk:"file_path"`
	FileName types.String `tfsdk:"file_name"`
}

func (r *VoiceSampleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_voice_sample"
}

func (r *VoiceSampleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Voice Sample resource for ElevenLabs. Allows adding and removing audio samples for custom voices.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"voice_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"file_path": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"file_name": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *VoiceSampleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VoiceSampleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data VoiceSampleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	addReq := &models.AddVoiceSampleRequest{
		FilePath: data.FilePath.ValueString(),
	}

	sample, err := r.client.AddVoiceSample(data.VoiceID.ValueString(), addReq)
	if err != nil {
		resp.Diagnostics.AddError("Error adding voice sample", err.Error())
		return
	}

	data.ID = types.StringValue(sample.SampleID)
	data.FileName = types.StringValue(sample.FileName)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VoiceSampleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// API doesn't have a single GET for sample, usually listed under Voice.
}

func (r *VoiceSampleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Handled by RequiresReplace
}

func (r *VoiceSampleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data VoiceSampleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteVoiceSample(data.VoiceID.ValueString(), data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting voice sample", err.Error())
		return
	}
}

func (r *VoiceSampleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
