package sources

import (
	"reflect"
	"strings"
)

type SourceField struct {
	Id          string      `json:"id"`
	Type        string      `json:"type"`
	Label       string      `json:"label,omitempty"`
	Placeholder string      `json:"placeholder,omitempty"`
	Default     interface{} `json:"default,omitempty"`
	Required    bool        `json:"required,omitempty"`
	Description string      `json:"description,omitempty"`
}

func GenFields(data interface{}) (fields []*SourceField) {
	if data != nil {
		dataMeta := reflect.Indirect(reflect.ValueOf(data)).Type()

		for i := 0; i < dataMeta.NumField(); i += 1 {
			field := dataMeta.Field(i)
			if field.IsExported() {
				fields = append(fields, GenField(field))
			}
		}
	}
	return
}

func GenField(field reflect.StructField) (f *SourceField) {
	f = &SourceField{
		Id:          field.Name,
		Type:        field.Type.String(),
		Label:       field.Tag.Get("label"),
		Placeholder: field.Tag.Get("placeholder"),
		Default:     field.Tag.Get("default"),
		Description: field.Tag.Get("description"),
	}

	jsonTag := field.Tag.Get("json")
	jsonTuples := strings.Split(jsonTag, ",")
	if len(jsonTuples) > 0 && len(jsonTuples[0]) > 0 {
		f.Id = jsonTuples[0]
	}

	if f.Label == "" {
		f.Label = field.Name
	}

	return
}
