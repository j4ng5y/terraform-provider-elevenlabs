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
	_ resource.Resource                = &DubbingProjectResource{}
	_ resource.ResourceWithConfigure   = &DubbingProjectResource{}
	_ resource.ResourceWithImportState = &DubbingProjectResource{}
)

func NewDubbingProjectResource() resource.Resource {
	return &DubbingProjectResource{}
}

type DubbingProjectResource struct {
	client *client.Client
}

type DubbingProjectResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	SourceURL       types.String `tfsdk:"source_url"`
	FilePath        types.String `tfsdk:"file_path"`
	SourceLang      types.String `tfsdk:"source_lang"`
	TargetLang      types.String `tfsdk:"target_lang"`
	NumSpeakers     types.Int64  `tfsdk:"num_speakers"`
	Watermark       types.Bool   `tfsdk:"watermark"`
	DubbingStudio   types.Bool   `tfsdk:"dubbing_studio"`
	Mode            types.String `tfsdk:"mode"`
	Status          types.String `tfsdk:"status"`
	TargetLanguages types.List   `tfsdk:"target_languages"`
}

func (r *DubbingProjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dubbing_project"
}

func (r *DubbingProjectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Dubbing Project resource for ElevenLabs. Allows translating and voicing video or audio files.",
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
			"source_url": schema.StringAttribute{
				Optional: true,
			},
			"file_path": schema.StringAttribute{
				Optional: true,
			},
			"source_lang": schema.StringAttribute{
				Optional: true,
			},
			"target_lang": schema.StringAttribute{
				Required: true,
			},
			"num_speakers": schema.Int64Attribute{
				Optional: true,
			},
			"watermark": schema.BoolAttribute{
				Optional: true,
			},
			"dubbing_studio": schema.BoolAttribute{
				Optional: true,
			},
			"mode": schema.StringAttribute{
				Optional: true,
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
			"target_languages": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

func (r *DubbingProjectResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DubbingProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DubbingProjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	addReq := &models.CreateDubbingRequest{
		Name:          data.Name.ValueString(),
		SourceURL:     data.SourceURL.ValueString(),
		SourceLang:    data.SourceLang.ValueString(),
		TargetLang:    data.TargetLang.ValueString(),
		NumSpeakers:   int(data.NumSpeakers.ValueInt64()),
		Watermark:     data.Watermark.ValueBool(),
		DubbingStudio: data.DubbingStudio.ValueBool(),
		Mode:          data.Mode.ValueString(),
		FilePath:      data.FilePath.ValueString(),
	}

	project, err := r.client.CreateDubbingProject(addReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating dubbing project", err.Error())
		return
	}

	data.ID = types.StringValue(project.DubbingID)
	data.Status = types.StringValue(project.Status)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DubbingProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DubbingProjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	project, err := r.client.GetDubbingProject(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading dubbing project", err.Error())
		return
	}

	data.Name = types.StringValue(project.Name)
	data.Status = types.StringValue(project.Status)

	targetLanguages, diag := types.ListValueFrom(ctx, types.StringType, project.TargetLanguages)
	resp.Diagnostics.Append(diag...)
	data.TargetLanguages = targetLanguages

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DubbingProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning("Update limited", "Updating dubbing projects might requires replacement.")
}

func (r *DubbingProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DubbingProjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteDubbingProject(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting dubbing project", err.Error())
		return
	}
}

func (r *DubbingProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
