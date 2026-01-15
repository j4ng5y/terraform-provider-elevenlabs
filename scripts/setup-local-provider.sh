#!/bin/bash

# Build the provider
go build -o terraform-provider-elevenlabs

# Create local registry mirrors for testing (Terraform + OpenTofu)
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/j4ng5y/elevenlabs/0.1.0/linux_amd64
mkdir -p ~/.terraform.d/plugins/registry.opentofu.org/j4ng5y/elevenlabs/0.1.0/linux_amd64

# Copy the built provider to the local registry locations
cp terraform-provider-elevenlabs ~/.terraform.d/plugins/registry.terraform.io/j4ng5y/elevenlabs/0.1.0/linux_amd64/
cp terraform-provider-elevenlabs ~/.terraform.d/plugins/registry.opentofu.org/j4ng5y/elevenlabs/0.1.0/linux_amd64/

# Make it executable
chmod +x ~/.terraform.d/plugins/registry.terraform.io/j4ng5y/elevenlabs/0.1.0/linux_amd64/terraform-provider-elevenlabs
chmod +x ~/.terraform.d/plugins/registry.opentofu.org/j4ng5y/elevenlabs/0.1.0/linux_amd64/terraform-provider-elevenlabs
