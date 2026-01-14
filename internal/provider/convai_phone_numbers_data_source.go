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
	_ datasource.DataSource              = &ConvAIPhoneNumbersDataSource{}
	_ datasource.DataSourceWithConfigure = &ConvAIPhoneNumbersDataSource{}
)

var convAIPhoneNumberAttrTypes = map[string]attr.Type{
	"phone_number_id":      types.StringType,
	"phone_number":         types.StringType,
	"provider":             types.StringType,
	"label":                types.StringType,
	"supports_inbound":     types.BoolType,
	"supports_outbound":    types.BoolType,
	"assigned_agent_id":    types.StringType,
	"assigned_agent_name":  types.StringType,
	"provider_config_json": types.StringType,
	"outbound_trunk_json":  types.StringType,
	"inbound_trunk_json":   types.StringType,
	"livekit_stack_json":   types.StringType,
}

var convAIPhoneNumberObjectType = types.ObjectType{AttrTypes: convAIPhoneNumberAttrTypes}

func NewConvAIPhoneNumbersDataSource() datasource.DataSource {
	return &ConvAIPhoneNumbersDataSource{}
}

type ConvAIPhoneNumbersDataSource struct {
	client *client.Client
}

type convaiPhoneNumbersDataSourceModel struct {
	PhoneNumbers types.List `tfsdk:"phone_numbers"`
}

func (d *ConvAIPhoneNumbersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_convai_phone_numbers"
}

func (d *ConvAIPhoneNumbersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches ElevenLabs ConvAI phone numbers.",
		Attributes: map[string]schema.Attribute{
			"phone_numbers": schema.ListAttribute{
				Computed:            true,
				ElementType:         convAIPhoneNumberObjectType,
				MarkdownDescription: "List of ConvAI phone numbers configured in the workspace.",
			},
		},
	}
}

func (d *ConvAIPhoneNumbersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ConvAIPhoneNumbersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data convaiPhoneNumbersDataSourceModel

	numbers, err := d.client.GetConvAIPhoneNumbers()
	if err != nil {
		resp.Diagnostics.AddError("Error fetching ConvAI phone numbers", err.Error())
		return
	}

	numbersList, diags := flattenConvAIPhoneNumbers(ctx, numbers)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.PhoneNumbers = numbersList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenConvAIPhoneNumbers(ctx context.Context, numbers []models.ConvAIPhoneNumber) (types.List, diag.Diagnostics) {
	if len(numbers) == 0 {
		return types.ListNull(convAIPhoneNumberObjectType), nil
	}

	values := make([]attr.Value, 0, len(numbers))
	for _, number := range numbers {
		obj, objDiags := types.ObjectValue(convAIPhoneNumberAttrTypes, map[string]attr.Value{
			"phone_number_id":      types.StringValue(number.PhoneNumberID),
			"phone_number":         types.StringValue(number.PhoneNumber),
			"provider":             types.StringValue(number.Provider),
			"label":                optionalStringValue(number.Label),
			"supports_inbound":     types.BoolValue(number.SupportsInbound),
			"supports_outbound":    types.BoolValue(number.SupportsOutbound),
			"assigned_agent_id":    optionalStringValue(number.AssignedAgent.AgentID),
			"assigned_agent_name":  optionalStringValue(number.AssignedAgent.AgentName),
			"provider_config_json": jsonStringValue(number.ProviderConfig),
			"outbound_trunk_json":  jsonStringValue(number.OutboundTrunk),
			"inbound_trunk_json":   jsonStringValue(number.InboundTrunk),
			"livekit_stack_json":   jsonStringValue(number.LivekitStack),
		})
		if objDiags.HasError() {
			return types.ListNull(convAIPhoneNumberObjectType), objDiags
		}

		values = append(values, obj)
	}

	return types.ListValue(convAIPhoneNumberObjectType, values)
}
