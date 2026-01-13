package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/j4ng5y/terraform-provider-elevenlabs/internal/client"
	"github.com/j4ng5y/terraform-provider-elevenlabs/internal/models"
)

var (
	_ resource.Resource                = &ConvAIAgentTestResource{}
	_ resource.ResourceWithConfigure   = &ConvAIAgentTestResource{}
	_ resource.ResourceWithImportState = &ConvAIAgentTestResource{}
)

func NewConvAIAgentTestResource() resource.Resource {
	return &ConvAIAgentTestResource{}
}

type ConvAIAgentTestResource struct {
	client *client.Client
}

type ConvAIAgentTestResourceModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	SuccessCondition types.String `tfsdk:"success_condition"`
}

func (r *ConvAIAgentTestResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_convai_agent_test"
}

func (r *ConvAIAgentTestResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Conversational AI Agent Test resource for ElevenLabs. Allows defining unit tests for agents.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"success_condition": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (r *ConvAIAgentTestResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ConvAIAgentTestResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ConvAIAgentTestResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	addReq := &models.CreateConvAIAgentTestRequest{
		Name:             data.Name.ValueString(),
		SuccessCondition: data.SuccessCondition.ValueString(),
		ChatHistory:      []interface{}{}, // Placeholder for simplicity
		SuccessExamples:  []string{},
		FailureExamples:  []string{},
	}

	test, err := r.client.CreateConvAIAgentTest(addReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating ConvAI agent test", err.Error())
		return
	}

	data.ID = types.StringValue(test.TestID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAIAgentTestResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ConvAIAgentTestResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	test, err := r.client.GetConvAIAgentTest(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading ConvAI agent test", err.Error())
		return
	}

	data.Name = types.StringValue(test.Name)
	data.SuccessCondition = types.StringValue(test.SuccessCondition)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConvAIAgentTestResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning("Update limited", "Updating ConvAI agent tests might require replacement.")
}

func (r *ConvAIAgentTestResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ConvAIAgentTestResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteConvAIAgentTest(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting ConvAI agent test", err.Error())
		return
	}
}

func (r *ConvAIAgentTestResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
