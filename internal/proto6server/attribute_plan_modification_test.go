package proto6server

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestAttributeModifyPlan(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		req          tfsdk.ModifyAttributePlanRequest
		resp         ModifySchemaPlanResponse // Plan automatically copied from req
		expectedResp ModifySchemaPlanResponse
	}{
		"config-error": {
			req: tfsdk.ModifyAttributePlanRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.ListType{ElemType: types.StringType},
								Required: true,
							},
						},
					},
				},
				Plan: tfsdk.Plan{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
							},
						},
					},
				},
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
							},
						},
					},
				},
			},
			resp: ModifySchemaPlanResponse{},
			expectedResp: ModifySchemaPlanResponse{
				Diagnostics: diag.Diagnostics{
					diag.NewAttributeErrorDiagnostic(
						tftypes.NewAttributePath().WithAttributeName("test"),
						"Configuration Read Error",
						"An unexpected error was encountered trying to read an attribute from the configuration. This is always an error in the provider. Please report the following to the provider developer:\n\n"+
							"can't use tftypes.String<\"testvalue\"> as value of List with ElementType types.primitive, can only use tftypes.String values",
					),
				},
				Plan: tfsdk.Plan{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
							},
						},
					},
				},
			},
		},
		"config-error-previous-error": {
			req: tfsdk.ModifyAttributePlanRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.ListType{ElemType: types.StringType},
								Required: true,
							},
						},
					},
				},
				Plan: tfsdk.Plan{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
							},
						},
					},
				},
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
							},
						},
					},
				},
			},
			resp: ModifySchemaPlanResponse{
				Diagnostics: diag.Diagnostics{
					diag.NewErrorDiagnostic(
						"Previous error diag",
						"This was a previous error",
					),
				},
			},
			expectedResp: ModifySchemaPlanResponse{
				Diagnostics: diag.Diagnostics{
					diag.NewErrorDiagnostic(
						"Previous error diag",
						"This was a previous error",
					),
					diag.NewAttributeErrorDiagnostic(
						tftypes.NewAttributePath().WithAttributeName("test"),
						"Configuration Read Error",
						"An unexpected error was encountered trying to read an attribute from the configuration. This is always an error in the provider. Please report the following to the provider developer:\n\n"+
							"can't use tftypes.String<\"testvalue\"> as value of List with ElementType types.primitive, can only use tftypes.String values",
					),
				},
				Plan: tfsdk.Plan{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
							},
						},
					},
				},
			},
		},
		"plan-error": {
			req: tfsdk.ModifyAttributePlanRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
							},
						},
					},
				},
				Plan: tfsdk.Plan{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.ListType{ElemType: types.StringType},
								Required: true,
							},
						},
					},
				},
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
							},
						},
					},
				},
			},
			resp: ModifySchemaPlanResponse{},
			expectedResp: ModifySchemaPlanResponse{
				Diagnostics: diag.Diagnostics{
					diag.NewAttributeErrorDiagnostic(
						tftypes.NewAttributePath().WithAttributeName("test"),
						"Plan Read Error",
						"An unexpected error was encountered trying to read an attribute from the plan. This is always an error in the provider. Please report the following to the provider developer:\n\n"+
							"can't use tftypes.String<\"testvalue\"> as value of List with ElementType types.primitive, can only use tftypes.String values",
					),
				},
				Plan: tfsdk.Plan{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.ListType{ElemType: types.StringType},
								Required: true,
							},
						},
					},
				},
			},
		},
		"plan-error-previous-error": {
			req: tfsdk.ModifyAttributePlanRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
							},
						},
					},
				},
				Plan: tfsdk.Plan{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.ListType{ElemType: types.StringType},
								Required: true,
							},
						},
					},
				},
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
							},
						},
					},
				},
			},
			resp: ModifySchemaPlanResponse{
				Diagnostics: diag.Diagnostics{
					diag.NewErrorDiagnostic(
						"Previous error diag",
						"This was a previous error",
					),
				},
			},
			expectedResp: ModifySchemaPlanResponse{
				Diagnostics: diag.Diagnostics{
					diag.NewErrorDiagnostic(
						"Previous error diag",
						"This was a previous error",
					),
					diag.NewAttributeErrorDiagnostic(
						tftypes.NewAttributePath().WithAttributeName("test"),
						"Plan Read Error",
						"An unexpected error was encountered trying to read an attribute from the plan. This is always an error in the provider. Please report the following to the provider developer:\n\n"+
							"can't use tftypes.String<\"testvalue\"> as value of List with ElementType types.primitive, can only use tftypes.String values",
					),
				},
				Plan: tfsdk.Plan{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.ListType{ElemType: types.StringType},
								Required: true,
							},
						},
					},
				},
			},
		},
		"state-error": {
			req: tfsdk.ModifyAttributePlanRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
							},
						},
					},
				},
				Plan: tfsdk.Plan{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
							},
						},
					},
				},
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.ListType{ElemType: types.StringType},
								Required: true,
							},
						},
					},
				},
			},
			resp: ModifySchemaPlanResponse{},
			expectedResp: ModifySchemaPlanResponse{
				Diagnostics: diag.Diagnostics{
					diag.NewAttributeErrorDiagnostic(
						tftypes.NewAttributePath().WithAttributeName("test"),
						"State Read Error",
						"An unexpected error was encountered trying to read an attribute from the state. This is always an error in the provider. Please report the following to the provider developer:\n\n"+
							"can't use tftypes.String<\"testvalue\"> as value of List with ElementType types.primitive, can only use tftypes.String values",
					),
				},
				Plan: tfsdk.Plan{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
							},
						},
					},
				},
			},
		},
		"state-error-previous-error": {
			req: tfsdk.ModifyAttributePlanRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
							},
						},
					},
				},
				Plan: tfsdk.Plan{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
							},
						},
					},
				},
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.ListType{ElemType: types.StringType},
								Required: true,
							},
						},
					},
				},
			},
			resp: ModifySchemaPlanResponse{
				Diagnostics: diag.Diagnostics{
					diag.NewErrorDiagnostic(
						"Previous error diag",
						"This was a previous error",
					),
				},
			},
			expectedResp: ModifySchemaPlanResponse{
				Diagnostics: diag.Diagnostics{
					diag.NewErrorDiagnostic(
						"Previous error diag",
						"This was a previous error",
					),
					diag.NewAttributeErrorDiagnostic(
						tftypes.NewAttributePath().WithAttributeName("test"),
						"State Read Error",
						"An unexpected error was encountered trying to read an attribute from the state. This is always an error in the provider. Please report the following to the provider developer:\n\n"+
							"can't use tftypes.String<\"testvalue\"> as value of List with ElementType types.primitive, can only use tftypes.String values",
					),
				},
				Plan: tfsdk.Plan{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
							},
						},
					},
				},
			},
		},
		"no-plan-modifiers": {
			req: tfsdk.ModifyAttributePlanRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
							},
						},
					},
				},
				Plan: tfsdk.Plan{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
							},
						},
					},
				},
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
							},
						},
					},
				},
			},
			resp: ModifySchemaPlanResponse{},
			expectedResp: ModifySchemaPlanResponse{
				Plan: tfsdk.Plan{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
							},
						},
					},
				},
			},
		},
		"attribute-plan": {
			req: tfsdk.ModifyAttributePlanRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "TESTATTRONE"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									testAttrPlanValueModifierOne{},
									testAttrPlanValueModifierTwo{},
								},
							},
						},
					},
				},
				Plan: tfsdk.Plan{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "TESTATTRONE"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									testAttrPlanValueModifierOne{},
									testAttrPlanValueModifierTwo{},
								},
							},
						},
					},
				},
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "TESTATTRONE"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									testAttrPlanValueModifierOne{},
									testAttrPlanValueModifierTwo{},
								},
							},
						},
					},
				},
			},
			resp: ModifySchemaPlanResponse{},
			expectedResp: ModifySchemaPlanResponse{
				Plan: tfsdk.Plan{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "MODIFIED_TWO"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									testAttrPlanValueModifierOne{},
									testAttrPlanValueModifierTwo{},
								},
							},
						},
					},
				},
			},
		},
		"attribute-plan-previous-error": {
			req: tfsdk.ModifyAttributePlanRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "TESTATTRONE"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									testAttrPlanValueModifierOne{},
									testAttrPlanValueModifierTwo{},
								},
							},
						},
					},
				},
				Plan: tfsdk.Plan{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "TESTATTRONE"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									testAttrPlanValueModifierOne{},
									testAttrPlanValueModifierTwo{},
								},
							},
						},
					},
				},
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "TESTATTRONE"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									testAttrPlanValueModifierOne{},
									testAttrPlanValueModifierTwo{},
								},
							},
						},
					},
				},
			},
			resp: ModifySchemaPlanResponse{
				Diagnostics: diag.Diagnostics{
					diag.NewErrorDiagnostic(
						"Previous error diag",
						"This was a previous error",
					),
				},
			},
			expectedResp: ModifySchemaPlanResponse{
				Diagnostics: diag.Diagnostics{
					diag.NewErrorDiagnostic(
						"Previous error diag",
						"This was a previous error",
					),
				},
				Plan: tfsdk.Plan{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "MODIFIED_TWO"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									testAttrPlanValueModifierOne{},
									testAttrPlanValueModifierTwo{},
								},
							},
						},
					},
				},
			},
		},
		"requires-replacement": {
			req: tfsdk.ModifyAttributePlanRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "newtestvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									tfsdk.RequiresReplace(),
								},
							},
						},
					},
				},
				Plan: tfsdk.Plan{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "newtestvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									tfsdk.RequiresReplace(),
								},
							},
						},
					},
				},
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									tfsdk.RequiresReplace(),
								},
							},
						},
					},
				},
			},
			resp: ModifySchemaPlanResponse{},
			expectedResp: ModifySchemaPlanResponse{
				Plan: tfsdk.Plan{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "newtestvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
							},
						},
					},
				},
				RequiresReplace: []*tftypes.AttributePath{
					tftypes.NewAttributePath().WithAttributeName("test"),
				},
			},
		},
		"requires-replacement-previous-error": {
			req: tfsdk.ModifyAttributePlanRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "newtestvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									tfsdk.RequiresReplace(),
								},
							},
						},
					},
				},
				Plan: tfsdk.Plan{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "newtestvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									tfsdk.RequiresReplace(),
								},
							},
						},
					},
				},
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									tfsdk.RequiresReplace(),
								},
							},
						},
					},
				},
			},
			resp: ModifySchemaPlanResponse{
				Diagnostics: diag.Diagnostics{
					diag.NewErrorDiagnostic(
						"Previous error diag",
						"This was a previous error",
					),
				},
			},
			expectedResp: ModifySchemaPlanResponse{
				Diagnostics: diag.Diagnostics{
					diag.NewErrorDiagnostic(
						"Previous error diag",
						"This was a previous error",
					),
				},
				Plan: tfsdk.Plan{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "newtestvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
							},
						},
					},
				},
				RequiresReplace: []*tftypes.AttributePath{
					tftypes.NewAttributePath().WithAttributeName("test"),
				},
			},
		},
		"requires-replacement-passthrough": {
			req: tfsdk.ModifyAttributePlanRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "TESTATTRONE"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									testAttrPlanValueModifierOne{},
									tfsdk.RequiresReplace(),
								},
							},
						},
					},
				},
				Plan: tfsdk.Plan{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "TESTATTRONE"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									tfsdk.RequiresReplace(),
									testAttrPlanValueModifierOne{},
								},
							},
						},
					},
				},
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "TESTATTRONE"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									tfsdk.RequiresReplace(),
									testAttrPlanValueModifierOne{},
								},
							},
						},
					},
				},
			},
			resp: ModifySchemaPlanResponse{},
			expectedResp: ModifySchemaPlanResponse{
				Plan: tfsdk.Plan{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "TESTATTRTWO"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
							},
						},
					},
				},
				RequiresReplace: []*tftypes.AttributePath{
					tftypes.NewAttributePath().WithAttributeName("test"),
				},
			},
		},
		"requires-replacement-unset": {
			req: tfsdk.ModifyAttributePlanRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									tfsdk.RequiresReplace(),
									testRequiresReplaceFalseModifier{},
								},
							},
						},
					},
				},
				Plan: tfsdk.Plan{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									tfsdk.RequiresReplace(),
									testRequiresReplaceFalseModifier{},
								},
							},
						},
					},
				},
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									tfsdk.RequiresReplace(),
									testRequiresReplaceFalseModifier{},
								},
							},
						},
					},
				},
			},
			resp: ModifySchemaPlanResponse{},
			expectedResp: ModifySchemaPlanResponse{
				Plan: tfsdk.Plan{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
							},
						},
					},
				},
			},
		},
		"warnings": {
			req: tfsdk.ModifyAttributePlanRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "TESTDIAG"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									testWarningDiagModifier{},
									testWarningDiagModifier{},
								},
							},
						},
					},
				},
				Plan: tfsdk.Plan{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "TESTDIAG"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									testWarningDiagModifier{},
									testWarningDiagModifier{},
								},
							},
						},
					},
				},
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "TESTDIAG"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									testWarningDiagModifier{},
									testWarningDiagModifier{},
								},
							},
						},
					},
				},
			},
			resp: ModifySchemaPlanResponse{},
			expectedResp: ModifySchemaPlanResponse{
				Diagnostics: diag.Diagnostics{
					// Diagnostics.Append() deduplicates, so the warning will only
					// be here once unless the test implementation is changed to
					// different modifiers or the modifier itself is changed.
					diag.NewWarningDiagnostic(
						"Warning diag",
						"This is a warning",
					),
				},
				Plan: tfsdk.Plan{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "TESTDIAG"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									testWarningDiagModifier{},
									testWarningDiagModifier{},
								},
							},
						},
					},
				},
			},
		},
		"warnings-previous-error": {
			req: tfsdk.ModifyAttributePlanRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "TESTDIAG"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									testWarningDiagModifier{},
									testWarningDiagModifier{},
								},
							},
						},
					},
				},
				Plan: tfsdk.Plan{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "TESTDIAG"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									testWarningDiagModifier{},
									testWarningDiagModifier{},
								},
							},
						},
					},
				},
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "TESTDIAG"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									testWarningDiagModifier{},
									testWarningDiagModifier{},
								},
							},
						},
					},
				},
			},
			resp: ModifySchemaPlanResponse{
				Diagnostics: diag.Diagnostics{
					diag.NewErrorDiagnostic(
						"Previous error diag",
						"This was a previous error",
					),
				},
			},
			expectedResp: ModifySchemaPlanResponse{
				Diagnostics: diag.Diagnostics{
					diag.NewErrorDiagnostic(
						"Previous error diag",
						"This was a previous error",
					),
					// Diagnostics.Append() deduplicates, so the warning will only
					// be here once unless the test implementation is changed to
					// different modifiers or the modifier itself is changed.
					diag.NewWarningDiagnostic(
						"Warning diag",
						"This is a warning",
					),
				},
				Plan: tfsdk.Plan{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "TESTDIAG"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									testWarningDiagModifier{},
									testWarningDiagModifier{},
								},
							},
						},
					},
				},
			},
		},
		"error": {
			req: tfsdk.ModifyAttributePlanRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "TESTDIAG"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									testErrorDiagModifier{},
									testErrorDiagModifier{},
								},
							},
						},
					},
				},
				Plan: tfsdk.Plan{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "TESTDIAG"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									testErrorDiagModifier{},
									testErrorDiagModifier{},
								},
							},
						},
					},
				},
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "TESTDIAG"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									testErrorDiagModifier{},
									testErrorDiagModifier{},
								},
							},
						},
					},
				},
			},
			resp: ModifySchemaPlanResponse{},
			expectedResp: ModifySchemaPlanResponse{
				Diagnostics: diag.Diagnostics{
					diag.NewErrorDiagnostic(
						"Error diag",
						"This is an error",
					),
				},
				Plan: tfsdk.Plan{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "TESTDIAG"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									testErrorDiagModifier{},
									testErrorDiagModifier{},
								},
							},
						},
					},
				},
			},
		},
		"error-previous-error": {
			req: tfsdk.ModifyAttributePlanRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "TESTDIAG"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									testErrorDiagModifier{},
									testErrorDiagModifier{},
								},
							},
						},
					},
				},
				Plan: tfsdk.Plan{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "TESTDIAG"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									testErrorDiagModifier{},
									testErrorDiagModifier{},
								},
							},
						},
					},
				},
				State: tfsdk.State{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "TESTDIAG"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									testErrorDiagModifier{},
									testErrorDiagModifier{},
								},
							},
						},
					},
				},
			},
			resp: ModifySchemaPlanResponse{
				Diagnostics: diag.Diagnostics{
					diag.NewErrorDiagnostic(
						"Previous error diag",
						"This was a previous error",
					),
				},
			},
			expectedResp: ModifySchemaPlanResponse{
				Diagnostics: diag.Diagnostics{
					diag.NewErrorDiagnostic(
						"Previous error diag",
						"This was a previous error",
					),
					diag.NewErrorDiagnostic(
						"Error diag",
						"This is an error",
					),
				},
				Plan: tfsdk.Plan{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "TESTDIAG"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									testErrorDiagModifier{},
									testErrorDiagModifier{},
								},
							},
						},
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			attribute, err := tc.req.Config.Schema.AttributeAtPath(tc.req.AttributePath)

			if err != nil {
				t.Fatalf("Unexpected error getting %s", err)
			}

			tc.resp.Plan = tc.req.Plan

			AttributeModifyPlan(context.Background(), attribute, tc.req, &tc.resp)

			if diff := cmp.Diff(tc.expectedResp, tc.resp); diff != "" {
				t.Errorf("Unexpected response (-wanted, +got): %s", diff)
			}
		})
	}
}