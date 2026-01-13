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
	_ resource.Resource                = &PronunciationDictionaryResource{}
	_ resource.ResourceWithConfigure   = &PronunciationDictionaryResource{}
	_ resource.ResourceWithImportState = &PronunciationDictionaryResource{}
)

func NewPronunciationDictionaryResource() resource.Resource {
	return &PronunciationDictionaryResource{}
}

type PronunciationDictionaryResource struct {
	client *client.Client
}

type PronunciationDictionaryResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	LatestVersionID types.String `tfsdk:"latest_version_id"`
	FilePath        types.String `tfsdk:"file_path"`
	Rules           []RuleModel  `tfsdk:"rules"`
}

type RuleModel struct {
	Type            types.String `tfsdk:"type"`
	StringToReplace types.String `tfsdk:"string_to_replace"`
	Alias           types.String `tfsdk:"alias"`
	Phoneme         types.String `tfsdk:"phoneme"`
	Alphabet        types.String `tfsdk:"alphabet"`
}

func (r *PronunciationDictionaryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pronunciation_dictionary"
}

func (r *PronunciationDictionaryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Pronunciation Dictionary resource for ElevenLabs. Allows defining custom pronunciations for words.",
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
				Optional: true,
			},
			"latest_version_id": schema.StringAttribute{
				Computed: true,
			},
			"file_path": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Path to a .pls file. Mutually exclusive with `rules`.",
			},
			"rules": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Type of rule: `alias` or `phoneme`.",
						},
						"string_to_replace": schema.StringAttribute{
							Required: true,
						},
						"alias": schema.StringAttribute{
							Optional: true,
						},
						"phoneme": schema.StringAttribute{
							Optional: true,
						},
						"alphabet": schema.StringAttribute{
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func (r *PronunciationDictionaryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PronunciationDictionaryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PronunciationDictionaryResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var dict *models.PronunciationDictionary
	var err error

	if !data.FilePath.IsNull() {
		addReq := &models.AddPronunciationDictionaryFromFileRequest{
			Name:        data.Name.ValueString(),
			Description: data.Description.ValueString(),
			FilePath:    data.FilePath.ValueString(),
		}
		dict, err = r.client.AddPronunciationDictionaryFromFile(addReq)
	} else if len(data.Rules) > 0 {
		rules := make([]models.PronunciationRule, len(data.Rules))
		for i, rule := range data.Rules {
			rules[i] = models.PronunciationRule{
				Type:            rule.Type.ValueString(),
				StringToReplace: rule.StringToReplace.ValueString(),
				Alias:           rule.Alias.ValueString(),
				Phoneme:         rule.Phoneme.ValueString(),
				Alphabet:        rule.Alphabet.ValueString(),
			}
		}
		addReq := &models.AddPronunciationDictionaryFromRulesRequest{
			Name:        data.Name.ValueString(),
			Description: data.Description.ValueString(),
			Rules:       rules,
		}
		dict, err = r.client.AddPronunciationDictionaryFromRules(addReq)
	} else {
		resp.Diagnostics.AddError("Invalid Configuration", "Either `file_path` or `rules` must be provided.")
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Error creating pronunciation dictionary", err.Error())
		return
	}

	data.ID = types.StringValue(dict.ID)
	data.LatestVersionID = types.StringValue(dict.LatestVersionID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PronunciationDictionaryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PronunciationDictionaryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dict, err := r.client.GetPronunciationDictionary(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading pronunciation dictionary", err.Error())
		return
	}

	data.Name = types.StringValue(dict.Name)
	data.LatestVersionID = types.StringValue(dict.LatestVersionID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PronunciationDictionaryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// The ElevenLabs API doesn't seem to have a direct "Update Metadata" for dictionaries
	// but rules can be added/removed. For now, we'll mark it as needing replacement
	// or implement rule delta logic if requested.
	resp.Diagnostics.AddWarning("Update limited", "Updating pronunciation dictionaries might require replacement or manual rule management via API.")
}

func (r *PronunciationDictionaryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PronunciationDictionaryResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.ArchivePronunciationDictionary(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error archiving pronunciation dictionary", err.Error())
		return
	}
}

func (r *PronunciationDictionaryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
