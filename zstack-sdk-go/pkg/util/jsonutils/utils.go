// Copyright (c) ZStack.io, Inc.

package jsonutils

import (
	"fmt"
	"time"

	"zstack.io/zstack-sdk-go/pkg/errors"
)

func NewStringArray(arr []string) *JSONArray {
	ret := NewArray()
	for _, a := range arr {
		ret.Add(NewString(a))
	}
	return ret
}

func (th *JSONArray) GetStringArray() []string {
	ret := make([]string, len(th.data))
	for i, obj := range th.data {
		s, e := obj.GetString()
		if e == nil {
			ret[i] = s
		}
	}
	return ret
}

func JSONArray2StringArray(arr []JSONObject) []string {
	ret := make([]string, len(arr))
	for i, o := range arr {
		s, e := o.GetString()
		if e == nil {
			ret[i] = s
		}
	}
	return ret
}

func GetStringArray(o JSONObject, key ...string) ([]string, error) {
	arr, err := o.GetArray(key...)
	if err != nil {
		return nil, err
	}
	return JSONArray2StringArray(arr), nil
}

func NewTimeString(tm time.Time) *JSONString {
	return NewString(tm.UTC().Format("2006-01-02T15:04:05Z"))
}

func GetQueryStringArray(query JSONObject, key string) []string {
	if query == nil {
		return nil
	}
	var arr []string
	searchObj, _ := query.Get(key)
	if searchObj != nil {
		switch searchObjTmp := searchObj.(type) {
		case *JSONArray:
			searchArr := searchObjTmp
			arr = searchArr.GetStringArray()
		case *JSONString:
			searchText, _ := searchObjTmp.GetString()
			arr = []string{searchText}
		case *JSONDict:
			arr = make([]string, 0)
			idx := 0
			for {
				searchText, err := searchObj.GetString(fmt.Sprintf("%d", idx))
				if err != nil {
					break
				}
				arr = append(arr, searchText)
				idx += 1
			}
		}
	}
	return arr
}

func CheckRequiredFields(data JSONObject, fields []string) error {
	jsonMap, err := data.GetMap()
	if err != nil {
		return errors.Wrap(err, "data.GetMap") //fmt.Errorf("fail to convert input to map")
	}
	for _, f := range fields {
		jsonVal, ok := jsonMap[f]
		if !ok {
			return errors.Wrap(ErrMissingInputField, f)
		}
		if jsonVal == JSONNull {
			return errors.Wrap(ErrNilInputField, f)
		}
	}
	return nil
}

func GetAnyString(json JSONObject, keys []string) string {
	val, _ := GetAnyString2(json, keys)
	return val
}

func GetAnyString2(json JSONObject, keys []string) (string, string) {
	if json == nil {
		return "", ""
	}
	for _, key := range keys {
		val, _ := json.GetString(key)
		if len(val) > 0 {
			return val, key
		}
	}
	return "", ""
}

func GetArrayOfPrefix(json JSONObject, prefix string) []JSONObject {
	if json == nil {
		return nil
	}

	retArray := make([]JSONObject, 0)
	idx := 0
	for {
		obj, _ := json.Get(fmt.Sprintf("%s.%d", prefix, idx))
		if obj == nil || obj == JSONNull {
			break
		}
		retArray = append(retArray, obj)
		idx += 1
	}
	return retArray
}
