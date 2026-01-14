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

var workspaceWebhookAttrTypes = map[string]attr.Type{
	"webhook_id": types.StringType,
	"url":        types.StringType,
	"events":     types.ListType{ElemType: types.StringType},
	"secret":     types.StringType,
}

var workspaceWebhookObjectType = types.ObjectType{AttrTypes: workspaceWebhookAttrTypes}

func NewWorkspaceWebhooksDataSource() datasource.DataSource {
	return &WorkspaceWebhooksDataSource{}
}

type WorkspaceWebhooksDataSource struct {
	client *client.Client
}

type workspaceWebhooksDataSourceModel struct {
	Webhooks types.List `tfsdk:"webhooks"`
}

func (d *WorkspaceWebhooksDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace_webhooks"
}

func (d *WorkspaceWebhooksDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches workspace webhooks.",
		Attributes: map[string]schema.Attribute{
			"webhooks": schema.ListAttribute{
				Computed:            true,
				ElementType:         workspaceWebhookObjectType,
				MarkdownDescription: "List of workspace webhooks.",
			},
		},
	}
}

func (d *WorkspaceWebhooksDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	d.client = c
}

func (d *WorkspaceWebhooksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data workspaceWebhooksDataSourceModel

	webhooks, err := d.client.ListWorkspaceWebhooks()
	if err != nil {
		resp.Diagnostics.AddError("Error fetching workspace webhooks", err.Error())
		return
	}

	webhooksList, diags := flattenWorkspaceWebhooks(ctx, webhooks)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Webhooks = webhooksList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenWorkspaceWebhooks(ctx context.Context, webhooks []models.WorkspaceWebhook) (types.List, diag.Diagnostics) {
	if len(webhooks) == 0 {
		return types.ListNull(workspaceWebhookObjectType), nil
	}

	values := make([]attr.Value, 0, len(webhooks))
	for _, webhook := range webhooks {
		events, diags := stringsToListValue(ctx, webhook.Events)
		if diags.HasError() {
			return types.ListNull(workspaceWebhookObjectType), diags
		}

		obj, objDiags := types.ObjectValue(workspaceWebhookAttrTypes, map[string]attr.Value{
			"webhook_id": types.StringValue(webhook.WebhookID),
			"url":        types.StringValue(webhook.URL),
			"events":     events,
			"secret":     optionalStringValue(webhook.Secret),
		})
		if objDiags.HasError() {
			return types.ListNull(workspaceWebhookObjectType), objDiags
		}

		values = append(values, obj)
	}

	return types.ListValue(workspaceWebhookObjectType, values)
}
