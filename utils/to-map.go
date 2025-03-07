package utils

import "reflect"

// use for put (effect 0 value)
// ToMap takes in an interface{} and converts it into a map[string]interface{} representation.
func ToMap(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	reflectValue := reflect.ValueOf(data)
	reflectType := reflect.TypeOf(data)

	for i := 0; i < reflectValue.NumField(); i++ {
		fieldValue := reflectValue.Field(i)
		fieldType := reflectType.Field(i)

		if fieldValue.Kind() != reflect.Struct {
			result[fieldType.Name] = fieldValue.Interface()
		}

		if fieldValue.Kind() == reflect.Struct {
			nestedFields := ToMap(fieldValue.Interface())
			for nestedKey, nestedValue := range nestedFields {
				result[nestedKey] = nestedValue
			}
		}
	}
	return result
}
