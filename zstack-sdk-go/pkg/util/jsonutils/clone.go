// Copyright (c) ZStack.io, Inc.

package jsonutils

import (
	"reflect"

	"zstack.io/zstack-sdk-go/pkg/util/sortedmap"
	"zstack.io/zstack-sdk-go/pkg/util/utils"
)

func (th *JSONDict) Copy(excludes ...string) *JSONDict {
	return th.CopyExcludes(excludes...)
}

func (th *JSONDict) CopyExcludes(excludes ...string) *JSONDict {
	dict := NewDict()
	for iter := sortedmap.NewIterator(th.data); iter.HasMore(); iter.Next() {
		k, v := iter.Get()
		exists, _ := utils.InStringArray(k, excludes)
		if !exists {
			dict.Set(k, v.(JSONObject))
		}
	}
	return dict
}

func (th *JSONDict) CopyIncludes(includes ...string) *JSONDict {
	dict := NewDict()
	for iter := sortedmap.NewIterator(th.data); iter.HasMore(); iter.Next() {
		k, v := iter.Get()
		exists, _ := utils.InStringArray(k, includes)
		if exists {
			dict.Set(k, v.(JSONObject))
		}
	}
	return dict
}

func (th *JSONArray) Copy() *JSONArray {
	arr := NewArray()
	arr.data = append(arr.data, th.data...)
	return arr
}

func (th *JSONString) DeepCopy() interface{} {
	return DeepCopy(th)
}

func (th *JSONInt) DeepCopy() interface{} {
	return DeepCopy(th)
}

func (th *JSONFloat) DeepCopy() interface{} {
	return DeepCopy(th)
}

func (th *JSONBool) DeepCopy() interface{} {
	return DeepCopy(th)
}

func (th *JSONArray) DeepCopy() interface{} {
	return DeepCopy(th)
}

func (th *JSONDict) DeepCopy() interface{} {
	return DeepCopy(th)
}

func DeepCopy(obj JSONObject) JSONObject {
	if obj == nil || reflect.ValueOf(obj).IsNil() {
		return nil
	}
	switch v := obj.(type) {
	case *JSONString:
		vc := *v
		return &vc
	case *JSONInt:
		vc := *v
		return &vc
	case *JSONFloat:
		vc := *v
		return &vc
	case *JSONBool:
		vc := *v
		return &vc
	case *JSONArray:
		elemsC := make([]JSONObject, v.Length())
		for i := 0; i < v.Length(); i++ {
			elem, _ := v.GetAt(i)
			elemC := DeepCopy(elem)
			elemsC[i] = elemC
		}
		vc := NewArray(elemsC...)
		return vc
	case *JSONDict:
		vc := NewDict()
		m, _ := v.GetMap()
		for mk, mv := range m {
			mvc := DeepCopy(mv)
			vc.Set(mk, mvc)
		}
		return vc
	case *JSONValue:
		vc := *v
		return &vc
	}
	return nil
}
