// # Copyright (c) ZStack.io, Inc.

package utils

import "github.com/hashicorp/terraform-plugin-framework/types"

func TfInt64ToIntPointer(number types.Int64) *int {
	if number.IsNull() {
		return nil
	}
	intNumber := int(number.ValueInt64())
	return &intNumber
}

func TfInt64ToInt64Pointer(number types.Int64) *int64 {
	if number.IsNull() {
		return nil
	}
	int64Number := number.ValueInt64()
	return &int64Number
}

func MBToBytes(mb int64) int64 {
	return mb * 1024 * 1024
}

func BytesToMB(bytes int64) int64 {
	return bytes / (1024 * 1024)
}

func BytesToGB(bytes int64) int64 {
	return bytes / (1024 * 1024 * 1024)
}

func GBToBytes(gb int64) int64 {
	return gb * 1024 * 1024 * 1024
}
