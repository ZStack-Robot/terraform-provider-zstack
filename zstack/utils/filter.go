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

			field := fieldByAPIName(resourceValue, apiFieldName)
			if !field.IsValid() {
				diags.AddError(
					"Invalid Filter Key",
					fmt.Sprintf("Field '%s' does not exist in resource", key),
				)
				return nil, diags
			}

			var fieldValue string
			switch field.Kind() {
			case reflect.Struct:
				if field.Type() == reflect.TypeOf(types.String{}) {
					strValue := field.Interface().(types.String)
					fieldValue = strValue.ValueString()
				} else {
					diags.AddError(
						"Unsupported Field Type",
						fmt.Sprintf("Field '%s' has unsupported type: %s", key, field.Type()),
					)
					return nil, diags
				}
			case reflect.String:
				fieldValue = field.String()
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				switch key {
				case "memory_size":
					fieldValue = fmt.Sprintf("%d", BytesToMB(field.Int()))
				case "disk_size", "volume_size":
					fieldValue = fmt.Sprintf("%d", BytesToGB(field.Int()))
				default:
					fieldValue = fmt.Sprintf("%d", field.Int())
				}
			case reflect.Bool:
				fieldValue = fmt.Sprintf("%t", field.Bool())
			default:
				diags.AddError(
					"Unsupported Field Type",
					fmt.Sprintf("Field '%s' has unsupported type: %s", key, field.Kind()),
				)
				return nil, diags
			}

			// Check if the field value matches any of the filter values
			valueMatch := false
			for _, value := range values {
				if fieldValue == value {
					valueMatch = true
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

func fieldByAPIName(resourceValue reflect.Value, apiFieldName string) reflect.Value {
	resourceValue = indirectValue(resourceValue)
	if !resourceValue.IsValid() || resourceValue.Kind() != reflect.Struct {
		return reflect.Value{}
	}

	if field := resourceValue.FieldByName(apiFieldName); field.IsValid() {
		return field
	}

	if field := fieldByTag(resourceValue, "json", apiFieldName); field.IsValid() {
		return field
	}

	if field := fieldByTag(resourceValue, "tfsdk", apiFieldName); field.IsValid() {
		return field
	}

	fieldName := exportedFieldName(apiFieldName)
	if field := resourceValue.FieldByName(fieldName); field.IsValid() {
		return field
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
