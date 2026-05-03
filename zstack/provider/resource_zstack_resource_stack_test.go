// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestResourceStackResource_Schema(t *testing.T) {
	var r resourceStackResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}
	// Check required attributes
	required := []string{"name"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}
	// Check computed attributes
	computed := []string{"uuid", "version", "status"}
	for _, attr := range computed {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing computed attribute %q", attr)
		}
		if !a.IsComputed() {
			t.Errorf("attribute %q should be computed", attr)
		}
	}
}

func TestResourceStackResource_Metadata(t *testing.T) {
	var r resourceStackResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_resource_stack" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestValidateResourceStackTemplateInputs(t *testing.T) {
	tests := []struct {
		name string
		plan resourceStackModel
		want bool
	}{
		{
			name: "template content with marker",
			plan: resourceStackModel{
				TemplateContent: types.StringValue(testStackTemplateContent),
			},
			want: true,
		},
		{
			name: "template uuid without inline content",
			plan: resourceStackModel{
				TemplateUuid: types.StringValue("test-template-uuid"),
			},
			want: true,
		},
		{
			name: "missing template content and uuid",
			plan: resourceStackModel{},
			want: false,
		},
		{
			name: "template content without marker",
			plan: resourceStackModel{
				TemplateContent: types.StringValue(`{"Resources":{}}`),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var diags diag.Diagnostics
			got := validateResourceStackTemplateInputs(tt.plan, &diags)
			if got != tt.want {
				t.Fatalf("validateResourceStackTemplateInputs() = %t, want %t", got, tt.want)
			}
			if got == diags.HasError() {
				t.Fatalf("diagnostics error state should be inverse of result, got result=%t diagnostics=%v", got, diags)
			}
		})
	}
}
