package mongolia

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"time"
	"strings"
)

func ValidateFindResults[T Model](findResults []T) error {
	for _, elem := range findResults {
		if err := elem.ValidateRead(); err != nil {
			return err
		}
	}
	return nil
}


func camelCaseString(input string) string {
	if len(input) < 1 {
		return input
	}
	output := strings.ToLower(string(input[0])) + input[1:len(input)]
	return output
}

func splitCommaTag(input string) string {
	// Split tag on commas, assume name is first
	all_tags := strings.Split(input, ",")
	name := all_tags[0]
	return name
}

func GetFieldName(field reflect.StructField) string {
	//Priority 1: bson tag
	bson_tag := field.Tag.Get("bson")
	if bson_tag != "" {
		return splitCommaTag(bson_tag)
	}

	//Priority 2: json tag
	json_tag := field.Tag.Get("json")
	if json_tag != "" {
		return splitCommaTag(json_tag)
	}

	//Priority 2: camelCase
	return camelCaseString(field.Name)
}

func RecurseTagReference(t reflect.Type, tagName string, rootName string, result map[string]string) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return
	}

	// Iterate over all available fields and read the tag value
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		ft := field.Type
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}

		if ft.Kind() == reflect.Struct {
			recurseName := rootName + GetFieldName(field) + "."
			RecurseTagReference(ft, tagName, recurseName, result)
		}
		if ft.Kind() == reflect.Slice && ft.Elem().Kind() == reflect.Struct {
			recurseName := rootName + GetFieldName(field) + "."
			RecurseTagReference(ft.Elem(), tagName, recurseName, result)
		}

		reference := field.Tag.Get(tagName)
		if reference != "" {
			refName := rootName + GetFieldName(field)
			result[refName] = reference
		}
	}
}

func GetStructTags(model interface{}, tagName string) map[string]string {
	result := make(map[string]string)

	// TypeOf returns the reflection Type that represents the dynamic type of variable.
	// If variable is a nil interface value, TypeOf returns nil.
	t := reflect.TypeOf(model)
	if t.Kind() != reflect.Struct {
		return result
	}

	// Recursively find tag values and append
	RecurseTagReference(t, tagName, "", result)

	return result
}

func BSONUpdateAtHook(update bson.D) {
	var setElements bson.D
	ok := false
	for i, elem := range update {
		if elem.Key == "$set" {
			//If this is a bson.D set
			if setElements, ok = elem.Value.(bson.D); ok {
				updateset := true
				for _, set := range setElements {
					if set.Key == "updatedAt" {
						updateset = false
					}
				}
				if updateset {
					setElements = append(setElements, bson.E{"updatedAt", time.Now().UTC()})
				}
				update[i].Value = setElements
			}
		}
	}
}

// Generic conversion to map
func CastToMap(update any) map[string]any {
	bsonmap := make(map[string]any)
	if update != nil {
		switch update.(type) {
		case bson.D:
			return CastBDToMap(update.(bson.D))
		case bson.M:
			return CastBMToMap(update.(bson.M))
		default:
			return bsonmap
		}
	}
	return bsonmap
}

// Convert bson.D Set to map for easy validation
func CastBDToMap(update bson.D) map[string]any {
	bsonmap := make(map[string]any)
	for _, elem := range update {
		if elem.Key == "$set" {
			if setElements, ok := elem.Value.(bson.D); ok {
				for _, set := range setElements {
					bsonmap[set.Key] = set.Value
				}
			} else if modelElements, ok := elem.Value.(Model); ok {
				//if its a full model set do this
				jsonmap := StructToMap(modelElements)
				return jsonmap
			}
		}
	}
	return bsonmap
}

// Convert bson.M Set to map for easy validation
func CastBMToMap(update bson.M) map[string]any {
	bsonmap := make(map[string]any)
	for key, value := range update {
		bsonmap[key] = value
	}
	return bsonmap
}

func CastMapToDB(update map[string]any) bson.D {
	set := bson.D{}
	for key, value := range update {
		set = append(set, bson.E{key, value})
	}
	return bson.D{{"$set", set}}
}

func StructToMap(model any) map[string]any {
	var model_map map[string]any
	model_json, _ := json.Marshal(model)
	json.Unmarshal(model_json, &model_map)
	return model_map
}
