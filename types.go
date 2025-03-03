package graphql

import (
	"fmt"
)

type (
	IntrospectionResponse struct {
		Data Data `json:"data"`
	}

	Data struct {
		Schema Schema `json:"__schema"`
	}

	NameStruct struct {
		Name string `json:"name,omitempty"`
	}

	Schema struct {
		QueryType        NameStruct `json:"queryType,omitempty"`
		MutationType     NameStruct `json:"mutationType,omitempty"`
		SubscriptionType NameStruct `json:"subscriptionType,omitempty"`
		Types            []FullType `json:"types,omitempty"`
		Directives       []Field    `json:"directives,omitempty"`
		Queries          []Field    `json:"queries,omitempty"`
		Mutations        []Field    `json:"mutations,omitempty"`
	}
)

func (s *Schema) GetQueries() []Field {
	tmp := make([]Field, 0)
	for _, v := range s.Types {
		if v.Kind == "OBJECT" && v.Name == s.QueryType.Name {
			for _, f := range v.Fields {
				tmp = append(tmp, f)
			}
		}
	}
	return tmp
}

func (s *Schema) GetQuery(name string) Field {
	for _, v := range s.Types {
		if v.Kind == "OBJECT" && v.Name == s.QueryType.Name {
			for _, f := range v.Fields {
				if f.Name == name {
					return f
				}
			}
		}
	}
	return Field{}
}

func (s *Schema) GetMutations() []Field {
	tmp := make([]Field, 0)
	for _, v := range s.Types {
		if v.Kind == "OBJECT" && v.Name == s.MutationType.Name {
			for _, f := range v.Fields {
				tmp = append(tmp, f)
			}
		}
	}
	return tmp
}

func (s *Schema) GetSubscriptions() []Field {
	tmp := make([]Field, 0)
	for _, v := range s.Types {
		if v.Kind == "OBJECT" && v.Name == s.SubscriptionType.Name {
			for _, f := range v.Fields {
				tmp = append(tmp, f)
			}
		}
	}
	return tmp
}

type InputValue struct {
	FieldProperties
	DefaultValue string `json:"defaultValue,omitempty"`
}

type FullType struct {
	Kind          string       `json:"kind,omitempty"`
	Name          string       `json:"name,omitempty"`
	Description   interface{}  `json:"description,omitempty"`
	Fields        []Field      `json:"fields,omitempty"`
	InputFields   []InputValue `json:"inputFields,omitempty"`
	Interfaces    []TypeRef    `json:"interfaces,omitempty"`
	EnumValues    []EnumValue  `json:"enumValues,omitempty"`
	PossibleTypes []TypeRef    `json:"possibleTypes,omitempty"`
}

func (typ FullType) String() string {
	var schemaSnippet string
	if typ.Kind == "OBJECT" {
		schemaSnippet = fmt.Sprintf("type %s {\n", typ.Name)
	} else if typ.Kind == "INTERFACE" {
		schemaSnippet = fmt.Sprintf("interface %s {\n", typ.Name)
	} else if typ.Kind == "INPUT_OBJECT" {
		schemaSnippet = fmt.Sprintf("input %s {\n", typ.Name)
	} else if typ.Kind == "ENUM" {
		schemaSnippet = fmt.Sprintf("enum %s {\n", typ.Name)
	} else if typ.Kind == "SCALAR" {
		if typ.Name == "Boolean" || typ.Name == "Float" || typ.Name == "ID" || typ.Name == "Int" || typ.Name == "String" {
			return ""
		}
		schemaSnippet = fmt.Sprintf("scalar %s", typ.Name)
		return schemaSnippet
	} else {
		// schemaSnippet = fmt.Sprintf("#%s\ninput %s {\n", typ.Kind, typ.Name)
		schemaSnippet = fmt.Sprintf("#%#v", typ)
	}

	for _, field := range typ.InputFields {
		schemaSnippet += fmt.Sprintf("\t%s: %s\n", field.Name, field.Type.String())
	}
	for _, field := range typ.Fields {
		schemaSnippet += fmt.Sprintf("\t%s: %s\n", field.Name, field.Type.String())
	}
	for _, field := range typ.EnumValues {
		schemaSnippet += fmt.Sprintf("\t%s\n", field.Name)
	}

	schemaSnippet += "}"
	return schemaSnippet
}

