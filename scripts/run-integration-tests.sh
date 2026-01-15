#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TEST_DIR="$ROOT_DIR/test"

if [[ -z "${TF_VAR_api_key:-}" ]]; then
  echo "TF_VAR_api_key is required to run integration tests." >&2
  exit 1
fi

TOFU_BIN="${TOFU_BIN:-tofu}"

"$ROOT_DIR/scripts/setup-local-provider.sh"

TF_CLI_CONFIG_FILE="$(mktemp)"
trap 'rm -f "$TF_CLI_CONFIG_FILE"' EXIT

cat > "$TF_CLI_CONFIG_FILE" <<EOF
provider_installation {
  filesystem_mirror {
    path    = "${HOME}/.terraform.d/plugins"
    include = ["j4ng5y/elevenlabs"]
  }
  direct {
    exclude = ["j4ng5y/elevenlabs"]
  }
}
EOF

TF_CLI_CONFIG_FILE="$TF_CLI_CONFIG_FILE" "$TOFU_BIN" -chdir="$TEST_DIR" init -input=false
TF_CLI_CONFIG_FILE="$TF_CLI_CONFIG_FILE" "$TOFU_BIN" -chdir="$TEST_DIR" plan -input=false
TF_CLI_CONFIG_FILE="$TF_CLI_CONFIG_FILE" "$TOFU_BIN" -chdir="$TEST_DIR" apply -auto-approve -input=false
TF_CLI_CONFIG_FILE="$TF_CLI_CONFIG_FILE" "$TOFU_BIN" -chdir="$TEST_DIR" destroy -auto-approve -input=false
