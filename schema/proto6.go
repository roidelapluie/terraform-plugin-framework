package schema

import (
	"context"
	"errors"
	"sort"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// ToProto6 returns the *tfprotov6.Schema equivalent of a Schema. At least
// one attribute must be set in the schema, or an error will be returned.
func (s Schema) ToProto6(ctx context.Context) (*tfprotov6.Schema, error) {
	result := &tfprotov6.Schema{
		Version: s.Version,
	}
	var attrs []*tfprotov6.SchemaAttribute
	for name, attr := range s.Attributes {
		a, err := attr.ToProto6(ctx, name, tftypes.NewAttributePath().WithAttributeName(name))
		if err != nil {
			return nil, err
		}
		attrs = append(attrs, a)
	}
	sort.Slice(attrs, func(i, j int) bool {
		if attrs[i] == nil {
			return true
		}
		if attrs[j] == nil {
			return false
		}
		return attrs[i].Name < attrs[j].Name
	})
	if len(attrs) < 1 {
		return nil, errors.New("must have at least one attribute in the schema")
	}
	result.Block = &tfprotov6.SchemaBlock{
		// core doesn't do anything with version, as far as I can tell,
		// so let's not set it.
		Attributes: attrs,
		Deprecated: s.DeprecationMessage != "",
	}
	if s.Description != "" {
		result.Block.Description = s.Description
		result.Block.DescriptionKind = tfprotov6.StringKindPlain
	}
	if s.MarkdownDescription != "" {
		result.Block.Description = s.MarkdownDescription
		result.Block.DescriptionKind = tfprotov6.StringKindMarkdown
	}
	return result, nil
}

// Attribute returns the *tfprotov6.SchemaAttribute equivalent of a
// Attribute. Errors will be tftypes.AttributePathErrors based on
// `path`. `name` is the name of the attribute.
func (t Attribute) ToProto6(ctx context.Context, name string, path *tftypes.AttributePath) (*tfprotov6.SchemaAttribute, error) {
	a := &tfprotov6.SchemaAttribute{
		Name:      name,
		Required:  t.Required,
		Optional:  t.Optional,
		Computed:  t.Computed,
		Sensitive: t.Sensitive,
	}
	if t.DeprecationMessage != "" {
		a.Deprecated = true
	}
	if t.Description != "" {
		a.Description = t.Description
		a.DescriptionKind = tfprotov6.StringKindPlain
	}
	if t.MarkdownDescription != "" {
		a.Description = t.MarkdownDescription
		a.DescriptionKind = tfprotov6.StringKindMarkdown
	}
	if t.Type != nil && t.Attributes == nil {
		a.Type = t.Type.TerraformType(ctx)
	} else if t.Attributes != nil && len(t.Attributes.GetAttributes()) > 0 && t.Type == nil {
		object := &tfprotov6.SchemaObject{
			MinItems: t.Attributes.GetMinItems(),
			MaxItems: t.Attributes.GetMaxItems(),
		}
		nm := t.Attributes.GetNestingMode()
		switch nm {
		case NestingModeSingle:
			object.Nesting = tfprotov6.SchemaObjectNestingModeSingle
		case NestingModeList:
			object.Nesting = tfprotov6.SchemaObjectNestingModeList
		case NestingModeSet:
			object.Nesting = tfprotov6.SchemaObjectNestingModeSet
		case NestingModeMap:
			object.Nesting = tfprotov6.SchemaObjectNestingModeMap
		default:
			return nil, path.NewErrorf("unrecognized nesting mode %v", nm)
		}
		attrs := t.Attributes.GetAttributes()
		for nestedName, nestedAttr := range attrs {
			nestedA, err := nestedAttr.ToProto6(ctx, nestedName, path.WithAttributeName(nestedName))
			if err != nil {
				return nil, err
			}
			object.Attributes = append(object.Attributes, nestedA)
		}
		sort.Slice(object.Attributes, func(i, j int) bool {
			if object.Attributes[i] == nil {
				return true
			}
			if object.Attributes[j] == nil {
				return false
			}
			return object.Attributes[i].Name < object.Attributes[j].Name
		})
		a.NestedType = object
	} else if t.Attributes != nil && len(t.Attributes.GetAttributes()) > 0 && t.Type != nil {
		return nil, path.NewErrorf("can't have both Attributes and Type set")
	} else if (t.Attributes == nil || len(t.Attributes.GetAttributes()) < 1) && t.Type == nil {
		return nil, path.NewErrorf("must have Attributes or Type set")
	}
	return a, nil
}
