// Package models defines the domain models that model field relations and manipulation.
package models

import (
	"fmt"
)

// Field represents a field to be copied to/from.
// A field's struct properties are set in the parser unless its stated otherwise.
type Field struct {
	// The variable name the field is assigned for assignment.
	// This value will always be unique in the context of the application.
	// Type variable names do not contain '.' (i.e 'tA' in 'tA.UserID')
	// Field variable names are defined by their specifier (i.e '.UserID' in 'domain.Account.UserID').
	VariableName string

	// The package the field is defined in.
	Package string

	// Import path
	Import string

	// The name of the field (i.e `ID` in `ID int`).
	Name string

	// The type definition of the field (i.e `int` in `ID int`, string, `Logger` in `log.Logger`).
	Definition string

	// The type definition of the field (i.e `int` in `ID int`, string, `Logger` in `log.Logger`).
	OrigDefinition string

	// The pointer(s) of the field in string format (i.e **).
	Pointer string

	// The type or field that contains this field.
	Parent *Field

	// The field that this field will be copied from (or nil).
	// Set in the matcher.
	From *Field

	// The field that this field will be copied to (or nil).
	// Set in the matcher.
	To *Field

	// The fields of this field.
	Fields []*Field

	// The custom options of a field.
	Options FieldOptions
}

// FieldOptions represent options for a Field.
type FieldOptions struct {
	// The function the field is converted with (as a parameter).
	Convert string

	// The field to map this field to, if any.
	Map string

	// The level at which sub-fields are discovered.
	Depth int

	// Whether the field should be deepcopied.
	Deepcopy bool
}

// IsType returns whether the field is a type.
func (f *Field) IsType() bool {
	return f.Parent == nil
}

// IsPointer returns the type is a pointer of a type definition.
func (f *Field) IsPointer() bool {
	return f.Pointer == ""
}

// FullName gets the full name of a field including its parents (i.e domain.Account.User.ID).
func (f *Field) FullName(name string) string {
	if !f.IsType() {
		// add names in reverse order
		if name == "" {
			name = f.Name
		} else {
			name = f.Name + "." + name
		}

		return f.Parent.FullName(name)
	}

	if name != "" {
		name = "." + name
	}

	return fmt.Sprintf("%v%v.%v%v", f.Pointer, f.Package, f.Name, name)
}

// FullVariableName gets the full variable name of a field (i.e tA.User.UserID).
func (f *Field) FullVariableName(name string) string {
	if !f.IsType() {
		return f.Parent.FullVariableName(f.VariableName + name)
	}

	return f.VariableName + name
}

// AllFields gets all the fields in the scope of a field (including itself).
func (f *Field) AllFields(fields []*Field) []*Field {
	fields = append(fields, f)

	if len(f.Fields) != 0 {
		for i := 0; i < len(f.Fields); i++ {
			fields = f.Fields[i].AllFields(fields)
		}
	}

	return fields
}

func (f *Field) String() string {
	var direction, parent string

	if f.From != nil {
		direction = "To"
	}

	if f.To != nil {
		if direction != "" {
			direction += " and "
		}

		direction += "From"
	}

	if direction == "" {
		direction = "Unpointed"
	}

	if f.Parent != nil {
		parent = f.Parent.FullName("")
	}

	return fmt.Sprintf("%v Field %q of Definition %q: Parent %q Fields[%v]", direction, f.FullName(""), f.Definition, parent, len(f.Fields))
}
