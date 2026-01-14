package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/j4ng5y/terraform-provider-elevenlabs/internal/client"
)

var (
	_ resource.Resource              = &ConvAIAgentTestRunnerResource{}
	_ resource.ResourceWithConfigure = &ConvAIAgentTestRunnerResource{}
)

func NewConvAIAgentTestRunnerResource() resource.Resource {
	return &ConvAIAgentTestRunnerResource{}
}

type ConvAIAgentTestRunnerResource struct {
	client *client.Client
}

type ConvAIAgentTestRunnerResourceModel struct {
	AgentID          types.String   `tfsdk:"agent_id"`
	TestIDs          []types.String `tfsdk:"test_ids"`
	AgentConfig      types.Map      `tfsdk:"agent_configuration"`
	TestInvocationID types.String   `tfsdk:"test_invocation_id"`
	Status           types.String   `tfsdk:"status"`
	Results          types.List     `tfsdk:"results"`
	StartedAt        types.String   `tfsdk:"started_at"`
	CompletedAt      types.String   `tfsdk:"completed_at"`
}

type TestResultModel struct {
	TestID   types.String `tfsdk:"test_id"`
	Status   types.String `tfsdk:"status"`
	Passed   types.Bool   `tfsdk:"passed"`
	Message  types.String `tfsdk:"message"`
	Duration types.Int64  `tfsdk:"duration_ms"`
}

func (r *ConvAIAgentTestRunnerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_convai_agent_test_runner"
}

func (r *ConvAIAgentTestRunnerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource for running ElevenLabs ConvAI agent tests. Executes test suites on agents with optional configuration overrides.",
		Attributes: map[string]schema.Attribute{
			"agent_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the agent to run tests on.",
			},
			"test_ids": schema.ListAttribute{
				Required:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of test IDs to execute.",
			},
			"agent_configuration": schema.MapAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Optional agent configuration to override during testing.",
			},
			"test_invocation_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the test invocation.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The status of the test run.",
			},
			"results": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"test_id": schema.StringAttribute{
							Computed: true,
						},
						"status": schema.StringAttribute{
							Computed: true,
						},
						"passed": schema.BoolAttribute{
							Computed: true,
						},
						"message": schema.StringAttribute{
							Computed: true,
						},
						"duration_ms": schema.Int64Attribute{
							Computed: true,
						},
					},
				},
			},
			"started_at": schema.StringAttribute{
				Computed: true,
			},
			"completed_at": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *ConvAIAgentTestRunnerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *ConvAIAgentTestRunnerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ConvAIAgentTestRunnerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	testIDs := make([]string, len(plan.TestIDs))
	for i, testID := range plan.TestIDs {
		testIDs[i] = testID.ValueString()
	}

	var agentConfig map[string]interface{}
	if !plan.AgentConfig.IsNull() {
		agentConfigElements := plan.AgentConfig.Elements()
		agentConfig = make(map[string]interface{})
		for key, value := range agentConfigElements {
			agentConfig[key] = value
		}
	}

	result, err := r.client.RunConvAIAgentTests(plan.AgentID.ValueString(), testIDs, agentConfig)
	if err != nil {
		resp.Diagnostics.AddError("Error running agent tests", err.Error())
		return
	}

	plan.TestInvocationID = types.StringValue(result["test_invocation_id"].(string))
	plan.Status = types.StringValue(result["status"].(string))

	if startedAt, ok := result["started_at"].(string); ok {
		plan.StartedAt = types.StringValue(startedAt)
	}
	if completedAt, ok := result["completed_at"].(string); ok {
		plan.CompletedAt = types.StringValue(completedAt)
	}

	// Process results
	if results, ok := result["results"].([]interface{}); ok {
		testResults := make([]TestResultModel, len(results))
		for i, res := range results {
			r := res.(map[string]interface{})
			testResults[i] = TestResultModel{
				TestID:   types.StringValue(r["test_id"].(string)),
				Status:   types.StringValue(r["status"].(string)),
				Passed:   types.BoolValue(r["passed"].(bool)),
				Message:  types.StringValue(r["message"].(string)),
				Duration: types.Int64Value(int64(r["duration_ms"].(float64))),
			}
		}

		resultsValue, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: map[string]attr.Type{
			"test_id":     types.StringType,
			"status":      types.StringType,
			"passed":      types.BoolType,
			"message":     types.StringType,
			"duration_ms": types.Int64Type,
		}}, testResults)
		resp.Diagnostics.Append(diags...)
		plan.Results = resultsValue
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ConvAIAgentTestRunnerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// This is an action resource, no persistent state needed
}

func (r *ConvAIAgentTestRunnerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ConvAIAgentTestRunnerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	testIDs := make([]string, len(plan.TestIDs))
	for i, testID := range plan.TestIDs {
		testIDs[i] = testID.ValueString()
	}

	var agentConfig map[string]interface{}
	if !plan.AgentConfig.IsNull() {
		agentConfigElements := plan.AgentConfig.Elements()
		agentConfig = make(map[string]interface{})
		for key, value := range agentConfigElements {
			agentConfig[key] = value
		}
	}

	result, err := r.client.RunConvAIAgentTests(plan.AgentID.ValueString(), testIDs, agentConfig)
	if err != nil {
		resp.Diagnostics.AddError("Error running agent tests", err.Error())
		return
	}

	plan.TestInvocationID = types.StringValue(result["test_invocation_id"].(string))
	plan.Status = types.StringValue(result["status"].(string))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ConvAIAgentTestRunnerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// This is an action resource, no deletion needed
}
