package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/j4ng5y/terraform-provider-elevenlabs/internal/client"
	"github.com/j4ng5y/terraform-provider-elevenlabs/internal/models"
)

var (
	_ datasource.DataSource              = &ConvAIKnowledgeBasesDataSource{}
	_ datasource.DataSourceWithConfigure = &ConvAIKnowledgeBasesDataSource{}
)

var knowledgeBaseDocumentAttrTypes = map[string]attr.Type{
	"documentation_id": types.StringType,
	"name":             types.StringType,
	"type":             types.StringType,
	"status":           types.StringType,
	"metadata_json":    types.StringType,
	"supported_usages": types.ListType{ElemType: types.StringType},
	"folder_parent_id": types.StringType,
	"folder_path":      types.ListType{ElemType: types.StringType},
	"access_info_json": types.StringType,
}

var knowledgeBaseDocumentObjectType = types.ObjectType{AttrTypes: knowledgeBaseDocumentAttrTypes}

func NewConvAIKnowledgeBasesDataSource() datasource.DataSource {
	return &ConvAIKnowledgeBasesDataSource{}
}

type ConvAIKnowledgeBasesDataSource struct {
	client *client.Client
}

type convaiKnowledgeBasesDataSourceModel struct {
	Search                 types.String `tfsdk:"search"`
	ShowOnlyOwnedDocuments types.Bool   `tfsdk:"show_only_owned_documents"`
	Types                  types.List   `tfsdk:"types"`
	PageSize               types.Int64  `tfsdk:"page_size"`
	Cursor                 types.String `tfsdk:"cursor"`

	Documents  types.List   `tfsdk:"documents"`
	HasMore    types.Bool   `tfsdk:"has_more"`
	NextCursor types.String `tfsdk:"next_cursor"`
}

func (d *ConvAIKnowledgeBasesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_convai_knowledge_bases"
}

func (d *ConvAIKnowledgeBasesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches ElevenLabs ConvAI knowledge base documents.",
		Attributes: map[string]schema.Attribute{
			"search": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Full-text search to filter documents by name.",
			},
			"show_only_owned_documents": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Limit results to documents owned by the authenticated workspace.",
			},
			"types": schema.ListAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Restrict results to the provided document types (e.g., `url`, `file`, `text`, `folder`).",
			},
			"page_size": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Number of records to request per page (max 100).",
			},
			"cursor": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Opaque pagination cursor returned by the API.",
			},
			"documents": schema.ListAttribute{
				Computed:            true,
				ElementType:         knowledgeBaseDocumentObjectType,
				MarkdownDescription: "Knowledge base entries that satisfy the supplied filters.",
			},
			"has_more": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "True when another page of results is available.",
			},
			"next_cursor": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Cursor to pass to the next query when `has_more` is true.",
			},
		},
	}
}

func (d *ConvAIKnowledgeBasesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *ConvAIKnowledgeBasesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data convaiKnowledgeBasesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := &models.ListConvAIKnowledgeBaseDocumentsParams{}

	if !data.PageSize.IsNull() && !data.PageSize.IsUnknown() {
		value := int(data.PageSize.ValueInt64())
		params.PageSize = &value
	}

	if !data.Search.IsNull() && !data.Search.IsUnknown() {
		params.Search = data.Search.ValueString()
	}

	if !data.Cursor.IsNull() && !data.Cursor.IsUnknown() {
		params.Cursor = data.Cursor.ValueString()
	}

	if !data.ShowOnlyOwnedDocuments.IsNull() && !data.ShowOnlyOwnedDocuments.IsUnknown() {
		value := data.ShowOnlyOwnedDocuments.ValueBool()
		params.ShowOnlyOwned = &value
	}

	typeFilters, diags := stringSliceFromList(ctx, data.Types)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	params.Types = typeFilters

	result, err := d.client.ListConvAIKnowledgeBaseDocuments(params)
	if err != nil {
		resp.Diagnostics.AddError("Error listing ConvAI knowledge bases", err.Error())
		return
	}

	docs, diags := flattenKnowledgeBaseDocuments(ctx, result.Documents)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Documents = docs
	data.HasMore = types.BoolValue(result.HasMore)
	if result.NextCursor == "" {
		data.NextCursor = types.StringNull()
	} else {
		data.NextCursor = types.StringValue(result.NextCursor)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenKnowledgeBaseDocuments(ctx context.Context, docs []models.ConvAIKnowledgeBase) (types.List, diag.Diagnostics) {
	if len(docs) == 0 {
		return types.ListNull(knowledgeBaseDocumentObjectType), nil
	}

	values := make([]attr.Value, 0, len(docs))
	for _, doc := range docs {
		folderNames := make([]string, 0, len(doc.FolderPath))
		for _, segment := range doc.FolderPath {
			folderNames = append(folderNames, segment.FolderName)
		}

		supported, diags := stringsToListValue(ctx, doc.SupportedUsages)
		if diags.HasError() {
			return types.ListNull(knowledgeBaseDocumentObjectType), diags
		}

		folders, diags := stringsToListValue(ctx, folderNames)
		if diags.HasError() {
			return types.ListNull(knowledgeBaseDocumentObjectType), diags
		}

		obj, objDiags := types.ObjectValue(knowledgeBaseDocumentAttrTypes, map[string]attr.Value{
			"documentation_id": types.StringValue(doc.DocumentationID),
			"name":             types.StringValue(doc.Name),
			"type":             types.StringValue(doc.Type),
			"status":           types.StringValue(doc.Status),
			"metadata_json":    jsonStringValue(doc.Metadata),
			"supported_usages": supported,
			"folder_parent_id": optionalStringValue(doc.FolderParentID),
			"folder_path":      folders,
			"access_info_json": jsonStringValue(doc.AccessInfo),
		})
		if objDiags.HasError() {
			return types.ListNull(knowledgeBaseDocumentObjectType), objDiags
		}

		values = append(values, obj)
	}

	return types.ListValue(knowledgeBaseDocumentObjectType, values)
}
