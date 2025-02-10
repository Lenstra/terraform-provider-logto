package utils

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ConvertListToSlice(ctx context.Context, list types.List) []string {
	if list.IsNull() || list.IsUnknown() {
		return []string{}
	}
	var result []string
	list.ElementsAs(ctx, &result, false)
	return result
}

func StringSliceToList(slice []string) types.List {
	values := make([]attr.Value, len(slice))
	for i, s := range slice {
		values[i] = types.StringValue(s)
	}
	return types.ListValueMust(types.StringType, values)
}
