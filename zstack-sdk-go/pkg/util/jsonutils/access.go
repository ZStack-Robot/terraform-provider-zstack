// Copyright (c) ZStack.io, Inc.

package jsonutils

import (
	"strconv"
	"strings"
	"time"

	"zstack.io/zstack-sdk-go/pkg/errors"
	"zstack.io/zstack-sdk-go/pkg/util/sortedmap"
	"zstack.io/zstack-sdk-go/pkg/util/timeutils"
)

type JSONPair struct {
	key string
	val JSONObject
}

func NewDict(objs ...JSONPair) *JSONDict {
	dict := JSONDict{data: sortedmap.NewSortedMapWithCapa(len(objs))}
	for _, o := range objs {
		dict.Set(o.key, o.val)
	}
	return &dict
}

func NewArray(objs ...JSONObject) *JSONArray {
	arr := JSONArray{data: make([]JSONObject, 0, len(objs))}
	arr.data = append(arr.data, objs...)
	return &arr
}

func NewString(val string) *JSONString {
	return &JSONString{data: val}
}

func NewInt(val int64) *JSONInt {
	return &JSONInt{data: val}
}

// deprecated
func NewFloat(val float64) *JSONFloat {
	return &JSONFloat{data: val, bit: 64}
}

func NewFloat64(val float64) *JSONFloat {
	return &JSONFloat{data: val, bit: 64}
}

func NewFloat32(val float32) *JSONFloat {
	return &JSONFloat{data: float64(val), bit: 32}
}

func NewBool(val bool) *JSONBool {
	if val {
		return JSONTrue
	}
	return JSONFalse
}

func (th *JSONDict) Set(key string, value JSONObject) {
	th.data = sortedmap.Add(th.data, key, value)
}

func (th *JSONDict) Remove(key string) bool {
	return th.remove(key, true)
}

func (th *JSONDict) RemoveIgnoreCase(key string) bool {
	someRemoved := false
	for {
		removed := th.remove(key, false)
		if !removed {
			break
		}
		if !someRemoved {
			someRemoved = true
		}
	}
	return someRemoved
}

func (th *JSONDict) remove(key string, caseSensitive bool) bool {
	exist := false
	if !caseSensitive {
		th.data, _, exist = sortedmap.DeleteIgnoreCase(th.data, key)
	} else {
		th.data, exist = sortedmap.Delete(th.data, key)
	}
	return exist
}

func (th *JSONDict) Add(o JSONObject, keys ...string) error {
	obj := th
	for i := 0; i < len(keys); i++ {
		if i == len(keys)-1 {
			obj.Set(keys[i], o)
		} else {
			o, ok := obj.data.Get(keys[i])
			if !ok || o == JSONNull {
				obj.Set(keys[i], NewDict())
				o, ok = obj.data.Get(keys[i])
			}
			if ok {
				obj, ok = o.(*JSONDict)
				if !ok {
					return ErrInvalidJsonDict
				}
			} else {
				return ErrJsonDictFailInsert
			}
		}
	}
	return nil
}

func (th *JSONArray) SetAt(idx int, obj JSONObject) {
	th.data[idx] = obj
}

func (th *JSONArray) Add(objs ...JSONObject) {
	th.data = append(th.data, objs...)
}

func (th *JSONValue) Contains(keys ...string) bool {
	return false
}

func (th *JSONValue) ContainsIgnoreCases(keys ...string) bool {
	return false
}

func (th *JSONValue) Get(keys ...string) (JSONObject, error) {
	return nil, ErrUnsupported
}

func (th *JSONValue) GetIgnoreCases(keys ...string) (JSONObject, error) {
	return nil, ErrUnsupported
}

func (th *JSONValue) GetString(keys ...string) (string, error) {
	if len(keys) > 0 {
		return "", ErrOutOfKeyRange
	}
	return th.String(), nil
}

func (th *JSONValue) GetAt(i int, keys ...string) (JSONObject, error) {
	return nil, ErrUnsupported
}

func (th *JSONValue) Int(keys ...string) (int64, error) {
	return 0, ErrUnsupported
}

func (th *JSONValue) Float(keys ...string) (float64, error) {
	return 0.0, ErrUnsupported
}

func (th *JSONValue) Bool(keys ...string) (bool, error) {
	return false, ErrUnsupported
}

func (th *JSONValue) GetMap(keys ...string) (map[string]JSONObject, error) {
	return nil, ErrUnsupported
}

