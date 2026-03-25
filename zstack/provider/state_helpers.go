// Copyright (c) ZStack.io, Inc.

package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

func copyStringValues(values []types.String) []types.String {
	if len(values) == 0 {
		return nil
	}

	copied := make([]types.String, len(values))
	copy(copied, values)
	return copied
}

func stringValueOrNull(value string) types.String {
	if value == "" {
		return types.StringNull()
	}

	return types.StringValue(value)
}

func terraformStringsToSlice(values []types.String) []string {
	if len(values) == 0 {
		return nil
	}

	result := make([]string, 0, len(values))
	for _, value := range values {
		if value.IsNull() || value.ValueString() == "" {
			continue
		}
		result = append(result, value.ValueString())
	}

	return result
}

func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func intPtr(i int) *int {
	return &i
}

func int64Ptr(i int64) *int64 {
	return &i
}

func listToStringSlice(list types.List) []string {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}
	result := make([]string, 0, len(list.Elements()))
	for _, elem := range list.Elements() {
		if s, ok := elem.(types.String); ok && !s.IsNull() && s.ValueString() != "" {
			result = append(result, s.ValueString())
		}
	}
	return result
}
