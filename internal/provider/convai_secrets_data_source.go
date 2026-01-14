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

var convaiSecretAttrTypes = map[string]attr.Type{
	"secret_id": types.StringType,
	"name":      types.StringType,
}

var convaiSecretObjectType = types.ObjectType{AttrTypes: convaiSecretAttrTypes}

func NewConvAISecretsDataSource() datasource.DataSource {
	return &ConvAISecretsDataSource{}
}

type ConvAISecretsDataSource struct {
	client *client.Client
}

type convaiSecretsDataSourceModel struct {
	Secrets types.List `tfsdk:"secrets"`
}

func (d *ConvAISecretsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_convai_secrets"
}

func (d *ConvAISecretsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches ElevenLabs ConvAI secrets.",
		Attributes: map[string]schema.Attribute{
			"secrets": schema.ListAttribute{
				Computed:            true,
				ElementType:         convaiSecretObjectType,
				MarkdownDescription: "List of ConvAI secrets configured in the workspace.",
			},
		},
	}
}

func (d *ConvAISecretsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ConvAISecretsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data convaiSecretsDataSourceModel

	secrets, err := d.client.GetConvAISecrets()
	if err != nil {
		resp.Diagnostics.AddError("Error fetching ConvAI secrets", err.Error())
		return
	}

	secretsList, diags := flattenConvAISecrets(ctx, secrets)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Secrets = secretsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenConvAISecrets(ctx context.Context, secrets []models.ConvAISecret) (types.List, diag.Diagnostics) {
	if len(secrets) == 0 {
		return types.ListNull(convaiSecretObjectType), nil
	}

	values := make([]attr.Value, 0, len(secrets))
	for _, secret := range secrets {
		obj, objDiags := types.ObjectValue(convaiSecretAttrTypes, map[string]attr.Value{
			"secret_id": types.StringValue(secret.SecretID),
			"name":      types.StringValue(secret.Name),
		})
		if objDiags.HasError() {
			return types.ListNull(convaiSecretObjectType), objDiags
		}

		values = append(values, obj)
	}

	return types.ListValue(convaiSecretObjectType, values)
}
