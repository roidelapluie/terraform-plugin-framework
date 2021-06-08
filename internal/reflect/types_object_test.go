package reflect

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type testObjectType struct {
	AttrTypes map[string]attr.Type
}

func (t testObjectType) WithAttributeTypes(types map[string]attr.Type) attr.AttributesType {
	return testObjectType{
		AttrTypes: types,
	}
}

func (t testObjectType) AttributeTypes() map[string]attr.Type {
	return t.AttrTypes
}

func (t testObjectType) TerraformType(ctx context.Context) tftypes.Type {
	types := map[string]tftypes.Type{}
	for k, v := range t.AttrTypes {
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
		AttrTypes: t.AttrTypes,
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
	if len(inVals) != len(t.AttrTypes) {
		return nil, fmt.Errorf("expected value to have %d attributes, has %d: %s", len(t.AttrTypes), len(inVals), in)
	}
	result.Attributes = map[string]attr.Value{}
	for k, v := range t.AttrTypes {
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
	if len(o.AttrTypes) != len(t.AttrTypes) {
		return false
	}
	for k, v := range t.AttrTypes {
		ov, ok := o.AttrTypes[k]
		if !ok {
			return false
		}
		if !ov.Equal(v) {
			return false
		}
	}
	return true
}

var _ attr.AttributesType = testObjectType{}

type testObjectValue struct {
	AttrTypes  map[string]attr.Type
	Attributes map[string]attr.Value
	Unknown    bool
	Null       bool
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
		typ, ok := t.AttrTypes[k]
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
	if len(t.AttrTypes) != len(o.AttrTypes) {
		return false
	}
	for k, v := range t.AttrTypes {
		ov, ok := o.AttrTypes[k]
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
