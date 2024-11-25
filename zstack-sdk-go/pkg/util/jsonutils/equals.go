// Copyright (c) ZStack.io, Inc.

package jsonutils

import (
	"zstack.io/zstack-sdk-go/pkg/util/sortedmap"
)

func (dict *JSONDict) Equals(json JSONObject) bool {
	dict2, ok := json.(*JSONDict)
	if !ok {
		return false
	}
	if len(dict.data) != len(dict2.data) {
		return false
	}
	aNoB, aB, bA, bNoA := sortedmap.Split(dict.data, dict2.data)
	if len(aNoB) > 0 || len(bNoA) > 0 {
		return false
	}
	for _, k := range aB.Keys() {
		aVal, _ := aB.Get(k)
		bVal, _ := bA.Get(k)
		aJson := aVal.(JSONObject)
		bJson := bVal.(JSONObject)
		if !aJson.Equals(bJson) {
			return false
		}
	}
	return true
}

func (arr *JSONArray) Equals(json JSONObject) bool {
	arr2, ok := json.(*JSONArray)
	if !ok {
		return false
	}
	if len(arr.data) != len(arr2.data) {
		return false
	}
	for i, v := range arr.data {
		if !v.Equals(arr2.data[i]) {
			return false
		}
	}
	return true
}

func (o *JSONString) Equals(json JSONObject) bool {
	o2, ok := json.(*JSONString)
	if !ok {
		return false
	}
	return o.data == o2.data
}

func (o *JSONInt) Equals(json JSONObject) bool {
	o2, ok := json.(*JSONInt)
	if !ok {
		return false
	}
	return o.data == o2.data
}

func (o *JSONFloat) Equals(json JSONObject) bool {
	o2, ok := json.(*JSONFloat)
	if !ok {
		return false
	}
	return o.data == o2.data
}

func (o *JSONBool) Equals(json JSONObject) bool {
	o2, ok := json.(*JSONBool)
	if !ok {
		return false
	}
	return o.data == o2.data
}

func (o *JSONValue) Equals(json JSONObject) bool {
	_, ok := json.(*JSONValue)
	return ok
}
