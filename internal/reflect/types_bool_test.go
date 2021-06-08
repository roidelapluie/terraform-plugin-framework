package reflect

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type testBoolType struct{}

var _ attr.Type = testBoolType{}

func (t testBoolType) TerraformType(_ context.Context) tftypes.Type {
	return tftypes.Bool
}

func (t testBoolType) ValueFromTerraform(_ context.Context, in tftypes.Value) (attr.Value, error) {
	if !in.Type().Is(tftypes.Bool) {
		return nil, fmt.Errorf("unexpected type %s", tftypes.Bool)
	}
	if !in.IsKnown() {
		return &testBoolValue{Unknown: true}, nil
	}
	if in.IsNull() {
		return &testBoolValue{Null: true}, nil
	}
	var a bool
	err := in.As(&a)
	if err != nil {
		return nil, err
	}
	return &testBoolValue{Value: a}, nil
}

func (t testBoolType) Equal(other attr.Type) bool {
	_, ok := other.(testBoolType)
	return ok
}

type testBoolValue struct {
	Unknown bool
	Null    bool
	Value   bool
	setErr  error
}

var _ attr.Value = &testBoolValue{}

func (t *testBoolValue) ToTerraformValue(_ context.Context) (interface{}, error) {
	if t.Unknown {
		return tftypes.UnknownValue, nil
	}
	if t.Null {
		return nil, nil
	}
	return t.Value, nil
}

func (t *testBoolValue) SetTerraformValue(_ context.Context, in tftypes.Value) error {
	if t.setErr != nil {
		return t.setErr
	}
	t.Unknown = false
	t.Null = false
	t.Value = false
	if !in.Type().Is(tftypes.Bool) {
		return fmt.Errorf("unexpected type %s", tftypes.Bool)
	}
	if !in.IsKnown() {
		t.Unknown = true
		return nil
	}
	if !in.IsNull() {
		t.Null = true
		return nil
	}
	return in.As(&t.Value)
}

func (t *testBoolValue) Equal(other attr.Value) bool {
	if t == nil && other == nil {
		return true
	}
	if t == nil || other == nil {
		return false
	}
	o, ok := other.(*testBoolValue)
	if !ok {
		return false
	}
	if t.Unknown != o.Unknown {
		return false
	}
	if t.Null != o.Null {
		return false
	}
	if t.Value != o.Value {
		return false
	}
	return true
}
