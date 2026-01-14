package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/j4ng5y/terraform-provider-elevenlabs/internal/client"
	"github.com/j4ng5y/terraform-provider-elevenlabs/internal/models"
)

var (
	_ resource.Resource              = &PronunciationDictionaryRulesResource{}
	_ resource.ResourceWithConfigure = &PronunciationDictionaryRulesResource{}
)

func NewPronunciationDictionaryRulesResource() resource.Resource {
	return &PronunciationDictionaryRulesResource{}
}

type PronunciationDictionaryRulesResource struct {
	client *client.Client
}

type PronunciationDictionaryRulesResourceModel struct {
	DictionaryID types.String `tfsdk:"dictionary_id"`
	Rules        []RuleModel  `tfsdk:"rules"`
	Action       types.String `tfsdk:"action"`
}

func (r *PronunciationDictionaryRulesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pronunciation_dictionary_rules"
}

func (r *PronunciationDictionaryRulesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for managing pronunciation dictionary rules in ElevenLabs. Allows adding or removing rules from an existing dictionary.",
		Attributes: map[string]schema.Attribute{
			"dictionary_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the pronunciation dictionary to modify.",
			},
			"rules": schema.ListNestedAttribute{
				Required: true,
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
			"action": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Action to perform: `add` or `remove`.",
			},
		},
	}
}

func (r *PronunciationDictionaryRulesResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PronunciationDictionaryRulesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PronunciationDictionaryRulesResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rules := make([]models.PronunciationRule, 0, len(data.Rules))
	for _, rule := range data.Rules {
		rules = append(rules, models.PronunciationRule{
			StringToReplace: rule.StringToReplace.ValueString(),
			Type:            rule.Type.ValueString(),
		})
	}

	var err error
	action := data.Action.ValueString()
	switch action {
	case "add":
		err = r.client.AddPronunciationDictionaryRules(data.DictionaryID.ValueString(), rules)
	case "remove":
		err = r.client.RemovePronunciationDictionaryRules(data.DictionaryID.ValueString(), rules)
	default:
		resp.Diagnostics.AddError("Invalid Action", "Action must be 'add' or 'remove'.")
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Error modifying pronunciation dictionary rules", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PronunciationDictionaryRulesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data PronunciationDictionaryRulesResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

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

	var err error
	action := data.Action.ValueString()
	switch action {
	case "add":
		err = r.client.AddPronunciationDictionaryRules(data.DictionaryID.ValueString(), rules)
	case "remove":
		err = r.client.RemovePronunciationDictionaryRules(data.DictionaryID.ValueString(), rules)
	default:
		resp.Diagnostics.AddError("Invalid Action", "Action must be 'add' or 'remove'.")
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Error modifying pronunciation dictionary rules", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PronunciationDictionaryRulesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// This is an action resource, no persistent state needed
}

func (r *PronunciationDictionaryRulesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This is an action resource, no deletion needed
}
