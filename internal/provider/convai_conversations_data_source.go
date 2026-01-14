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
)

var (
	_ datasource.DataSource              = &ConvAIConversationsDataSource{}
	_ datasource.DataSourceWithConfigure = &ConvAIConversationsDataSource{}
)

var convaiConversationAttrTypes = map[string]attr.Type{
	"conversation_id": types.StringType,
	"name":            types.StringType,
	"agent_id":        types.StringType,
	"agent_name":      types.StringType,
	"created_at":      types.StringType,
}

var convaiConversationObjectType = types.ObjectType{AttrTypes: convaiConversationAttrTypes}

func NewConvAIConversationsDataSource() datasource.DataSource {
	return &ConvAIConversationsDataSource{}
}

type ConvAIConversationsDataSource struct {
	client *client.Client
}

type convaiConversationsDataSourceModel struct {
	Conversations types.List `tfsdk:"conversations"`
}

func (d *ConvAIConversationsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_convai_conversations"
}

func (d *ConvAIConversationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches ElevenLabs ConvAI conversations.",
		Attributes: map[string]schema.Attribute{
			"conversations": schema.ListAttribute{
				Computed:            true,
				ElementType:         convaiConversationObjectType,
				MarkdownDescription: "List of ConvAI conversations.",
			},
		},
	}
}

func (d *ConvAIConversationsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ConvAIConversationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data convaiConversationsDataSourceModel

	conversations, err := d.client.GetConvAIConversations()
	if err != nil {
		resp.Diagnostics.AddError("Error fetching ConvAI conversations", err.Error())
		return
	}

	convsList, diags := flattenConvAIConversations(ctx, conversations)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Conversations = convsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenConvAIConversations(ctx context.Context, conversations []map[string]interface{}) (types.List, diag.Diagnostics) {
	if len(conversations) == 0 {
		return types.ListNull(convaiConversationObjectType), nil
	}

	values := make([]attr.Value, 0, len(conversations))
	for _, conv := range conversations {
		obj, objDiags := types.ObjectValue(convaiConversationAttrTypes, map[string]attr.Value{
			"conversation_id": types.StringValue(conv["conversation_id"].(string)),
			"name":            types.StringValue(conv["name"].(string)),
			"agent_id":        types.StringValue(conv["agent_id"].(string)),
			"agent_name":      types.StringValue(conv["agent_name"].(string)),
			"created_at":      types.StringValue(conv["created_at"].(string)),
		})
		if objDiags.HasError() {
			return types.ListNull(convaiConversationObjectType), objDiags
		}

		values = append(values, obj)
	}

	return types.ListValue(convaiConversationObjectType, values)
}
