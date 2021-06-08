package reflect

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type testMapType struct {
	ElemType attr.Type
}

func (t testMapType) WithElementType(typ attr.Type) attr.ElementType {
	return testMapType{
		ElemType: typ,
	}
}

func (t testMapType) ElementType() attr.Type {
	return t.ElemType
}

func (t testMapType) TerraformType(ctx context.Context) tftypes.Type {
	return tftypes.Map{
		AttributeType: t.ElemType.TerraformType(ctx),
	}
}

func (t testMapType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if !in.Type().Is(tftypes.Map{}) {
		return nil, fmt.Errorf("unexpected type %s", in.Type())
	}
	result := &testMapValue{
		ElemType: t.ElemType,
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
	result.Elements = map[string]attr.Value{}
	for k, v := range inVals {
		res, err := t.ElemType.ValueFromTerraform(ctx, v)
		if err != nil {
			return nil, fmt.Errorf("error converting %q to attr.Value: %w", k, err)
		}
		result.Elements[k] = res
	}
	return result, nil
}

func (t testMapType) Equal(other attr.Type) bool {
	o, ok := other.(testMapType)
	if !ok {
		return false
	}
	return t.ElemType.Equal(o.ElemType)
}

var _ attr.ElementType = testMapType{}

type testMapValue struct {
	ElemType attr.Type
	Elements map[string]attr.Value
	Unknown  bool
	Null     bool
}

var _ attr.Value = &testMapValue{}

func (t *testMapValue) ToTerraformValue(ctx context.Context) (interface{}, error) {
	if t.Unknown {
		return tftypes.UnknownValue, nil
	}
	if t.Null {
		return nil, nil
	}
	resultVals := map[string]tftypes.Value{}
	for k, v := range t.Elements {
		val, err := v.ToTerraformValue(ctx)
		if err != nil {
			return nil, fmt.Errorf("error generating terraform value for %q: %w", k, err)
		}
		err = tftypes.ValidateValue(t.ElemType.TerraformType(ctx), val)
		if err != nil {
			return nil, fmt.Errorf("error validating terraform value for %q: %w", k, err)
		}
		resultVals[k] = tftypes.NewValue(t.ElemType.TerraformType(ctx), val)
	}
	return resultVals, nil
}

func (t *testMapValue) Equal(other attr.Value) bool {
	if t == nil && other == nil {
		return true
	}
	if t == nil || other == nil {
		return false
	}
	o, ok := other.(*testMapValue)
	if !ok {
		return false
	}
	if t.Null != o.Null {
		return false
	}
	if t.Unknown != o.Unknown {
		return false
	}
	if !t.ElemType.Equal(o.ElemType) {
		return false
	}
	if len(t.Elements) != len(o.Elements) {
		return false
	}
	for k, v := range t.Elements {
		ov, ok := o.Elements[k]
		if !ok {
			return false
		}
		if !v.Equal(ov) {
			return false
		}
	}
	return true
}
