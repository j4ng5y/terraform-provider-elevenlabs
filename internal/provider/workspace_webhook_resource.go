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
	_ resource.Resource                = &WorkspaceWebhookResource{}
	_ resource.ResourceWithConfigure   = &WorkspaceWebhookResource{}
	_ resource.ResourceWithImportState = &WorkspaceWebhookResource{}
)

func NewWorkspaceWebhookResource() resource.Resource {
	return &WorkspaceWebhookResource{}
}

type WorkspaceWebhookResource struct {
	client *client.Client
}

type WorkspaceWebhookResourceModel struct {
	ID     types.String `tfsdk:"id"`
	URL    types.String `tfsdk:"url"`
	Events types.List   `tfsdk:"events"`
	Secret types.String `tfsdk:"secret"`
}

func (r *WorkspaceWebhookResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace_webhook"
}

func (r *WorkspaceWebhookResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Workspace Webhook resource for ElevenLabs.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"url": schema.StringAttribute{
				Required: true,
			},
			"events": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
			"secret": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func (r *WorkspaceWebhookResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *WorkspaceWebhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data WorkspaceWebhookResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var events []string
	resp.Diagnostics.Append(data.Events.ElementsAs(ctx, &events, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	addReq := &models.CreateWorkspaceWebhookRequest{
		URL:    data.URL.ValueString(),
		Events: events,
	}

	webhook, err := r.client.CreateWorkspaceWebhook(addReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating workspace webhook", err.Error())
		return
	}

	data.ID = types.StringValue(webhook.WebhookID)
	data.Secret = types.StringValue(webhook.Secret)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkspaceWebhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data WorkspaceWebhookResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	webhook, err := r.client.GetWorkspaceWebhook(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading workspace webhook", err.Error())
		return
	}

	data.URL = types.StringValue(webhook.URL)

	events, diag := types.ListValueFrom(ctx, types.StringType, webhook.Events)
	resp.Diagnostics.Append(diag...)
	data.Events = events

	if webhook.Secret != "" {
		data.Secret = types.StringValue(webhook.Secret)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkspaceWebhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data WorkspaceWebhookResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var events []string
	resp.Diagnostics.Append(data.Events.ElementsAs(ctx, &events, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &models.CreateWorkspaceWebhookRequest{
		URL:    data.URL.ValueString(),
		Events: events,
	}

	err := r.client.UpdateWorkspaceWebhook(data.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating workspace webhook", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkspaceWebhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data WorkspaceWebhookResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteWorkspaceWebhook(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting workspace webhook", err.Error())
		return
	}
}

func (r *WorkspaceWebhookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