func (th *JSONValue) GetArray(keys ...string) ([]JSONObject, error) {
	return nil, ErrUnsupported
}

func (th *JSONDict) Contains(keys ...string) bool {
	obj, err := th._get(true, keys)
	if err == nil && obj != nil {
		return true
	}
	return false
}

func (th *JSONDict) ContainsIgnoreCases(keys ...string) bool {
	obj, err := th._get(false, keys)
	if err == nil && obj != nil {
		return true
	}
	return false
}

func (th *JSONDict) Get(keys ...string) (JSONObject, error) {
	return th._get(true, keys)
}

func (th *JSONDict) GetIgnoreCases(keys ...string) (JSONObject, error) {
	return th._get(false, keys)
}

func (th *JSONDict) _get(caseSensitive bool, keys []string) (JSONObject, error) {
	if len(keys) == 0 {
		return th, nil
	}
	for i := 0; i < len(keys); i++ {
		key := keys[i]
		var val interface{}
		var ok bool
		if caseSensitive {
			val, ok = th.data.Get(key)
		} else {
			val, _, ok = th.data.GetIgnoreCase(key)
		}
		if ok {
			if i == len(keys)-1 {
				return val.(JSONObject), nil
			} else {
				th, ok = val.(*JSONDict)
				if !ok {
					return nil, ErrInvalidJsonDict
				}
			}
		} else {
			return nil, ErrJsonDictKeyNotFound
		}
	}
	return nil, ErrOutOfKeyRange
}

func (th *JSONDict) GetString(keys ...string) (string, error) {
	if len(keys) > 0 {
		obj, err := th.Get(keys...)
		if err != nil {
			return "", errors.Wrap(err, "Get")
		}
		return obj.GetString()
	} else {
		return th.String(), nil
	}
}

func (th *JSONDict) GetMap(keys ...string) (map[string]JSONObject, error) {
	obj, err := th.Get(keys...)
	if err != nil {
		return nil, errors.Wrap(err, "Get")
	}
	dict, ok := obj.(*JSONDict)
	if !ok {
		return nil, ErrInvalidJsonDict
	}
	// allocate a map to hold the return map
	ret := make(map[string]JSONObject, len(th.data))
	for iter := sortedmap.NewIterator(dict.data); iter.HasMore(); iter.Next() {
		k, v := iter.Get()
		ret[k] = v.(JSONObject)
	}
	return ret, nil
}

func (th *JSONArray) GetAt(i int, keys ...string) (JSONObject, error) {
	if len(keys) > 0 {
		return nil, ErrOutOfKeyRange
	}
	if i < 0 {
		i = len(th.data) + i
	}
	if i >= 0 && i < len(th.data) {
		return th.data[i], nil
	} else {
		return nil, ErrOutOfIndexRange
	}
}

func (th *JSONArray) GetString(keys ...string) (string, error) {
	if len(keys) > 0 {
		return "", ErrOutOfKeyRange
	}
	return th.String(), nil
}

func (th *JSONDict) GetAt(i int, keys ...string) (JSONObject, error) {
	obj, err := th.Get(keys...)
	if err != nil {
		return nil, errors.Wrap(err, "Get")
	}
	arr, ok := obj.(*JSONArray)
	if !ok {
		return nil, ErrInvalidJsonArray
	}
	return arr.GetAt(i)
}

func (th *JSONArray) GetArray(keys ...string) ([]JSONObject, error) {
	if len(keys) > 0 {
		return nil, ErrOutOfKeyRange
	}
	// Allocate a new array to hold the return array
	ret := make([]JSONObject, len(th.data))
	for i := range th.data {
		ret[i] = th.data[i]
	}
	return ret, nil
}

func (th *JSONDict) GetArray(keys ...string) ([]JSONObject, error) {
	obj, err := th.Get(keys...)
	if err != nil {
		return nil, errors.Wrap(err, "Get")
	}
	if _, ok := obj.(*JSONDict); ok {
		return nil, ErrInvalidJsonArray
	}
	return obj.GetArray()
}

func _getarray(obj JSONObject, keys ...string) ([]JSONObject, error) {
	if len(keys) > 0 {
		return nil, ErrOutOfKeyRange
	}
	return []JSONObject{obj}, nil
}

func (th *JSONString) GetArray(keys ...string) ([]JSONObject, error) {
	return _getarray(th, keys...)
}

func (th *JSONInt) GetArray(keys ...string) ([]JSONObject, error) {
	return _getarray(th, keys...)
}

func (th *JSONFloat) GetArray(keys ...string) ([]JSONObject, error) {
	return _getarray(th, keys...)
}

