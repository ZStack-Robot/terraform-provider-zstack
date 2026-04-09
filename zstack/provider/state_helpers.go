// Copyright (c) ZStack.io, Inc.

package provider

import (
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	zstackerrors "github.com/zstackio/zstack-sdk-go-v2/pkg/errors"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

// ErrResourceNotFound is returned by finder functions when the requested
// resource does not exist. Callers should check this with errors.Is and
// call resp.State.RemoveResource when true.
var ErrResourceNotFound = errors.New("resource not found")

// findResourceByGet wraps a Get-style SDK method (returns *T, error) into a
// standard finder. It converts ZStack not-found errors into ErrResourceNotFound.
func findResourceByGet[T any](getFunc func(uuid string) (*T, error), uuid string) (*T, error) {
	result, err := getFunc(uuid)
	if err != nil {
		if isZStackNotFoundError(err) {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}
	return result, nil
}

// findResourceByQuery wraps a Query-style SDK method (returns []T, error) into
// a standard finder. It adds a uuid filter to the query and returns
// ErrResourceNotFound when the result set is empty.
func findResourceByQuery[T any](queryFunc func(params *param.QueryParam) ([]T, error), uuid string) (*T, error) {
	q := param.NewQueryParam()
	q.AddQ("uuid=" + uuid)
	results, err := queryFunc(&q)
	if err != nil {
		if isZStackNotFoundError(err) {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}
	if len(results) == 0 {
		return nil, ErrResourceNotFound
	}
	return &results[0], nil
}

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

func stringPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
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

func stringSliceToList(values []string) types.List {
	if len(values) == 0 {
		return types.ListValueMust(types.StringType, []attr.Value{})
	}
	elems := make([]attr.Value, len(values))
	for i, v := range values {
		elems[i] = types.StringValue(v)
	}
	return types.ListValueMust(types.StringType, elems)
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

func isZStackNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	if zstackerrors.Cause(err) == zstackerrors.ErrNotFound {
		return true
	}

	return strings.Contains(err.Error(), "status code 404")
}
