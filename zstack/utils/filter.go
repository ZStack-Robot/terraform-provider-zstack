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
	filters map[string]string,
) ([]T, diag.Diagnostics) {
	var diags diag.Diagnostics
	var filteredResources []T

	for _, resource := range resources {
		match := true
		resourceValue := reflect.ValueOf(resource)

		for key, value := range filters {

			fieldName := strings.Title(key)
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
				fieldValue = fmt.Sprintf("%d", field.Int())
			case reflect.Bool:
				fieldValue = fmt.Sprintf("%t", field.Bool())
			default:
				diags.AddError(
					"Unsupported Field Type",
					fmt.Sprintf("Field '%s' has unsupported type: %s", key, field.Kind()),
				)
				return nil, diags
			}

			if fieldValue != value {
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
