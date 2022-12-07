package mongolia

import (
	"reflect"
)

//Returns a map of request tags for each field in struct
//e.g. GetStructTags(UserModel, 'ref') will return all the 'ref' values in a structs tags
func GetStructTags(model interface{}, tagName string) map[string]string {
	result := make(map[string]string)

	// TypeOf returns the reflection Type that represents the dynamic type of variable.
	// If variable is a nil interface value, TypeOf returns nil.
	t := reflect.TypeOf(model)
	if t.Kind() != reflect.Struct {
		return result
	}

	// Iterate over all available fields and read the tag value
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		reference := field.Tag.Get(tagName)
		if reference != "" {
			result[field.Name] = reference
		}
		//Also available: Data Type = field.Type.Name()
	}

	return result
}
