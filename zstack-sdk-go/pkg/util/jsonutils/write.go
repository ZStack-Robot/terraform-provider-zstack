// Copyright (c) ZStack.io, Inc.

package jsonutils

import (
	"strings"

	"zstack.io/zstack-sdk-go/pkg/util/sortedmap"
)

type writeSource interface {
	buildString(sb *strings.Builder)
}

func (th *JSONString) buildString(sb *strings.Builder) {
	sb.WriteString(th.String())
}

func (th *JSONValue) buildString(sb *strings.Builder) {
	sb.WriteString(th.String())
}

func (th *JSONInt) buildString(sb *strings.Builder) {
	sb.WriteString(th.String())
}

func (th *JSONFloat) buildString(sb *strings.Builder) {
	sb.WriteString(th.String())
}

func (th *JSONBool) buildString(sb *strings.Builder) {
	sb.WriteString(th.String())
}

func (th *JSONDict) buildString(sb *strings.Builder) {
	sb.WriteByte('{')
	var idx = 0
	for iter := sortedmap.NewIterator(th.data); iter.HasMore(); iter.Next() {
		k, vinf := iter.Get()
		v := vinf.(JSONObject)
		if idx > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(quoteString(k))
		sb.WriteByte(':')

		v.buildString(sb)
		idx++
	}
	sb.WriteByte('}')
}

func (th *JSONArray) buildString(sb *strings.Builder) {
	sb.WriteByte('[')
	for idx, v := range th.data {
		if idx > 0 {
			sb.WriteByte(',')
		}
		v.buildString(sb)
	}
	sb.WriteByte(']')
}
