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
	_ resource.Resource                = &AudioNativeResource{}
	_ resource.ResourceWithConfigure   = &AudioNativeResource{}
	_ resource.ResourceWithImportState = &AudioNativeResource{}
)

func NewAudioNativeResource() resource.Resource {
	return &AudioNativeResource{}
}

type AudioNativeResource struct {
	client *client.Client
}

type AudioNativeResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	FilePath        types.String `tfsdk:"file_path"`
	VoiceID         types.String `tfsdk:"voice_id"`
	ModelID         types.String `tfsdk:"model_id"`
	Title           types.String `tfsdk:"title"`
	Author          types.String `tfsdk:"author"`
	TextColor       types.String `tfsdk:"text_color"`
	BackgroundColor types.String `tfsdk:"background_color"`
	AutoConvert     types.Bool   `tfsdk:"auto_convert"`
	HTMLSnippet     types.String `tfsdk:"html_snippet"`
	Status          types.String `tfsdk:"status"`
}

func (r *AudioNativeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_audio_native"
}

func (r *AudioNativeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Audio Native resource for ElevenLabs. Allows creating an automated TTS player for websites.",
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
			"file_path": schema.StringAttribute{
				Required: true,
			},
			"voice_id": schema.StringAttribute{
				Optional: true,
			},
			"model_id": schema.StringAttribute{
				Optional: true,
			},
			"title": schema.StringAttribute{
				Optional: true,
			},
			"author": schema.StringAttribute{
				Optional: true,
			},
			"text_color": schema.StringAttribute{
				Optional: true,
			},
			"background_color": schema.StringAttribute{
				Optional: true,
			},
			"auto_convert": schema.BoolAttribute{
				Optional: true,
			},
			"html_snippet": schema.StringAttribute{
				Computed: true,
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *AudioNativeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AudioNativeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AudioNativeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	addReq := &models.CreateAudioNativeRequest{
		Name:            data.Name.ValueString(),
		FilePath:        data.FilePath.ValueString(),
		VoiceID:         data.VoiceID.ValueString(),
		ModelID:         data.ModelID.ValueString(),
		Title:           data.Title.ValueString(),
		Author:          data.Author.ValueString(),
		TextColor:       data.TextColor.ValueString(),
		BackgroundColor: data.BackgroundColor.ValueString(),
		AutoConvert:     data.AutoConvert.ValueBool(),
	}

	project, err := r.client.CreateAudioNative(addReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating audio native project", err.Error())
		return
	}

	data.ID = types.StringValue(project.ProjectID)
	data.HTMLSnippet = types.StringValue(project.HTMLSnippet)

	settings, err := r.client.GetAudioNativeSettings(project.ProjectID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading audio native settings", err.Error())
		return
	}
	data.Title = types.StringValue(settings.Title)
	data.Author = types.StringValue(settings.Author)
	data.TextColor = types.StringValue(settings.TextColor)
	data.BackgroundColor = types.StringValue(settings.BackgroundColor)
	data.Status = types.StringValue(settings.Status)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AudioNativeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AudioNativeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	settings, err := r.client.GetAudioNativeSettings(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading audio native settings", err.Error())
		return
	}

	data.Title = types.StringValue(settings.Title)
	data.Author = types.StringValue(settings.Author)
	data.TextColor = types.StringValue(settings.TextColor)
	data.BackgroundColor = types.StringValue(settings.BackgroundColor)
	data.Status = types.StringValue(settings.Status)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AudioNativeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Update logic could involve updating content or settings if endpoints exist.
	resp.Diagnostics.AddWarning("Update limited", "Updating audio native projects might requires replacement.")
}

func (r *AudioNativeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AudioNativeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Assuming standard project delete works since Audio Native returns a project ID
	err := r.client.DeleteProject(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting audio native project", err.Error())
		return
	}
}

func (r *AudioNativeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
