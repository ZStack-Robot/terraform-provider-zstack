// Copyright (c) ZStack.io, Inc.

package provider

import "strings"

const stackTemplateFormatVersionMarker = "ZStackTemplateFormatVersion"

func hasStackTemplateFormatVersionMarker(templateContent string) bool {
	return strings.Contains(templateContent, stackTemplateFormatVersionMarker)
}
