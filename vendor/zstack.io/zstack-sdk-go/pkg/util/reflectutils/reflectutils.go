package reflectutils

import (
	"fmt"
	"reflect"

	"github.com/kataras/golog"
)

func FindStructFieldValue(dataValue reflect.Value, name string) (reflect.Value, bool) {
	set := FetchStructFieldValueSet(dataValue)
	val, find := set.GetValue(name)
	if find && val.CanSet() {
		return val, true
	}
	return reflect.Value{}, false
}

func FindStructFieldInterface(dataValue reflect.Value, name string) (interface{}, bool) {
	set := FetchStructFieldValueSet(dataValue)
	return set.GetInterface(name)
}

func FillEmbededStructValue(container reflect.Value, embed reflect.Value) bool {
	containerType := container.Type()
	for i := 0; i < containerType.NumField(); i += 1 {
		fieldType := containerType.Field(i)
		fieldValue := container.Field(i)
		if fieldType.Type.Kind() == reflect.Struct && fieldType.Anonymous {
			if fieldType.Type == embed.Type() {
				fieldValue.Set(embed)
				return true
			} else {
				filled := FillEmbededStructValue(fieldValue, embed)
				if filled {
					return true
				}
			}
		}

	}
	return false
}

func SetStructFieldValue(structValue reflect.Value, fieldName string, val reflect.Value) bool {
	set := FetchStructFieldValueSet(structValue)
	target, find := set.GetValue(fieldName)
	if !find {
		return false
	}
	if !target.CanSet() {
		return false
	}
	target.Set(val)
	return true
}

func ExpandInterface(val interface{}) []interface{} {
	value := reflect.Indirect(reflect.ValueOf(val))
	if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
		ret := make([]interface{}, value.Len())
		for i := 0; i < len(ret); i += 1 {
			ret[i] = value.Index(i).Interface()
		}
		return ret
	} else {
		return []interface{}{val}
	}
}

// tagetType must not be a pointer
func getAnonymouStructPointer(structValue reflect.Value, targetType reflect.Type) interface{} {
	structType := structValue.Type()
	if structType == targetType {
		return structValue.Addr().Interface()
	}
	for i := 0; i < structValue.NumField(); i += 1 {
		fieldType := structType.Field(i)
		if fieldType.Anonymous && fieldType.Type.Kind() == reflect.Struct {
			ptr := getAnonymouStructPointer(structValue.Field(i), targetType)
			if ptr != nil {
				return ptr
			}
		}
	}
	return nil
}

func FindAnonymouStructPointer(data interface{}, targetPtr interface{}) error {
	targetValue := reflect.ValueOf(targetPtr)
	if targetValue.Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer to pointer")
	}
	targetValue = targetValue.Elem()
	if targetValue.Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer to pointer")
	}
	targetType := targetValue.Type().Elem()
	if targetType.Kind() != reflect.Struct {
		return fmt.Errorf("target type must be a struct")
	}
	structValue := reflect.ValueOf(data)
	if structValue.Kind() != reflect.Ptr {
		return fmt.Errorf("data type must be a pointer to struct")
	}
	structValue = reflect.ValueOf(data).Elem()
	if structValue.Kind() != reflect.Struct {
		return fmt.Errorf("data type must be a pointer to struct")
	}
	ptr := getAnonymouStructPointer(structValue, targetType)
	if ptr == nil {
		return fmt.Errorf("no anonymous struct found")
	}
	targetValue.Set(reflect.ValueOf(ptr))
	return nil
}

func StructContains(type1 reflect.Type, type2 reflect.Type) bool {
	if type1.Kind() != reflect.Struct || type2.Kind() != reflect.Struct {
		golog.Errorf("types should be struct!")
		return false
	}
	if type1 == type2 {
		return true
	}
	for i := 0; i < type1.NumField(); i += 1 {
		field := type1.Field(i)
		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			contains := StructContains(field.Type, type2)
			if contains {
				return true
			}
		}
	}
	return false
}