type EnumValue struct {
	Name              string `json:"name,omitempty"`
	Description       string `json:"description,omitempty"`
	IsDeprecated      bool   `json:"isDeprecated,omitempty"`
	DeprecationReason bool   `json:"deprecationReason,omitempty"`
}

type FieldProperties struct {
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Type        TypeRef `json:"type,omitempty"`
}

type Field struct {
	FieldProperties
	Args []InputValue `json:"args,omitempty"`
}

type TypeRef struct {
	Kind   string  `json:"kind,omitempty"`
	Name   string  `json:"name,omitempty"`
	OfType *OfType `json:"ofType,omitempty"`
}

func (t TypeRef) String() string {
	return resolveType(t.Kind, t.Name, func() string {
		if t.OfType == nil {
			return ""
		}
		if t.IsMultiple() {
			return fmt.Sprintf("[%s]", t.OfType.String())
		}
		return t.OfType.String()
	})
}

func (t TypeRef) IsMultiple() bool {
	if t.Kind == "LIST" {
		return true
	}
	if t.OfType == nil {
		return false
	}
	return t.OfType.IsMultiple()
}

type OfType struct {
	Kind   string   `json:"kind,omitempty"`
	Name   string   `json:"name,omitempty"`
	OfType *OfType2 `json:"ofType,omitempty"`
}

func (t OfType) String() string {
	return resolveType(t.Kind, t.Name, func() string { return t.OfType.String() })
}

func (t OfType) IsMultiple() bool {
	if t.Kind == "LIST" {
		return true
	}
	if t.OfType == nil {
		return false
	}
	return t.OfType.IsMultiple()
}

type OfType2 struct {
	Kind   string   `json:"kind,omitempty"`
	Name   string   `json:"name,omitempty"`
	OfType *OfType3 `json:"ofType,omitempty"`
}

func (t OfType2) String() string {
	return resolveType(t.Kind, t.Name, func() string { return t.OfType.String() })
}

func (t OfType2) IsMultiple() bool {
	if t.Kind == "LIST" {
		return true
	}
	if t.OfType == nil {
		return false
	}
	return t.OfType.IsMultiple()
}

type OfType3 struct {
	Kind   string   `json:"kind,omitempty"`
	Name   string   `json:"name,omitempty"`
	OfType *OfType4 `json:"ofType,omitempty"`
}

func (t OfType3) String() string {
	return resolveType(t.Kind, t.Name, func() string { return t.OfType.String() })
}

func (t OfType3) IsMultiple() bool {
	if t.Kind == "LIST" {
		return true
	}
	if t.OfType == nil {
		return false
	}
	return t.OfType.IsMultiple()
}

type OfType4 struct {
	Kind string `json:"kind,omitempty"`
	Name string `json:"name,omitempty"`
}

func (t OfType4) String() string {
	return resolveType(t.Kind, t.Name, func() string { return "" })
}

func (t OfType4) IsMultiple() bool {
	if t.Kind == "LIST" {
		return true
	}
	return false
}

func resolveType(kind string, name string, next func() string) string {
	if kind == "SCALAR" || kind == "INTERFACE" {
		return name
	}
	if kind == "OBJECT" || kind == "INPUT_OBJECT" || kind == "ENUM" {
		return name
	}
	if kind == "LIST" {
		// return "[" + next() + "]"
		return next()
	}
	if kind == "NON_NULL" {
		return next() + "!"
		// return next()
	}
	return next()
}
