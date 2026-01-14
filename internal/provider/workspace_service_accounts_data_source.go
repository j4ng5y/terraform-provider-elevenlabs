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

var workspaceServiceAccountAttrTypes = map[string]attr.Type{
	"service_account_user_id": types.StringType,
	"name":                    types.StringType,
}

var workspaceServiceAccountObjectType = types.ObjectType{AttrTypes: workspaceServiceAccountAttrTypes}

func NewWorkspaceServiceAccountsDataSource() datasource.DataSource {
	return &WorkspaceServiceAccountsDataSource{}
}

type WorkspaceServiceAccountsDataSource struct {
	client *client.Client
}

type workspaceServiceAccountsDataSourceModel struct {
	ServiceAccounts types.List `tfsdk:"service_accounts"`
}

func (d *WorkspaceServiceAccountsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace_service_accounts"
}

func (d *WorkspaceServiceAccountsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches workspace service accounts.",
		Attributes: map[string]schema.Attribute{
			"service_accounts": schema.ListAttribute{
				Computed:            true,
				ElementType:         workspaceServiceAccountObjectType,
				MarkdownDescription: "List of workspace service accounts.",
			},
		},
	}
}

func (d *WorkspaceServiceAccountsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *WorkspaceServiceAccountsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data workspaceServiceAccountsDataSourceModel

	accounts, err := d.client.GetWorkspaceServiceAccounts()
	if err != nil {
		resp.Diagnostics.AddError("Error fetching workspace service accounts", err.Error())
		return
	}

	accountsList, diags := flattenWorkspaceServiceAccounts(ctx, accounts)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ServiceAccounts = accountsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenWorkspaceServiceAccounts(ctx context.Context, accounts []models.WorkspaceServiceAccount) (types.List, diag.Diagnostics) {
	if len(accounts) == 0 {
		return types.ListNull(workspaceServiceAccountObjectType), nil
	}

	values := make([]attr.Value, 0, len(accounts))
	for _, account := range accounts {
		obj, objDiags := types.ObjectValue(workspaceServiceAccountAttrTypes, map[string]attr.Value{
			"service_account_user_id": types.StringValue(account.ServiceAccountUserID),
			"name":                    types.StringValue(account.Name),
		})
		if objDiags.HasError() {
			return types.ListNull(workspaceServiceAccountObjectType), objDiags
		}

		values = append(values, obj)
	}

	return types.ListValue(workspaceServiceAccountObjectType, values)
}
