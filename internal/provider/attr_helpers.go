package provider

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func stringSliceFromList(ctx context.Context, list types.List) ([]string, diag.Diagnostics) {
	var result []string
	if list.IsNull() || list.IsUnknown() {
		return result, nil
	}

	diags := list.ElementsAs(ctx, &result, false)
	return result, diags
}

func stringsToListValue(ctx context.Context, items []string) (types.List, diag.Diagnostics) {
	if len(items) == 0 {
		return types.ListNull(types.StringType), nil
	}

	return types.ListValueFrom(ctx, types.StringType, items)
}

func optionalStringValue(value string) types.String {
	if value == "" {
		return types.StringNull()
	}

	return types.StringValue(value)
}

func jsonStringValue(value interface{}) types.String {
	if value == nil {
		return types.StringNull()
	}

	data, err := json.Marshal(value)
	if err != nil || len(data) == 0 {
		return types.StringNull()
	}

	return types.StringValue(string(data))
}
