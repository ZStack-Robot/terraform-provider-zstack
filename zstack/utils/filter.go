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

			fieldName := strings.Title(apiFieldName)
			field := resourceValue.FieldByName(fieldName)

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
				if key == "memory_size" {
					fieldValue = fmt.Sprintf("%d", BytesToMB(field.Int()))
				} else if key == "disk_size" {
					fieldValue = fmt.Sprintf("%d", BytesToGB(field.Int()))
				} else if key == "volume_size" {
					fieldValue = fmt.Sprintf("%d", BytesToGB(field.Int()))
				} else {
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
