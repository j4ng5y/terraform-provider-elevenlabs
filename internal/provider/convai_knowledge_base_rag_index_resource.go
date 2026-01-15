package provider

import (
	"context"
	"fmt"
	"strings"

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
	_ resource.Resource                = &ConvAIKnowledgeBaseRAGIndexResource{}
	_ resource.ResourceWithConfigure   = &ConvAIKnowledgeBaseRAGIndexResource{}
	_ resource.ResourceWithImportState = &ConvAIKnowledgeBaseRAGIndexResource{}
)

func NewConvAIKnowledgeBaseRAGIndexResource() resource.Resource {
	return &ConvAIKnowledgeBaseRAGIndexResource{}
}

type ConvAIKnowledgeBaseRAGIndexResource struct {
	client *client.Client
}

type ConvAIKnowledgeBaseRAGIndexResourceModel struct {
	ID                 types.String  `tfsdk:"id"`
	DocumentationID    types.String  `tfsdk:"documentation_id"`
	Model              types.String  `tfsdk:"model"`
	Status             types.String  `tfsdk:"status"`
	ProgressPercentage types.Float64 `tfsdk:"progress_percentage"`
	UsedBytes          types.Int64   `tfsdk:"used_bytes"`
}

func (r *ConvAIKnowledgeBaseRAGIndexResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_convai_knowledge_base_rag_index"
}

func (r *ConvAIKnowledgeBaseRAGIndexResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Conversational AI Knowledge Base RAG index resource for ElevenLabs.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"documentation_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"model": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
			"progress_percentage": schema.Float64Attribute{
				Computed: true,
			},
			"used_bytes": schema.Int64Attribute{
				Computed: true,
			},
		},
	}
}

func (r *ConvAIKnowledgeBaseRAGIndexResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ConvAIKnowledgeBaseRAGIndexResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ConvAIKnowledgeBaseRAGIndexResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	addReq := &models.RAGIndexRequest{
		Model: data.Model.ValueString(),
	}

	index, err := r.client.CreateConvAIKnowledgeBaseRAGIndex(data.DocumentationID.ValueString(), addReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating ConvAI knowledge base RAG index", err.Error())
		return
	}

	applyRAGIndexToState(&data, index)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAIKnowledgeBaseRAGIndexResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ConvAIKnowledgeBaseRAGIndexResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	list, err := r.client.GetConvAIKnowledgeBaseRAGIndexes(data.DocumentationID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading ConvAI knowledge base RAG indexes", err.Error())
		return
	}

	index := findRAGIndexByID(list.Indexes, data.ID.ValueString())
	if index == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	applyRAGIndexToState(&data, index)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAIKnowledgeBaseRAGIndexResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ConvAIKnowledgeBaseRAGIndexResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAIKnowledgeBaseRAGIndexResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ConvAIKnowledgeBaseRAGIndexResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteConvAIKnowledgeBaseRAGIndex(data.DocumentationID.ValueString(), data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting ConvAI knowledge base RAG index", err.Error())
		return
	}
}

func (r *ConvAIKnowledgeBaseRAGIndexResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError("Invalid import ID", "expected ID in format documentation_id:rag_index_id")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("documentation_id"), parts[0])...)
}

func findRAGIndexByID(indexes []models.RAGDocumentIndexResponse, id string) *models.RAGDocumentIndexResponse {
	for i := range indexes {
		if indexes[i].ID == id {
			return &indexes[i]
		}
	}
	return nil
}

func applyRAGIndexToState(data *ConvAIKnowledgeBaseRAGIndexResourceModel, index *models.RAGDocumentIndexResponse) {
	data.ID = types.StringValue(index.ID)
	data.Model = types.StringValue(index.Model)
	data.Status = types.StringValue(index.Status)
	data.ProgressPercentage = types.Float64Value(index.ProgressPercentage)
	data.UsedBytes = types.Int64Value(index.DocumentModelIndexUsage.UsedBytes)
}
