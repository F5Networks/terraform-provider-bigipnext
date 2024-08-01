//go:build !test

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Int64DefaultValue is assign default values for Int Type
func Int64DefaultValue(v types.Int64) planmodifier.Int64 {
	return &int64DefaultValuePlanModifier{v}
}

type int64DefaultValuePlanModifier struct {
	DefaultValue types.Int64
}

var _ planmodifier.Int64 = (*int64DefaultValuePlanModifier)(nil)

func (apm *int64DefaultValuePlanModifier) Description(ctx context.Context) string {
	return apm.MarkdownDescription(ctx)
}

func (apm *int64DefaultValuePlanModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Sets the default value %q (%s) if the attribute is not set", apm.DefaultValue, apm.DefaultValue.Type(ctx))
}

func (apm *int64DefaultValuePlanModifier) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, res *planmodifier.Int64Response) {
	// If the attribute configuration is not null, we are done here
	if !req.ConfigValue.IsNull() {
		return
	}

	// If the attribute plan is "known" and "not null", then a previous plan modifier in the sequence
	// has already been applied, and we don't want to interfere.
	if !req.PlanValue.IsUnknown() && !req.PlanValue.IsNull() {
		return
	}
	res.PlanValue = apm.DefaultValue
}

// StringDefaultValue is assign default values for String Type
func StringDefaultValue(v types.String) planmodifier.String {
	return &stringDefaultValuePlanModifier{v}
}

type stringDefaultValuePlanModifier struct {
	DefaultValue types.String
}

var _ planmodifier.String = (*stringDefaultValuePlanModifier)(nil)

func (apm *stringDefaultValuePlanModifier) Description(ctx context.Context) string {
	return apm.MarkdownDescription(ctx)
}

func (apm *stringDefaultValuePlanModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Sets the default value %q (%s) if the attribute is not set", apm.DefaultValue, apm.DefaultValue.Type(ctx))
}

func (apm *stringDefaultValuePlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, res *planmodifier.StringResponse) {
	// If the attribute configuration is not null, we are done here
	if !req.ConfigValue.IsNull() {
		return
	}

	// If the attribute plan is "known" and "not null", then a previous plan modifier in the sequence
	// has already been applied, and we don't want to interfere.
	if !req.PlanValue.IsUnknown() && !req.PlanValue.IsNull() {
		return
	}
	res.PlanValue = apm.DefaultValue
}

// BoolDefaultValue is assign default values for Bool Type
func BoolDefaultValue(v types.Bool) planmodifier.Bool {
	return &boolDefaultValuePlanModifier{v}
}

type boolDefaultValuePlanModifier struct {
	DefaultValue types.Bool
}

var _ planmodifier.Bool = (*boolDefaultValuePlanModifier)(nil)

func (apm *boolDefaultValuePlanModifier) Description(ctx context.Context) string {
	return apm.MarkdownDescription(ctx)
}

func (apm *boolDefaultValuePlanModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Sets the default value %q (%s) if the attribute is not set", apm.DefaultValue, apm.DefaultValue.Type(ctx))
}

func (apm *boolDefaultValuePlanModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, res *planmodifier.BoolResponse) {
	// If the attribute configuration is not null, we are done here
	if !req.ConfigValue.IsNull() {
		return
	}

	// If the attribute plan is "known" and "not null", then a previous plan modifier in the sequence
	// has already been applied, and we don't want to interfere.
	if !req.PlanValue.IsUnknown() && !req.PlanValue.IsNull() {
		return
	}
	res.PlanValue = apm.DefaultValue
}

type as3JsonPlanModifier struct{}

func (m as3JsonPlanModifier) Description(ctx context.Context) string {
	return "plan modifier for as3 json."
}

func (m as3JsonPlanModifier) MarkdownDescription(_ context.Context) string {
	return "plan modifier for as3 json."
}

func (m as3JsonPlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, res *planmodifier.StringResponse) {
	if req.State.Raw.IsNull() {
		return
	}
	as3Statedata := []byte(req.StateValue.ValueString())
	as3Plandata := []byte(req.PlanValue.ValueString())
	var stateAs3 map[string]interface{}
	var planAs3 map[string]interface{}
	_ = json.Unmarshal(as3Statedata, &stateAs3)
	_ = json.Unmarshal(as3Plandata, &planAs3)
	if reflect.DeepEqual(stateAs3, planAs3) {
		res.PlanValue = req.StateValue
	}
}

func SuppressDiffAs3() planmodifier.String {
	return as3JsonPlanModifier{}
}
