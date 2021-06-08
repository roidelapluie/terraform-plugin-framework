package reflect

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type testObjectType struct {
	AttributeTypes map[string]attr.Type
}

func (t testObjectType) TerraformType(ctx context.Context) tftypes.Type {
	types := map[string]tftypes.Type{}
	for k, v := range t.AttributeTypes {
		types[k] = v.TerraformType(ctx)
	}
	return tftypes.Object{
		AttributeTypes: types,
	}
}

func (t testObjectType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if !in.Type().Is(tftypes.Object{}) {
		return nil, fmt.Errorf("unexpected type %s", in.Type())
	}
	result := &testObjectValue{
		AttributeTypes: t.AttributeTypes,
	}
	if in.IsNull() {
		result.Null = true
		return result, nil
	}
	if !in.IsKnown() {
		result.Unknown = true
		return result, nil
	}
	inVals := map[string]tftypes.Value{}
	err := in.As(&inVals)
	if err != nil {
		return nil, err
	}
	if len(inVals) != len(t.AttributeTypes) {
		return nil, fmt.Errorf("expected value to have %d attributes, has %d: %s", len(t.AttributeTypes), len(inVals), in)
	}
	result.Attributes = map[string]attr.Value{}
	for k, v := range t.AttributeTypes {
		val, ok := inVals[k]
		if !ok {
			return nil, fmt.Errorf("expected value to have attribute %q, doesn't", k)
		}
		res, err := v.ValueFromTerraform(ctx, val)
		if err != nil {
			return nil, fmt.Errorf("error converting %q to attr.Value: %w", k, err)
		}
		result.Attributes[k] = res
	}
	return result, nil
}

func (t testObjectType) Equal(other attr.Type) bool {
	o, ok := other.(testObjectType)
	if !ok {
		return false
	}
	if len(o.AttributeTypes) != len(t.AttributeTypes) {
		return false
	}
	for k, v := range t.AttributeTypes {
		ov, ok := o.AttributeTypes[k]
		if !ok {
			return false
		}
		if !ov.Equal(v) {
			return false
		}
	}
	return true
}

var _ attr.Type = testObjectType{}

type testObjectValue struct {
	AttributeTypes map[string]attr.Type
	Attributes     map[string]attr.Value
	Unknown        bool
	Null           bool
}

var _ attr.Value = &testObjectValue{}

func (t *testObjectValue) ToTerraformValue(ctx context.Context) (interface{}, error) {
	if t.Unknown {
		return tftypes.UnknownValue, nil
	}
	if t.Null {
		return nil, nil
	}
	resultVals := map[string]tftypes.Value{}
	for k, v := range t.Attributes {
		typ, ok := t.AttributeTypes[k]
		if !ok {
			return nil, fmt.Errorf("no type for attribute %q", k)
		}
		val, err := v.ToTerraformValue(ctx)
		if err != nil {
			return nil, fmt.Errorf("error generating terraform value for %q: %w", k, err)
		}
		err = tftypes.ValidateValue(typ.TerraformType(ctx), val)
		if err != nil {
			return nil, fmt.Errorf("error validating terraform value for %q: %w", k, err)
		}
		resultVals[k] = tftypes.NewValue(typ.TerraformType(ctx), val)
	}
	return resultVals, nil
}

func (t *testObjectValue) SetTerraformValue(ctx context.Context, in tftypes.Value) error {
	t.Unknown = false
	t.Null = false
	t.Attributes = map[string]attr.Value{}
	if !in.Type().Is(tftypes.Object{}) {
		return fmt.Errorf("unexpected type %s", in.Type())
	}
	if !in.IsKnown() {
		t.Unknown = true
		return nil
	}
	if in.IsNull() {
		t.Null = true
		return nil
	}
	resultVals := map[string]tftypes.Value{}
	err := in.As(&resultVals)
	if err != nil {
		return err
	}
	if len(resultVals) != len(t.AttributeTypes) {
		return fmt.Errorf("expected %d attributes, got %d", len(t.AttributeTypes), len(resultVals))
	}
	attrs := map[string]attr.Value{}
	for k, v := range resultVals {
		typ, ok := t.AttributeTypes[k]
		if !ok {
			return fmt.Errorf("no type defined for %q", k)
		}
		a, err := typ.ValueFromTerraform(ctx, v)
		if err != nil {
			return fmt.Errorf("error building value for %q: %w", k, err)
		}
		attrs[k] = a
	}
	t.Attributes = attrs
	return nil
}

func (t *testObjectValue) Equal(other attr.Value) bool {
	if t == nil && other == nil {
		return true
	}
	if t == nil || other == nil {
		return false
	}
	o, ok := other.(*testObjectValue)
	if !ok {
		return false
	}
	if t.Null != o.Null {
		return false
	}
	if t.Unknown != o.Unknown {
		return false
	}
	if len(t.AttributeTypes) != len(o.AttributeTypes) {
		return false
	}
	for k, v := range t.AttributeTypes {
		ov, ok := o.AttributeTypes[k]
		if !ok {
			return false
		}
		if !v.Equal(ov) {
			return false
		}
	}
	if len(t.Attributes) != len(o.Attributes) {
		return false
	}
	for k, v := range t.Attributes {
		ov, ok := o.Attributes[k]
		if !ok {
			return false
		}
		if !v.Equal(ov) {
			return false
		}
	}
	return true
}
