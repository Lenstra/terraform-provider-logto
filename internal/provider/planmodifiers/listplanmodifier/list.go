package listplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NullIsEmpty() planmodifier.List {
	return nullIsEmptyModifier{}
}

// nullIsEmptyModifier implements the plan modifier.
type nullIsEmptyModifier struct{}

// Description returns a human-readable description of the plan modifier.
func (m nullIsEmptyModifier) Description(_ context.Context) string {
	return "If the list is null, it will be set to an empty list."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m nullIsEmptyModifier) MarkdownDescription(_ context.Context) string {
	return "If the list is null, it will be set to an empty list."
}

// PlanModifyList implements the plan modification logic.
func (m nullIsEmptyModifier) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if req.ConfigValue.IsNull() {
		resp.PlanValue = types.ListValueMust(req.ConfigValue.ElementType(ctx), nil)
	}
}
