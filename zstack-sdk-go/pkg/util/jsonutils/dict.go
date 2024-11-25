// Copyright (c) ZStack.io, Inc.

package jsonutils

import (
	"zstack.io/zstack-sdk-go/pkg/util/sortedmap"
)

func (dict *JSONDict) Update(json JSONObject) {
	dict2, ok := json.(*JSONDict)
	if !ok {
		return
	}
	dict.data = sortedmap.Merge(dict.data, dict2.data)
}

func (dict *JSONDict) UpdateDefault(json JSONObject) {
	dict2, ok := json.(*JSONDict)
	if !ok {
		return
	}
	dict.data = sortedmap.Merge(dict2.data, dict.data)
}

func Diff(a, b *JSONDict) (aNoB, aDiffB, aAndB, bNoA *JSONDict) {
	aNoB = NewDict()
	aDiffB = NewDict()
	aAndB = NewDict()
	bNoA = NewDict()

	var aData, bData sortedmap.SortedMap
	aNoB.data, aData, bData, bNoA.data = sortedmap.Split(a.data, b.data)
	for _, k := range aData.Keys() {
		aVal, _ := aData.Get(k)
		bVal, _ := bData.Get(k)
		aJson := aVal.(JSONObject)
		bJson := bVal.(JSONObject)
		if !aJson.Equals(bJson) {
			aDiffB.Set(k, NewArray(aJson, bJson))
		} else {
			aAndB.Set(k, aJson)
		}
	}

	return
}
