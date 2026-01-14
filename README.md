# Terraform Provider for ElevenLabs

[![CI](https://github.com/j4ng5y/terraform-provider-elevenlabs/actions/workflows/ci.yml/badge.svg)](https://github.com/j4ng5y/terraform-provider-elevenlabs/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/j4ng5y/terraform-provider-elevenlabs)](https://goreportcard.com/report/github.com/j4ng5y/terraform-provider-elevenlabs)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

The ElevenLabs Terraform provider allows you to manage your ElevenLabs resources (Voices, Projects, and more) as infrastructure. This provider is built using the modern Terraform Plugin Framework and is compatible with both Terraform and OpenTofu.

## Features

- **Voice Management**: Create, edit, and delete custom cloned voices with multipart file support.
- **ElevenLabs Studio Projects**: Manage long-form content projects programmatically.
- **Model Discovery**: Data source to fetch available models and their capabilities.
- **Modern Architecture**: Built with Terraform Plugin Framework for better performance and type safety.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21 (to build the provider plugin)
- ElevenLabs API Key (Grab one from your [Profile Settings](https://elevenlabs.io/app/subscription))

## Using the Provider

### Authentication

The provider requires an ElevenLabs API key. You can provide it in two ways:

1.  **Environment Variable** (Recommended): Set the `ELEVENLABS_API_KEY` environment variable.
    ```bash
    export ELEVENLABS_API_KEY="your-api-key"
    ```

2.  **Provider Block**: Set the `api_key` attribute in the provider block.
    ```hcl
    provider "elevenlabs" {
      api_key = "your-api-key"
    }
    ```

### Example Usage

```hcl
terraform {
  required_providers {
    elevenlabs = {
      source = "j4ng5y/elevenlabs"
      # For local development or OpenTofu, you can use local mirrors or specific registry paths
    }
  }
}

provider "elevenlabs" {}

# Get the latest Multilingual V2 model
data "elevenlabs_models" "all" {}

locals {
  # Find the specific model ID dynamically
  multilingual_v2 = [for m in data.elevenlabs_models.all.models : m.model_id if m.model_id == "eleven_multilingual_v2"][0]
}

# Create a cloned voice for a specific brand
resource "elevenlabs_voice" "brand_voice" {
  name        = "Global Brand Voice"
  description = "Primary voice for international marketing campaigns."
  files       = ["./samples/reference_audio.mp3"]
  labels = {
    "accent" = "neutral"
    "use_case" = "narrative"
  }
}

# Initialize a Studio project using the brand voice
resource "elevenlabs_project" "quarterly_update" {
  name                    = "Q4 Executive Update"
  default_model_id        = local.multilingual_v2
  default_paragraph_voice_id = elevenlabs_voice.brand_voice.id
}
```

## Documentation

Full documentation for resources and data sources can be found on the [Terraform Registry](https://registry.terraform.io/providers/j4ng5y/elevenlabs/latest/docs).

### Resources
- [elevenlabs_voice](./docs/resources/voice.md)
- [elevenlabs_project](./docs/resources/project.md)

### Data Sources
- [elevenlabs_models](./docs/data-sources/models.md)

## Development

### Building the Provider

1.  Clone the repository:
    ```bash
    git clone https://github.com/j4ng5y/terraform-provider-elevenlabs
    ```
2.  Enter the repository directory:
    ```bash
    cd terraform-provider-elevenlabs
    ```
3.  Build the provider:
    ```bash
    go build -o terraform-provider-elevenlabs
    ```

### Local Development

To use a locally built provider, you can use a [Terraform CLI configuration file](https://www.terraform.io/docs/cli/config/config-file.html#development-overrides-for-provider-developers) (`~/.terraformrc` or `%APPDATA%/terraform.rc`):

```hcl
provider_installation {
  dev_overrides {
    "j4ng5y/elevenlabs" = "<PATH_TO_YOUR_BUILD_DIRECTORY>"
  }
  direct {}
}
```

### Running Tests

To run the unit tests:
```bash
go test ./...
```

To run acceptance tests (requires `ELEVENLABS_API_KEY`):
```bash
export TF_ACC=1
export ELEVENLABS_API_KEY="..."
go test -v ./internal/provider/
```

### API Coverage Tracking

Use the OpenAPI coverage command to understand how much of the ElevenLabs surface area currently has Terraform parity:

```bash
go run ./cmd/coverage -exclude-tags "text-to-voice,text-to-speech,text-to-dialogue,voice-generation,speech-to-text,speech-to-speech,music-generation,audio-isolation,sound-generation,usage,speech-history,Single Use Token,forced-alignment,dubbing,studio"
```

Handy flags:

- `-include-tags "Agents Platform"` — focus on a single product area such as ConvAI.
- `-methods GET,POST` — constrain the report to specific HTTP verbs.
- `-details=false` — collapse the per-endpoint list when you only need counts.

The current backlog is dominated by the following CRUD-heavy domains (counts pulled from `go run ./cmd/coverage` on 2026-01-13):

1. **Agents Platform (~48 missing ops)** – conversations inventory, knowledge-base RAG indexes, MCP server configurations, agent tool approvals, and action endpoints (audio playback, outbound calls, test conversations).
2. **Workspace & Enterprise (13 combined)** – SSO groups, entitlements, and enterprise organization management.
3. **PVC Voices (~13 ops)** – voice model versions and fine-tuning operations.

Target Terraform types to add next (Agents Platform focus):

- Data sources: `elevenlabs_convai_conversations`, `elevenlabs_convai_knowledge_bases`, `elevenlabs_convai_tools`, `elevenlabs_convai_phone_numbers`, `elevenlabs_convai_whatsapp_accounts`, `elevenlabs_convai_mcp_servers`.
- Resources: `elevenlabs_convai_whatsapp_account`, `elevenlabs_convai_mcp_tool_config`, `elevenlabs_convai_mcp_tool_approval`, `elevenlabs_convai_knowledge_base_rag_index`, `elevenlabs_convai_conversation`.

These priorities cover every persistent CRUD API that is still missing from the provider. The remaining action-only endpoints (text-to-speech rendering, outbound call simulators, etc.) stay outside Terraform scope per the current design agreement.

## Roadmap

- [ ] Support for **Pronunciation Dictionaries**
- [ ] Support for **Audio Native** player settings
- [x] Dedicated **Voice Sample** resource for incremental updates
- [ ] **ConvAI Agents Platform** - Complete coverage (conversations, knowledge bases, tools, MCP servers)
- [x] **Workspace and Enterprise** management (members, groups, invites, service accounts)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

