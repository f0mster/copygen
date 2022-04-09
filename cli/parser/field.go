package parser

import (
	"fmt"
	"go/types"

	"github.com/switchupcb/copygen/cli/models"
	"github.com/switchupcb/copygen/cli/parser/options"
)

// fieldParser represents the parameters required to parse a types.Type into a *models.Field.
type fieldParser struct {
	// field represents the current field being built.
	field *models.Field

	// parent represents the parent of the field parse.
	parent *models.Field

	// cyclic is a key value cache used to prevent cyclic fields from unnecessary duplication or stack overflow.
	cyclic map[string]bool

	// container represents a field's container.
	container string

	// options represents the field options defined above the models.Function
	options []*options.Option
}

// parseField parses a types.Type into a *models.Field.
func (fp fieldParser) parseField(typ types.Type) *models.Field {
	if fp.field == nil {
		fp.field = &models.Field{Parent: fp.parent}
	}

	switch x := typ.(type) {

	// Basic Types
	// https://go.googlesource.com/example/+/HEAD/gotypes#basic-types
	case *types.Basic:
		setFieldVariableName(fp.field, "."+x.Name())
		fp.field.Definition = x.Name()

	// Named Types (Alias)
	// https://go.googlesource.com/example/+/HEAD/gotypes#named-types
	case *types.Named:
		setFieldImportAndPackage(fp.field, x.Obj().Pkg().Path(), x.Obj().Pkg().Name())
		setFieldVariableName(fp.field, "."+x.Obj().Name())
		fp.field.Name = x.Obj().Name()
		return fp.parseField(x.Underlying())

	// Simple Composite Types
	// https://go.googlesource.com/example/+/HEAD/gotypes#simple-composite-types
	case *types.Pointer:
		fp.field.Container += "*"
		return fp.parseField(x.Elem())

	case *types.Array:
		setFieldVariableName(fp.field, "."+alphastring(x.String()))
		fp.field.Definition = x.String()
		fp.field.Container += "[" + fmt.Sprint(x.Len()) + "]"

	case *types.Slice:
		setFieldVariableName(fp.field, "."+alphastring(x.String()))
		fp.field.Definition = x.String()
		fp.field.Container += "[]"

	case *types.Map:
		setFieldVariableName(fp.field, "."+alphastring(x.String()))
		fp.field.Definition = x.String()
		fp.field.Container += "map"

	case *types.Chan:
		setFieldVariableName(fp.field, "."+alphastring(x.String()))
		fp.field.Definition = x.String()
		fp.field.Container += "chan"

	// Struct Types
	// https://go.googlesource.com/example/+/HEAD/gotypes#struct-types
	case *types.Struct:
		fp.field.Definition = "struct"
		for i := 0; i < x.NumFields(); i++ {
			subfield := &models.Field{
				VariableName: "." + x.Field(i).Name(),
				Name:         x.Field(i).Name(),
				Tag:          x.Tag(i),
				Parent:       fp.field,
			}
			setFieldImportAndPackage(subfield, x.Field(i).Pkg().Path(), x.Field(i).Pkg().Name())

			// a cyclic subfield (with the same type as its parent) is never fully assigned.
			if !fp.cyclic[subfield.Import+subfield.Package+subfield.Name] {
				subfieldParser := &fieldParser{
					field:     subfield,
					parent:    nil,
					container: "",
					options:   fp.options,
					cyclic:    fp.cyclic,
				}

				// sets the definition, container, and fields.
				fp.cyclic[subfield.Import+subfield.Package+subfield.Name] = true
				subfield = subfieldParser.parseField(x.Field(i).Type())
			}

			fp.field.Fields = append(fp.field.Fields, subfield)
		}

	// Function
	// https://go.googlesource.com/example/+/HEAD/gotypes#function-and-method-types
	case *types.Signature:
		setFieldVariableName(fp.field, "."+alphastring(x.String()))
		fp.field.Definition = x.String()

	// Interface Types
	// https://go.googlesource.com/example/+/HEAD/gotypes#interface-types
	case *types.Interface:
		fp.field.Definition = "interface"
		for i := 0; i < x.NumMethods(); i++ {
			subfield := &models.Field{
				VariableName: "." + x.Method(i).Name(),
				Name:         x.Method(i).Name(),
				Parent:       fp.field,
			}
			setFieldImportAndPackage(subfield, x.Method(i).Pkg().Path(), x.Method(i).Pkg().Name())

			subfieldParser := &fieldParser{
				field:     subfield,
				parent:    nil,
				container: "",
				options:   fp.options,
				cyclic:    fp.cyclic,
			}

			// sets the definition, container, and fields.
			subfield = subfieldParser.parseField(x.Method(i).Type())
			fp.field.Fields = append(fp.field.Fields, subfield)
		}
	}

	setFieldOptions(fp.field, fp.options)
	fp.cyclic[fp.field.Import+fp.field.Package+fp.field.Name] = true
	return fp.field
}

// setFieldVariableName sets a field's variable name.
func setFieldVariableName(field *models.Field, varname string) {
	if field.VariableName == "" {
		field.VariableName = varname
	}
}

// alphastring only returns alphabetic characters (English) in a string.
func alphastring(s string) string {
	bytes := []byte(s)
	i := 0
	for _, b := range bytes {
		if ('a' <= b && b <= 'z') || ('A' <= b && b <= 'Z') || b == ' ' {
			bytes[i] = b
			i++
		}
	}

	return string(bytes[:i])
}

// setFieldImportAndPackage sets the import and package of a field.
func setFieldImportAndPackage(field *models.Field, path string, varname string) {
	field.Import = path
	if ignorepkgpath != path {
		if _, ok := aliasImportMap[path]; ok {
			field.Package = aliasImportMap[path]
		} else {
			field.Package = varname
		}
	}
}

// setFieldOptions sets a field's (and its subfields) options.
func setFieldOptions(field *models.Field, fieldoptions []*options.Option) {
	for _, option := range fieldoptions {
		switch option.Category {

		case options.CategoryConvert:
			options.SetConvert(field, *option)

		case options.CategoryDepth:
			options.SetDepth(field, *option)

		case options.CategoryDeepcopy:
			options.SetDeepcopy(field, *option)

		case options.CategoryMap:
			options.SetMap(field, *option)

		}
	}
}
