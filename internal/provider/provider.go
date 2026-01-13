package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/j4ng5y/terraform-provider-elevenlabs/internal/client"
)

// ElevenLabsProvider defines the provider implementation.
type ElevenLabsProvider struct {
	version string
}

// ElevenLabsProviderModel describes the provider data model.
type ElevenLabsProviderModel struct {
	ApiKey types.String `tfsdk:"api_key"`
}

func (p *ElevenLabsProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "elevenlabs"
	resp.Version = p.version
}

func (p *ElevenLabsProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "ElevenLabs API Key. May also be provided via ELEVENLABS_API_KEY environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *ElevenLabsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ElevenLabsProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiKey := os.Getenv("ELEVENLABS_API_KEY")

	if !data.ApiKey.IsNull() {
		apiKey = data.ApiKey.ValueString()
	}

	if apiKey == "" {
		resp.Diagnostics.AddError(
			"Missing ElevenLabs API Key",
			"The provider cannot create the ElevenLabs API client as there is no API key. "+
				"Please set the api_key attribute in the provider block or use the ELEVENLABS_API_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	c := client.NewClient(apiKey)

	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *ElevenLabsProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewVoiceResource,
		NewProjectResource,
		NewPronunciationDictionaryResource,
		NewAudioNativeResource,
		NewDubbingProjectResource,
		NewConvAIAgentResource,
		NewConvAIKnowledgeBaseResource,
		NewConvAIToolResource,
		NewConvAISecretResource,
		NewWorkspaceWebhookResource,
		NewWorkspaceInviteResource,
		NewServiceAccountKeyResource,
		NewVoiceSampleResource,
		NewStudioChapterResource,
		NewConvAIMCPServerResource,
		NewConvAIPhoneNumberResource,
		NewWorkspaceMemberResource,
		NewWorkspaceGroupMembershipResource,
		NewConvAISettingsResource,
		NewResourceShareResource,
		NewConvAIAgentTestResource,
		NewSharedVoiceResource,
	}
}

func (p *ElevenLabsProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewModelsDataSource,
		NewVoicesDataSource,
		NewProjectsDataSource,
		NewPronunciationDictionariesDataSource,
		NewConvAIAgentsDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ElevenLabsProvider{
			version: version,
		}
	}
}
