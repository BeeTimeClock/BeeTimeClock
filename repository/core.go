package repository

import "reflect"

func ToMap(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	reflectValue := reflect.ValueOf(data)

	if reflectValue.Kind() == reflect.Ptr && reflectValue.Elem().Kind() == reflect.Struct {
		reflectValue = reflectValue.Elem()
	}

	for i := 0; i < reflectValue.NumField(); i++ {
		fieldValue := reflectValue.Field(i)
		fieldType := reflectValue.Type().Field(i)

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
