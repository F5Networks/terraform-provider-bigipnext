package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func requiredIf(attribute path.Path, value string) stringRequiredIfValidator {
	return stringRequiredIfValidator{
		attribute: attribute,
		value:     value,
	}
}

type stringRequiredIfValidator struct {
	attribute path.Path
	value     string
}

func (v stringRequiredIfValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("This field is required if %q is set to %q", v.attribute.String(), v.value)
}

func (v stringRequiredIfValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v stringRequiredIfValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	var attr types.String
	req.Config.GetAttribute(ctx, v.attribute, &attr)

	requiredAttr := req.Path.String()

	if attr.ValueString() == v.value {
		if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() || req.ConfigValue.ValueString() == "" {
			resp.Diagnostics.AddError( // coverage-ignore
				"validation error: ",
				fmt.Sprintf("Attribute %v is required, when %v is set to %v", requiredAttr, v.attribute, v.value),
			)
		}
	}
}

func (v stringRequiredIfValidator) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	var attr types.String
	req.Config.GetAttribute(ctx, v.attribute, &attr)

	requiredAttr := req.Path.String()

	if attr.ValueString() == v.value {
		if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() { // coverage-ignore
			resp.Diagnostics.AddError(
				"validation error: ",
				fmt.Sprintf("Attribute %v is required, when %v is set to %v", requiredAttr, v.attribute, v.value),
			)
		}
	}
}
