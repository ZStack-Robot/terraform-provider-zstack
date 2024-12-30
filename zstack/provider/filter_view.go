// Copyright (c) ZStack.io, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type Filter struct {
	Name   types.String `tfsdk:"name"`
	Values types.Set    `tfsdk:"values"`
}
