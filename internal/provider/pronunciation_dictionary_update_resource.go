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
	_ resource.Resource                = &PronunciationDictionaryUpdateResource{}
	_ resource.ResourceWithConfigure   = &PronunciationDictionaryUpdateResource{}
	_ resource.ResourceWithImportState = &PronunciationDictionaryUpdateResource{}
)

func NewPronunciationDictionaryUpdateResource() resource.Resource {
	return &PronunciationDictionaryUpdateResource{}
}

type PronunciationDictionaryUpdateResource struct {
	client *client.Client
}

type PronunciationDictionaryUpdateResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Archived  types.Bool   `tfsdk:"archived"`
	VersionID types.String `tfsdk:"version_id"`
}

func (r *PronunciationDictionaryUpdateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pronunciation_dictionary_update"
}

func (r *PronunciationDictionaryUpdateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Update resource for ElevenLabs Pronunciation Dictionary. Allows updating name and archived status.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The name of the pronunciation dictionary.",
			},
			"archived": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether to archive the pronunciation dictionary.",
			},
			"version_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The version ID of the pronunciation dictionary after update.",
			},
		},
	}
}

func (r *PronunciationDictionaryUpdateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PronunciationDictionaryUpdateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// This resource doesn't create new resources, it updates existing ones
	resp.Diagnostics.AddError(
		"Invalid Usage",
		"This resource is for updating existing pronunciation dictionaries. Use `elevenlabs_pronunciation_dictionary` to create new ones.",
	)
}

func (r *PronunciationDictionaryUpdateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PronunciationDictionaryUpdateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the current state to verify the resource still exists
	dict, err := r.client.GetPronunciationDictionary(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading pronunciation dictionary", err.Error())
		return
	}

	data.Name = types.StringValue(dict.Name)
	data.VersionID = types.StringValue(dict.LatestVersionID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PronunciationDictionaryUpdateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan PronunciationDictionaryUpdateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state PronunciationDictionaryUpdateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var name string
	if !plan.Name.IsNull() && plan.Name.ValueString() != "" {
		name = plan.Name.ValueString()
	} else if !state.Name.IsNull() {
		name = state.Name.ValueString()
	}

	var archived *bool
	if !plan.Archived.IsNull() {
		archivedValue := plan.Archived.ValueBool()
		archived = &archivedValue
	}

	err := r.client.UpdatePronunciationDictionary(plan.ID.ValueString(), name, archived)
	if err != nil {
		resp.Diagnostics.AddError("Error updating pronunciation dictionary", err.Error())
		return
	}

	// Get the updated dictionary to get the new version ID
	dict, err := r.client.GetPronunciationDictionary(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading updated pronunciation dictionary", err.Error())
		return
	}

	plan.VersionID = types.StringValue(dict.LatestVersionID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *PronunciationDictionaryUpdateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This resource doesn't delete the underlying resource
	// Use the elevenlabs_pronunciation_dictionary resource for deletion
}

func (r *PronunciationDictionaryUpdateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
