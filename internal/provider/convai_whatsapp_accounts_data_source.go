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
	_ datasource.DataSource              = &ConvAIWhatsAppAccountsDataSource{}
	_ datasource.DataSourceWithConfigure = &ConvAIWhatsAppAccountsDataSource{}
)

var convAIWhatsAppAccountAttrTypes = map[string]attr.Type{
	"business_account_id":   types.StringType,
	"business_account_name": types.StringType,
	"phone_number_id":       types.StringType,
	"phone_number_name":     types.StringType,
	"phone_number":          types.StringType,
	"assigned_agent_id":     types.StringType,
	"assigned_agent_name":   types.StringType,
}

var convAIWhatsAppAccountObjectType = types.ObjectType{AttrTypes: convAIWhatsAppAccountAttrTypes}

func NewConvAIWhatsAppAccountsDataSource() datasource.DataSource {
	return &ConvAIWhatsAppAccountsDataSource{}
}

type ConvAIWhatsAppAccountsDataSource struct {
	client *client.Client
}

type convaiWhatsAppAccountsDataSourceModel struct {
	WhatsAppAccounts types.List `tfsdk:"whatsapp_accounts"`
}

func (d *ConvAIWhatsAppAccountsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_convai_whatsapp_accounts"
}

func (d *ConvAIWhatsAppAccountsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches ElevenLabs ConvAI WhatsApp accounts.",
		Attributes: map[string]schema.Attribute{
			"whatsapp_accounts": schema.ListAttribute{
				Computed:            true,
				ElementType:         convAIWhatsAppAccountObjectType,
				MarkdownDescription: "List of ConvAI WhatsApp accounts configured in the workspace.",
			},
		},
	}
}

func (d *ConvAIWhatsAppAccountsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ConvAIWhatsAppAccountsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data convaiWhatsAppAccountsDataSourceModel

	accounts, err := d.client.ListConvAIWhatsAppAccounts()
	if err != nil {
		resp.Diagnostics.AddError("Error fetching ConvAI WhatsApp accounts", err.Error())
		return
	}

	accountsList, diags := flattenConvAIWhatsAppAccounts(ctx, accounts)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.WhatsAppAccounts = accountsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenConvAIWhatsAppAccounts(ctx context.Context, accounts []models.ConvAIWhatsAppAccount) (types.List, diag.Diagnostics) {
	if len(accounts) == 0 {
		return types.ListNull(convAIWhatsAppAccountObjectType), nil
	}

	values := make([]attr.Value, 0, len(accounts))
	for _, account := range accounts {
		obj, objDiags := types.ObjectValue(convAIWhatsAppAccountAttrTypes, map[string]attr.Value{
			"business_account_id":   types.StringValue(account.BusinessAccountID),
			"business_account_name": types.StringValue(account.BusinessAccountName),
			"phone_number_id":       types.StringValue(account.PhoneNumberID),
			"phone_number_name":     types.StringValue(account.PhoneNumberName),
			"phone_number":          types.StringValue(account.PhoneNumber),
			"assigned_agent_id":     optionalStringValue(account.AssignedAgentID),
			"assigned_agent_name":   optionalStringValue(account.AssignedAgentName),
		})
		if objDiags.HasError() {
			return types.ListNull(convAIWhatsAppAccountObjectType), objDiags
		}

		values = append(values, obj)
	}

	return types.ListValue(convAIWhatsAppAccountObjectType, values)
}
