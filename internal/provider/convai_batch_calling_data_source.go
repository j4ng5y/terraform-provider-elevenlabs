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

var convaiBatchCallAttrTypes = map[string]attr.Type{
	"batch_id":        types.StringType,
	"agent_id":        types.StringType,
	"phone_number_id": types.StringType,
	"status":          types.StringType,
	"total_calls":     types.Int64Type,
	"completed_calls": types.Int64Type,
	"failed_calls":    types.Int64Type,
	"created_at":      types.StringType,
}

var convaiBatchCallObjectType = types.ObjectType{AttrTypes: convaiBatchCallAttrTypes}

func NewConvAIBatchCallingDataSource() datasource.DataSource {
	return &ConvAIBatchCallingDataSource{}
}

type ConvAIBatchCallingDataSource struct {
	client *client.Client
}

type convaiBatchCallingDataSourceModel struct {
	BatchCalls types.List `tfsdk:"batch_calls"`
}

func (d *ConvAIBatchCallingDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_convai_batch_calling"
}

func (d *ConvAIBatchCallingDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches ElevenLabs ConvAI batch calling requests.",
		Attributes: map[string]schema.Attribute{
			"batch_calls": schema.ListAttribute{
				Computed:            true,
				ElementType:         convaiBatchCallObjectType,
				MarkdownDescription: "List of ConvAI batch calling requests.",
			},
		},
	}
}

func (d *ConvAIBatchCallingDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ConvAIBatchCallingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data convaiBatchCallingDataSourceModel

	batches, err := d.client.GetConvAIBatchCalls()
	if err != nil {
		resp.Diagnostics.AddError("Error fetching batch calls", err.Error())
		return
	}

	batchesList, diags := flattenConvAIBatchCalls(ctx, batches)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.BatchCalls = batchesList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenConvAIBatchCalls(ctx context.Context, batches []map[string]interface{}) (types.List, diag.Diagnostics) {
	if len(batches) == 0 {
		return types.ListNull(convaiBatchCallObjectType), nil
	}

	values := make([]attr.Value, 0, len(batches))
	for _, batch := range batches {
		obj, objDiags := types.ObjectValue(convaiBatchCallAttrTypes, map[string]attr.Value{
			"batch_id":        types.StringValue(batch["batch_id"].(string)),
			"agent_id":        types.StringValue(batch["agent_id"].(string)),
			"phone_number_id": types.StringValue(batch["phone_number_id"].(string)),
			"status":          types.StringValue(batch["status"].(string)),
			"total_calls":     types.Int64Value(int64(batch["total_calls"].(float64))),
			"completed_calls": types.Int64Value(int64(batch["completed_calls"].(float64))),
			"failed_calls":    types.Int64Value(int64(batch["failed_calls"].(float64))),
			"created_at":      types.StringValue(batch["created_at"].(string)),
		})
		if objDiags.HasError() {
			return types.ListNull(convaiBatchCallObjectType), objDiags
		}

		values = append(values, obj)
	}

	return types.ListValue(convaiBatchCallObjectType, values)
}
