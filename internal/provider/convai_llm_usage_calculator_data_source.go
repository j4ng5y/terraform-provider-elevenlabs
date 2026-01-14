package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/j4ng5y/terraform-provider-elevenlabs/internal/client"
)

var (
	_ datasource.DataSource              = &ConvAILLMUsageCalculatorDataSource{}
	_ datasource.DataSourceWithConfigure = &ConvAILLMUsageCalculatorDataSource{}
)

func NewConvAILLMUsageCalculatorDataSource() datasource.DataSource {
	return &ConvAILLMUsageCalculatorDataSource{}
}

type ConvAILLMUsageCalculatorDataSource struct {
	client *client.Client
}

type ConvAILLMUsageCalculatorDataSourceModel struct {
	AgentID       types.String `tfsdk:"agent_id"`
	PromptLength  types.Int64  `tfsdk:"prompt_length"`
	NumberOfPages types.Int64  `tfsdk:"number_of_pages"`
	RAGEnabled    types.Bool   `tfsdk:"rag_enabled"`
	LLMPrices     types.List   `tfsdk:"llm_prices"`
}

type LLMPriceModel struct {
	ModelName    types.String  `tfsdk:"model_name"`
	InputCost    types.Float64 `tfsdk:"input_cost_per_million_tokens"`
	OutputCost   types.Float64 `tfsdk:"output_cost_per_million_tokens"`
	InputTokens  types.Int64   `tfsdk:"estimated_input_tokens"`
	OutputTokens types.Int64   `tfsdk:"estimated_output_tokens"`
	TotalCost    types.Float64 `tfsdk:"estimated_cost"`
}

func (d *ConvAILLMUsageCalculatorDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_convai_llm_usage_calculator"
}

func (d *ConvAILLMUsageCalculatorDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for calculating expected LLM usage and costs for ElevenLabs ConvAI agents.",
		Attributes: map[string]schema.Attribute{
			"agent_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the agent to calculate LLM usage for.",
			},
			"prompt_length": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Length of the prompt in characters.",
			},
			"number_of_pages": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Pages of content in PDF documents or URLs in agent's Knowledge Base.",
			},
			"rag_enabled": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether RAG is enabled.",
			},
			"llm_prices": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"model_name": schema.StringAttribute{
							Computed: true,
						},
						"input_cost_per_million_tokens": schema.Float64Attribute{
							Computed: true,
						},
						"output_cost_per_million_tokens": schema.Float64Attribute{
							Computed: true,
						},
						"estimated_input_tokens": schema.Int64Attribute{
							Computed: true,
						},
						"estimated_output_tokens": schema.Int64Attribute{
							Computed: true,
						},
						"estimated_cost": schema.Float64Attribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *ConvAILLMUsageCalculatorDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ConvAILLMUsageCalculatorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ConvAILLMUsageCalculatorDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	promptLength := 0
	if !data.PromptLength.IsNull() {
		promptLength = int(data.PromptLength.ValueInt64())
	}

	numberOfPages := 0
	if !data.NumberOfPages.IsNull() {
		numberOfPages = int(data.NumberOfPages.ValueInt64())
	}

	ragEnabled := false
	if !data.RAGEnabled.IsNull() {
		ragEnabled = data.RAGEnabled.ValueBool()
	}

	result, err := d.client.CalculateLLMUsage(
		data.AgentID.ValueString(),
		promptLength,
		numberOfPages,
		ragEnabled,
	)

	if err != nil {
		resp.Diagnostics.AddError("Error calculating LLM usage", err.Error())
		return
	}

	// Process LLM prices
	if llmPrices, ok := result["llm_prices"].([]interface{}); ok {
		prices := make([]LLMPriceModel, len(llmPrices))
		for i, price := range llmPrices {
			p := price.(map[string]interface{})
			prices[i] = LLMPriceModel{
				ModelName:    types.StringValue(p["model_name"].(string)),
				InputCost:    types.Float64Value(p["input_cost_per_million_tokens"].(float64)),
				OutputCost:   types.Float64Value(p["output_cost_per_million_tokens"].(float64)),
				InputTokens:  types.Int64Value(int64(p["estimated_input_tokens"].(float64))),
				OutputTokens: types.Int64Value(int64(p["estimated_output_tokens"].(float64))),
				TotalCost:    types.Float64Value(p["estimated_cost"].(float64)),
			}
		}

		pricesValue, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
			"model_name":                     types.StringType,
			"input_cost_per_million_tokens":  types.Float64Type,
			"output_cost_per_million_tokens": types.Float64Type,
			"estimated_input_tokens":         types.Int64Type,
			"estimated_output_tokens":        types.Int64Type,
			"estimated_cost":                 types.Float64Type,
		}}, prices)
		resp.Diagnostics.Append(diags...)
		if !diags.HasError() {
			data.LLMPrices = pricesValue
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
