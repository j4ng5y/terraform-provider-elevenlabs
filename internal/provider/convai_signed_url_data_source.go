package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/j4ng5y/terraform-provider-elevenlabs/internal/client"
)

var (
	_ datasource.DataSource              = &ConvAISignedUrlDataSource{}
	_ datasource.DataSourceWithConfigure = &ConvAISignedUrlDataSource{}
)

func NewConvAISignedUrlDataSource() datasource.DataSource {
	return &ConvAISignedUrlDataSource{}
}

type ConvAISignedUrlDataSource struct {
	client *client.Client
}

type ConvAISignedUrlDataSourceModel struct {
	AgentID               types.String `tfsdk:"agent_id"`
	IncludeConversationID types.Bool   `tfsdk:"include_conversation_id"`
	ConversationSignature types.String `tfsdk:"conversation_signature"`
	ConversationID        types.String `tfsdk:"conversation_id"`
}

func (d *ConvAISignedUrlDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_convai_signed_url"
}

func (d *ConvAISignedUrlDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for generating signed URLs for ElevenLabs ConvAI agent conversations.",
		Attributes: map[string]schema.Attribute{
			"agent_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the agent to generate a signed URL for.",
			},
			"include_conversation_id": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether to include a conversation ID with the response. If included, the conversation_signature cannot be used again.",
			},
			"conversation_signature": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The signed conversation signature for authentication.",
			},
			"conversation_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The conversation ID (only included if include_conversation_id is true).",
			},
		},
	}
}

func (d *ConvAISignedUrlDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *ConvAISignedUrlDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ConvAISignedUrlDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	includeConversationID := false
	if !data.IncludeConversationID.IsNull() {
		includeConversationID = data.IncludeConversationID.ValueBool()
	}

	signature, conversationID, err := d.client.GetConvAISignedUrl(
		data.AgentID.ValueString(),
		includeConversationID,
	)

	if err != nil {
		resp.Diagnostics.AddError("Error generating signed URL", err.Error())
		return
	}

	data.ConversationSignature = types.StringValue(signature)
	if conversationID != "" {
		data.ConversationID = types.StringValue(conversationID)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
