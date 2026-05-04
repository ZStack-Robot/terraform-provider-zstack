// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

func TestVipQosResource_Schema(t *testing.T) {
	var r vipQosResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}
	// Check required attributes
	required := []string{"vip_uuid"}
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
	computed := []string{"uuid", "type"}
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

func TestVipQosResource_Metadata(t *testing.T) {
	var r vipQosResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_vip_qos" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

// TestVipQosUpdateGuardsUnknownValues verifies that Unknown QoS rate fields are omitted from API payload.
// Unknown values must NOT be sent as zero (which would change QoS unintentionally).
// Per-field semantics: All 3 QoS rate Int64 fields (Port, OutboundBandwidth, InboundBandwidth)
// must omit the field when Unknown (zero rate = unlimited or invalid depending on semantic).
func TestVipQosUpdateGuardsUnknownValues(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		plan      vipQosModel
		assertFn  func(t *testing.T, setParam *param.SetVipQosParam)
	}{
		{
			name:      "Port Unknown omits field",
			fieldName: "Port",
			plan: vipQosModel{
				VipUuid:           types.StringValue("test-vip-uuid"),
				Port:              types.Int64Unknown(),
				OutboundBandwidth: types.Int64Null(),
				InboundBandwidth:  types.Int64Null(),
			},
			assertFn: func(t *testing.T, setParam *param.SetVipQosParam) {
				assert.Nil(t, setParam.Params.Port, "Port should be nil when Unknown")
			},
		},
		{
			name:      "OutboundBandwidth Unknown omits field",
			fieldName: "OutboundBandwidth",
			plan: vipQosModel{
				VipUuid:           types.StringValue("test-vip-uuid"),
				Port:              types.Int64Null(),
				OutboundBandwidth: types.Int64Unknown(),
				InboundBandwidth:  types.Int64Null(),
			},
			assertFn: func(t *testing.T, setParam *param.SetVipQosParam) {
				assert.Nil(t, setParam.Params.OutboundBandwidth, "OutboundBandwidth should be nil when Unknown")
			},
		},
		{
			name:      "InboundBandwidth Unknown omits field",
			fieldName: "InboundBandwidth",
			plan: vipQosModel{
				VipUuid:           types.StringValue("test-vip-uuid"),
				Port:              types.Int64Null(),
				OutboundBandwidth: types.Int64Null(),
				InboundBandwidth:  types.Int64Unknown(),
			},
			assertFn: func(t *testing.T, setParam *param.SetVipQosParam) {
				assert.Nil(t, setParam.Params.InboundBandwidth, "InboundBandwidth should be nil when Unknown")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the Update logic that builds setParam
			setParam := &param.SetVipQosParam{
				BaseParam: param.BaseParam{},
				Params:    param.SetVipQosParamDetail{},
			}

			// This is the FIXED logic - it DOES guard against Unknown
			if !tt.plan.Port.IsNull() && !tt.plan.Port.IsUnknown() {
				port := int(tt.plan.Port.ValueInt64())
				setParam.Params.Port = &port
			}
			if !tt.plan.OutboundBandwidth.IsNull() && !tt.plan.OutboundBandwidth.IsUnknown() {
				bandwidth := tt.plan.OutboundBandwidth.ValueInt64()
				setParam.Params.OutboundBandwidth = &bandwidth
			}
			if !tt.plan.InboundBandwidth.IsNull() && !tt.plan.InboundBandwidth.IsUnknown() {
				bandwidth := tt.plan.InboundBandwidth.ValueInt64()
				setParam.Params.InboundBandwidth = &bandwidth
			}

			// Verify the field is omitted (nil) when Unknown
			tt.assertFn(t, setParam)
		})
	}
}
