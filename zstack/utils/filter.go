// Copyright (c) ZStack.io, Inc.
// SPDX-License-Identifier: MPL-2.0

package utils

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FilterResource[T any](
	ctx context.Context,
	resources []T,
	filters map[string][]string,
	dataSourceName string,
) ([]T, diag.Diagnostics) {
	var diags diag.Diagnostics
	var filteredResources []T

	fieldMapping := GetFieldMapping(dataSourceName)

	for _, resource := range resources {
		match := true
		resourceValue := reflect.ValueOf(resource)

		for key, values := range filters {
			//  Terraform Schema map to API Attribute
			apiFieldName, ok := fieldMapping[key]
			if !ok {
				apiFieldName = key
			}

			fieldValues, found, unsupportedType := fieldValuesByAPIName(resourceValue, apiFieldName, key)
			if !found {
				diags.AddError(
					"Invalid Filter Key",
					fmt.Sprintf("Field '%s' does not exist in resource", key),
				)
				return nil, diags
			}
			if unsupportedType != "" {
				diags.AddError(
					"Unsupported Field Type",
					fmt.Sprintf("Field '%s' has unsupported type: %s", key, unsupportedType),
				)
				return nil, diags
			}

			// Check if the field value matches any of the filter values
			valueMatch := false
			for _, value := range values {
				for _, fieldValue := range fieldValues {
					if fieldValue == value {
						valueMatch = true
						break
					}
				}
				if valueMatch {
					break
				}
			}

			if !valueMatch {
				match = false
				break
			}
		}

		if match {
			filteredResources = append(filteredResources, resource)
		}
	}

	return filteredResources, diags
}

func fieldValuesByAPIName(resourceValue reflect.Value, apiFieldName string, key string) ([]string, bool, string) {
	if apiFieldName == "" {
		return nil, false, ""
	}

	return fieldValuesByAPIPath(resourceValue, strings.Split(apiFieldName, "."), key)
}

func fieldValuesByAPIPath(value reflect.Value, apiFieldPath []string, key string) ([]string, bool, string) {
	value = indirectValue(value)
	if !value.IsValid() {
		return nil, false, ""
	}

	if len(apiFieldPath) == 0 {
		values, unsupportedType := filterFieldValues(key, value)
		return values, true, unsupportedType
	}

	switch value.Kind() {
	case reflect.Struct:
		field := fieldByAPIName(value, apiFieldPath[0])
		if !field.IsValid() {
			return nil, false, ""
		}
		return fieldValuesByAPIPath(field, apiFieldPath[1:], key)
	case reflect.Slice, reflect.Array:
		values := make([]string, 0, value.Len())
		for i := 0; i < value.Len(); i++ {
			elemValues, found, unsupportedType := fieldValuesByAPIPath(value.Index(i), apiFieldPath, key)
			if unsupportedType != "" {
				return nil, true, unsupportedType
			}
			if found {
				values = append(values, elemValues...)
			}
		}
		return values, true, ""
	default:
		return nil, true, value.Kind().String()
	}
}

func filterFieldValues(key string, field reflect.Value) ([]string, string) {
	field = indirectValue(field)
	if !field.IsValid() {
		return nil, "invalid"
	}

	switch field.Kind() {
	case reflect.Struct:
		if field.Type() == reflect.TypeOf(types.String{}) {
			strValue := field.Interface().(types.String)
			return []string{strValue.ValueString()}, ""
		}
		return nil, field.Type().String()
	case reflect.String:
		return []string{field.String()}, ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch key {
		case "memory_size":
			return []string{fmt.Sprintf("%d", BytesToMB(field.Int()))}, ""
		case "disk_size", "volume_size":
			return []string{fmt.Sprintf("%d", BytesToGB(field.Int()))}, ""
		default:
			return []string{fmt.Sprintf("%d", field.Int())}, ""
		}
	case reflect.Bool:
		return []string{fmt.Sprintf("%t", field.Bool())}, ""
	case reflect.Slice, reflect.Array:
		values := make([]string, 0, field.Len())
		for i := 0; i < field.Len(); i++ {
			elemValues, unsupportedType := filterFieldValues(key, field.Index(i))
			if unsupportedType != "" {
				return nil, unsupportedType
			}
			values = append(values, elemValues...)
		}
		return values, ""
	default:
		return nil, field.Kind().String()
	}
}

func fieldByAPIName(resourceValue reflect.Value, apiFieldName string) reflect.Value {
	resourceValue = indirectValue(resourceValue)
	if !resourceValue.IsValid() || resourceValue.Kind() != reflect.Struct {
		return reflect.Value{}
	}

	for _, fieldName := range fieldNameCandidates(apiFieldName) {
		if field := resourceValue.FieldByName(fieldName); field.IsValid() {
			return field
		}

		if field := fieldByTag(resourceValue, "json", fieldName); field.IsValid() {
			return field
		}

		if field := fieldByTag(resourceValue, "tfsdk", fieldName); field.IsValid() {
			return field
		}

		exportedName := exportedFieldName(fieldName)
		if field := resourceValue.FieldByName(exportedName); field.IsValid() {
			return field
		}
	}

	return reflect.Value{}
}

func fieldByTag(resourceValue reflect.Value, tagKey string, tagValue string) reflect.Value {
	resourceValue = indirectValue(resourceValue)
	if !resourceValue.IsValid() || resourceValue.Kind() != reflect.Struct {
		return reflect.Value{}
	}

	resourceType := resourceValue.Type()
	for i := 0; i < resourceType.NumField(); i++ {
		structField := resourceType.Field(i)
		if structField.PkgPath != "" && !structField.Anonymous {
			continue
		}

		fieldTagValue := strings.Split(structField.Tag.Get(tagKey), ",")[0]
		if fieldTagValue != "" && fieldTagValue != "-" && fieldTagValue == tagValue {
			return resourceValue.Field(i)
		}

		if structField.Anonymous {
			if field := fieldByTag(resourceValue.Field(i), tagKey, tagValue); field.IsValid() {
				return field
			}
		}
	}

	return reflect.Value{}
}

func indirectValue(value reflect.Value) reflect.Value {
	for value.IsValid() && (value.Kind() == reflect.Interface || value.Kind() == reflect.Pointer) {
		if value.IsNil() {
			return reflect.Value{}
		}
		value = value.Elem()
	}

	return value
}

func exportedFieldName(name string) string {
	if name == "" {
		return ""
	}

	return strings.ToUpper(name[:1]) + name[1:]
}

func fieldNameCandidates(name string) []string {
	candidates := []string{name}
	if strings.Contains(name, "_") {
		candidates = append(candidates, snakeToLowerCamel(name), snakeToUpperCamel(name))
	}

	result := make([]string, 0, len(candidates))
	seen := make(map[string]struct{}, len(candidates))
	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		if _, ok := seen[candidate]; ok {
			continue
		}
		seen[candidate] = struct{}{}
		result = append(result, candidate)
	}
	return result
}

func snakeToLowerCamel(name string) string {
	upperCamel := snakeToUpperCamel(name)
	if upperCamel == "" {
		return ""
	}

	return strings.ToLower(upperCamel[:1]) + upperCamel[1:]
}

func snakeToUpperCamel(name string) string {
	parts := strings.Split(name, "_")
	var builder strings.Builder
	for _, part := range parts {
		if part == "" {
			continue
		}
		builder.WriteString(strings.ToUpper(part[:1]))
		if len(part) > 1 {
			builder.WriteString(part[1:])
		}
	}
	return builder.String()
}
