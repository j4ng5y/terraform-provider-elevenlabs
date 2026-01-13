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
	_ resource.Resource                = &ConvAIKnowledgeBaseResource{}
	_ resource.ResourceWithConfigure   = &ConvAIKnowledgeBaseResource{}
	_ resource.ResourceWithImportState = &ConvAIKnowledgeBaseResource{}
)

func NewConvAIKnowledgeBaseResource() resource.Resource {
	return &ConvAIKnowledgeBaseResource{}
}

type ConvAIKnowledgeBaseResource struct {
	client *client.Client
}

type ConvAIKnowledgeBaseResourceModel struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	URL      types.String `tfsdk:"url"`
	Content  types.String `tfsdk:"content"`
	FilePath types.String `tfsdk:"file_path"`
	Type     types.String `tfsdk:"type"`
	Status   types.String `tfsdk:"status"`
}

func (r *ConvAIKnowledgeBaseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_convai_knowledge_base"
}

func (r *ConvAIKnowledgeBaseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Conversational AI Knowledge Base resource for ElevenLabs.",
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
			"url": schema.StringAttribute{
				Optional: true,
			},
			"content": schema.StringAttribute{
				Optional: true,
			},
			"file_path": schema.StringAttribute{
				Optional: true,
			},
			"type": schema.StringAttribute{
				Computed: true,
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *ConvAIKnowledgeBaseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ConvAIKnowledgeBaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ConvAIKnowledgeBaseResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	addReq := &models.CreateConvAIKnowledgeBaseRequest{
		Name:     data.Name.ValueString(),
		URL:      data.URL.ValueString(),
		Content:  data.Content.ValueString(),
		FilePath: data.FilePath.ValueString(),
	}

	kb, err := r.client.CreateConvAIKnowledgeBase(addReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating ConvAI knowledge base", err.Error())
		return
	}

	data.ID = types.StringValue(kb.DocumentationID)
	data.Type = types.StringValue(kb.Type)
	data.Status = types.StringValue(kb.Status)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAIKnowledgeBaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ConvAIKnowledgeBaseResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	kb, err := r.client.GetConvAIKnowledgeBase(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading ConvAI knowledge base", err.Error())
		return
	}

	data.Name = types.StringValue(kb.Name)
	data.Type = types.StringValue(kb.Type)
	data.Status = types.StringValue(kb.Status)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAIKnowledgeBaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning("Update limited", "Updating ConvAI knowledge base might require replacement.")
}

func (r *ConvAIKnowledgeBaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ConvAIKnowledgeBaseResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteConvAIKnowledgeBase(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting ConvAI knowledge base", err.Error())
		return
	}
}

func (r *ConvAIKnowledgeBaseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
