package utils

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// ConvertList converts a Go slice of any type that implements attr.Value into a basetypes.ListValue.
// - If the slice is empty, returns an empty Terraform list (never null).
// - Returns Diagnostics from NewListValueFrom if any error occurs during conversion.
func ConvertList[E any](ctx context.Context, elementType attr.Type, list []E) (basetypes.ListValue, diag.Diagnostics) {
	if len(list) == 0 {
		return basetypes.NewListValueFrom(ctx, elementType, []attr.Value{})
	}
	return basetypes.NewListValueFrom(ctx, elementType, list)
}
