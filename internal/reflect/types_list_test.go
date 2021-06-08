package reflect

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type testListType struct {
	ElemType attr.Type
}

func (t testListType) TerraformType(ctx context.Context) tftypes.Type {
	return tftypes.List{
		ElementType: t.ElemType.TerraformType(ctx),
	}
}

func (t testListType) WithElementType(typ attr.Type) attr.ElementType {
	return testListType{
		ElemType: typ,
	}
}

func (t testListType) ElementType() attr.Type {
	return t.ElemType
}

func (t testListType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if !in.Type().Is(tftypes.List{}) {
		return nil, fmt.Errorf("unexpected type %s", in.Type())
	}
	result := &testListValue{
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
	inVals := []tftypes.Value{}
	err := in.As(&inVals)
	if err != nil {
		return nil, err
	}
	result.Elements = []attr.Value{}
	for pos, v := range inVals {
		res, err := t.ElemType.ValueFromTerraform(ctx, v)
		if err != nil {
			return nil, fmt.Errorf("error converting element %d to attr.Value: %w", pos, err)
		}
		result.Elements = append(result.Elements, res)
	}
	return result, nil
}

func (t testListType) Equal(other attr.Type) bool {
	o, ok := other.(testListType)
	if !ok {
		return false
	}
	return t.ElemType.Equal(o.ElemType)
}

var _ attr.ElementType = testListType{}

type testListValue struct {
	ElemType attr.Type
	Elements []attr.Value
	Unknown  bool
	Null     bool
}

var _ attr.Value = &testListValue{}

func (t *testListValue) ToTerraformValue(ctx context.Context) (interface{}, error) {
	if t.Unknown {
		return tftypes.UnknownValue, nil
	}
	if t.Null {
		return nil, nil
	}
	resultVals := []tftypes.Value{}
	for pos, v := range t.Elements {
		val, err := v.ToTerraformValue(ctx)
		if err != nil {
			return nil, fmt.Errorf("error generating terraform value for element %d: %w", pos, err)
		}
		err = tftypes.ValidateValue(t.ElemType.TerraformType(ctx), val)
		if err != nil {
			return nil, fmt.Errorf("error validating terraform value for element %d: %w", pos, err)
		}
		resultVals = append(resultVals, tftypes.NewValue(t.ElemType.TerraformType(ctx), val))
	}
	return resultVals, nil
}

func (t *testListValue) Equal(other attr.Value) bool {
	if t == nil && other == nil {
		return true
	}
	if t == nil || other == nil {
		return false
	}
	o, ok := other.(*testListValue)
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
	for pos, v := range t.Elements {
		if !v.Equal(o.Elements[pos]) {
			return false
		}
	}
	return true
}