func (th *JSONBool) GetArray(keys ...string) ([]JSONObject, error) {
	return _getarray(th, keys...)
}

func (th *JSONInt) Int(keys ...string) (int64, error) {
	if len(keys) > 0 {
		return 0, ErrOutOfKeyRange
	}
	return th.data, nil
}

func (th *JSONString) Int(keys ...string) (int64, error) {
	if len(keys) > 0 {
		return 0, ErrOutOfKeyRange
	}
	val, e := strconv.ParseInt(th.data, 10, 64)
	if e != nil {
		return 0, ErrInvalidJsonInt
	} else {
		return val, nil
	}
}

func (th *JSONInt) GetString(keys ...string) (string, error) {
	if len(keys) > 0 {
		return "", ErrOutOfKeyRange
	}
	return th.String(), nil
}

func (th *JSONDict) Int(keys ...string) (int64, error) {
	obj, err := th.Get(keys...)
	if err != nil {
		return 0, errors.Wrap(err, "Get")
	}
	return obj.Int()
}

func (th *JSONFloat) Float(keys ...string) (float64, error) {
	if len(keys) > 0 {
		return 0.0, ErrOutOfKeyRange
	}
	return th.data, nil
}

func (th *JSONInt) Float(keys ...string) (float64, error) {
	if len(keys) > 0 {
		return 0.0, ErrOutOfKeyRange
	}
	return float64(th.data), nil
}

func (th *JSONString) Float(keys ...string) (float64, error) {
	if len(keys) > 0 {
		return 0.0, ErrOutOfKeyRange
	}
	val, err := strconv.ParseFloat(th.data, 64)
	if err != nil {
		return 0.0, ErrInvalidJsonFloat
	} else {
		return val, nil
	}
}

func (th *JSONFloat) GetString(keys ...string) (string, error) {
	if len(keys) > 0 {
		return "", ErrOutOfKeyRange
	}
	return th.String(), nil
}

func (th *JSONDict) Float(keys ...string) (float64, error) {
	obj, err := th.Get(keys...)
	if err != nil {
		return 0, errors.Wrap(err, "Get")
	}
	return obj.Float()
}

func (th *JSONBool) Bool(keys ...string) (bool, error) {
	if len(keys) > 0 {
		return false, ErrOutOfKeyRange
	}
	return th.data, nil
}

func (th *JSONString) Bool(keys ...string) (bool, error) {
	if len(keys) > 0 {
		return false, ErrOutOfKeyRange
	}
	if strings.EqualFold(th.data, "true") || strings.EqualFold(th.data, "on") || strings.EqualFold(th.data, "yes") || th.data == "1" {
		return true, nil
	} else if strings.EqualFold(th.data, "false") || strings.EqualFold(th.data, "off") || strings.EqualFold(th.data, "no") || th.data == "0" {
		return false, nil
	} else {
		return false, ErrInvalidJsonBoolean
	}
}

func (th *JSONBool) GetString(keys ...string) (string, error) {
	if len(keys) > 0 {
		return "", ErrOutOfKeyRange
	}
	return th.String(), nil
}

func (th *JSONDict) Bool(keys ...string) (bool, error) {
	obj, err := th.Get(keys...)
	if err != nil {
		return false, errors.Wrap(err, "Get")
	}
	return obj.Bool()
}

func (th *JSONValue) GetTime(keys ...string) (time.Time, error) {
	return time.Time{}, ErrUnsupported
}

func (th *JSONString) GetTime(keys ...string) (time.Time, error) {
	if len(keys) > 0 {
		return time.Time{}, ErrOutOfKeyRange
	}
	t, e := timeutils.ParseTimeStr(th.data)
	if e == nil {
		return t, nil
	}
	for _, tf := range []string{time.RFC3339, time.RFC1123, time.UnixDate,
		time.RFC822,
	} {
		t, e := time.Parse(tf, th.data)
		if e == nil {
			return t, nil
		}
	}
	return th.JSONValue.GetTime()
}

func (th *JSONString) GetString(keys ...string) (string, error) {
	if len(keys) > 0 {
		return "", ErrOutOfKeyRange
	}
	return th.data, nil
}

func (th *JSONDict) GetTime(keys ...string) (time.Time, error) {
	obj, err := th.Get(keys...)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "Get")
	}
	str, ok := obj.(*JSONString)
	if !ok {
		return time.Time{}, ErrInvalidJsonString
	}
	return str.GetTime()
}
