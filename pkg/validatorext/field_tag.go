package validatorext

import (
	"reflect"
	"strings"
)

func FieldTag(fld reflect.StructField) string {
	name := fld.Tag.Get("json")
	if name == "-" {
		return ""
	}

	commaIdx := strings.Index(name, ",")
	if commaIdx > -1 {
		name = name[:commaIdx]
	}
	return name
}
