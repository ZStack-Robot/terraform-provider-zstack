// Copyright (c) ZStack.io, Inc.

package jsonutils

import (
	"fmt"
	"reflect"
	"time"

	"github.com/kataras/golog"

	"zstack.io/zstack-sdk-go/pkg/util/gotypes"
	"zstack.io/zstack-sdk-go/pkg/util/reflectutils"
	"zstack.io/zstack-sdk-go/pkg/util/timeutils"
	"zstack.io/zstack-sdk-go/pkg/util/tristate"
)

func marshalSlice(val reflect.Value, info *reflectutils.StructFieldInfo, omitEmpty bool) JSONObject {
	if val.Kind() == reflect.Slice && val.IsNil() {
		if !omitEmpty {
			return JSONNull
		} else {
			return nil
		}
	}
	if val.Len() == 0 && info != nil && info.OmitEmpty && omitEmpty {
		return nil
	}
	objs := make([]JSONObject, 0)
	for i := 0; i < val.Len(); i += 1 {
		val := marshalValue(val.Index(i), nil, omitEmpty)
		if val != nil {
			objs = append(objs, val)
		}
	}
	arr := NewArray(objs...)
	if info != nil && info.ForceString {
		return NewString(arr.String())
	} else {
		return arr
	}
}

func marshalMap(val reflect.Value, info *reflectutils.StructFieldInfo, omitEmpty bool) JSONObject {
	if val.IsNil() {
		if !omitEmpty {
			return JSONNull
		} else {
			return nil
		}
	}
	keys := val.MapKeys()
	if len(keys) == 0 && info != nil && info.OmitEmpty && omitEmpty {
		return nil
	}
	objPairs := make([]JSONPair, 0)
	for i := 0; i < len(keys); i += 1 {
		key := keys[i]
		val := marshalValue(val.MapIndex(key), nil, omitEmpty)
		if val != nil {
			objPairs = append(objPairs, JSONPair{key: fmt.Sprintf("%s", key), val: val})
		}
	}
	dict := NewDict(objPairs...)
	if info != nil && info.ForceString {
		return NewString(dict.String())
	} else {
		return dict
	}
}

func marshalStruct(val reflect.Value, info *reflectutils.StructFieldInfo, omitEmpty bool) JSONObject {
	objPairs := struct2JSONPairs(val, omitEmpty)
	if len(objPairs) == 0 && info != nil && info.OmitEmpty && omitEmpty {
		return nil
	}
	dict := NewDict(objPairs...)
	if info != nil && info.ForceString {
		return NewString(dict.String())
	} else {
		return dict
	}
}

func findValueByKey(pairs []JSONPair, key string) JSONObject {
	for i := range pairs {
		if pairs[i].key == key {
			return pairs[i].val
		}
	}
	return nil
}

func struct2JSONPairs(val reflect.Value, omitEmpty bool) []JSONPair {
	fields := reflectutils.FetchStructFieldValueSet(val)
	objPairs := make([]JSONPair, 0, len(fields))
	depFields := make(map[string]string)
	for i := 0; i < len(fields); i += 1 {
		jsonInfo := fields[i].Info
		if jsonInfo.Ignore {
			continue
		}
		key := jsonInfo.MarshalName()
		if deprecatedBy, ok := fields[i].Info.Tags[TagDeprecatedBy]; ok {
			depFields[key] = deprecatedBy
			continue
		}
		val := marshalValue(fields[i].Value, jsonInfo, omitEmpty)
		if val != nil {
			objPair := JSONPair{key: key, val: val}
			objPairs = append(objPairs, objPair)
		}
	}
	depPairs := make([]JSONPair, 0, len(depFields))
	for depKey, key := range depFields {
		findLoop := false
		for {
			if okey, ok := depFields[key]; ok {
				if okey == depKey {
					// loop detected
					findLoop = true
					break
				}
				key = okey
			} else {
				break
			}
		}
		if findLoop {
			continue
		}
		val := findValueByKey(objPairs, key)
		if val != nil {
			objPair := JSONPair{key: depKey, val: val}
			depPairs = append(depPairs, objPair)
		}
	}
	objPairs = append(objPairs, depPairs...)
	return objPairs
}

func marshalInt64(val int64, info *reflectutils.StructFieldInfo, omitEmpty bool) JSONObject {
	if val == 0 && info != nil && info.OmitZero && omitEmpty {
		return nil
	} else if info != nil && info.ForceString {
		return NewString(fmt.Sprintf("%d", val))
	} else {
		return NewInt(val)
	}
}

func marshalFloat64(val float64, info *reflectutils.StructFieldInfo, bit int, omitEmpty bool) JSONObject {
	if val == 0.0 && info != nil && info.OmitZero && omitEmpty {
		return nil
	} else if info != nil && info.ForceString {
		return NewString(fmt.Sprintf("%f", val))
	} else {
		return NewFloat64(val)
	}
}

func marshalFloat32(val float32, info *reflectutils.StructFieldInfo, bit int, omitEmpty bool) JSONObject {
	if val == 0.0 && info != nil && info.OmitZero && omitEmpty {
		return nil
	} else if info != nil && info.ForceString {
		return NewString(fmt.Sprintf("%f", val))
	} else {
		return NewFloat32(val)
	}
}

func marshalBoolean(val bool, info *reflectutils.StructFieldInfo, omitEmpty bool) JSONObject {
	if !val && info != nil && info.OmitFalse && omitEmpty {
		return nil
	} else if info != nil && info.ForceString {
		return NewString(fmt.Sprintf("%v", val))
	} else {
		if val {
			return JSONTrue
		} else {
			return JSONFalse
		}
	}
}

func marshalTristate(val tristate.TriState, info *reflectutils.StructFieldInfo, omitEmpty bool) JSONObject {
	if val.IsTrue() {
		return JSONTrue
	} else if val.IsFalse() {
		return JSONFalse
	} else {
		if omitEmpty {
			return nil
		} else {
			return JSONNull
		}
	}
}

func marshalString(val string, info *reflectutils.StructFieldInfo, omitEmpty bool) JSONObject {
	if len(val) == 0 && info != nil && info.OmitEmpty && omitEmpty {
		return nil
	} else {
		return NewString(val)
	}
}

func marshalTime(val time.Time, info *reflectutils.StructFieldInfo, omitEmpty bool) JSONObject {
	if val.IsZero() {
		if info != nil && info.OmitEmpty && omitEmpty {
			return nil
		}
		return NewString("")
	} else {
		return NewString(timeutils.FullIsoTime(val))
	}
}

func Marshal(obj interface{}) JSONObject {
	if obj == nil {
		return JSONNull
	}
	val := reflect.ValueOf(obj)
	if kind := val.Kind(); val.IsZero() && kind == reflect.Ptr {
		return JSONNull
	}
	objValue := reflect.Indirect(val)
	mval := marshalValue(objValue, nil, true)
	if mval == nil {
		return JSONNull
	}
	return mval
}

func marshalValue(objValue reflect.Value, info *reflectutils.StructFieldInfo, omitEmpty bool) JSONObject {
	switch objValue.Type() {
	case JSONDictPtrType, JSONArrayPtrType, JSONBoolPtrType, JSONIntPtrType, JSONFloatPtrType, JSONStringPtrType, JSONObjectType:
		if objValue.IsNil() {
			if omitEmpty {
				return nil
			} else {
				return JSONNull
			}
		}
		return objValue.Interface().(JSONObject)
	case JSONDictType:
		json, ok := objValue.Interface().(JSONDict)
		if ok {
			if len(json.data) == 0 && info != nil && info.OmitEmpty && omitEmpty {
				return nil
			} else {
				return &json
			}
		} else {
			return nil
		}
	case JSONArrayType:
		json, ok := objValue.Interface().(JSONArray)
		if ok {
			if len(json.data) == 0 && info != nil && info.OmitEmpty && omitEmpty {
				return nil
			} else {
				return &json
			}
		} else {
			return nil
		}
	case JSONBoolType:
		json, ok := objValue.Interface().(JSONBool)
		if ok {
			if !json.data && info != nil && info.OmitEmpty && omitEmpty {
				return nil
			} else {
				return &json
			}
		} else {
			return nil
		}
	case JSONIntType:
		json, ok := objValue.Interface().(JSONInt)
		if ok {
			if json.data == 0 && info != nil && info.OmitEmpty && omitEmpty {
				return nil
			} else {
				return &json
			}
		} else {
			return nil
		}
	case JSONFloatType:
		json, ok := objValue.Interface().(JSONFloat)
		if ok {
			if json.data == 0.0 && info != nil && info.OmitEmpty && omitEmpty {
				return nil
			} else {
				return &json
			}
		} else {
			return nil
		}
	case JSONStringType:
		json, ok := objValue.Interface().(JSONString)
		if ok {
			if len(json.data) == 0 && info != nil && info.OmitEmpty && omitEmpty {
				return nil
			} else {
				return &json
			}
		} else {
			return nil
		}
	case tristate.TriStateType:
		tri, ok := objValue.Interface().(tristate.TriState)
		if ok {
			return marshalTristate(tri, info, omitEmpty)
		} else {
			return nil
		}
	}
	switch objValue.Kind() {
	case reflect.Slice, reflect.Array:
		return marshalSlice(objValue, info, omitEmpty)
	case reflect.Struct:
		if objValue.Type() == gotypes.TimeType {
			return marshalTime(objValue.Interface().(time.Time), info, omitEmpty)
		} else {
			return marshalStruct(objValue, info, omitEmpty)
		}
	case reflect.Map:
		return marshalMap(objValue, info, omitEmpty)
	case reflect.String:
		strValue := objValue.Convert(gotypes.StringType)
		return marshalString(strValue.Interface().(string), info, omitEmpty)
	case reflect.Bool:
		return marshalBoolean(objValue.Interface().(bool), info, omitEmpty)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		intValue := objValue.Convert(gotypes.Int64Type)
		return marshalInt64(intValue.Interface().(int64), info, omitEmpty)
	case reflect.Float32:
		floatVal := objValue.Convert(gotypes.Float32Type)
		return marshalFloat32(floatVal.Interface().(float32), info, 32, omitEmpty)
	case reflect.Float64:
		floatVal := objValue.Convert(gotypes.Float64Type)
		return marshalFloat64(floatVal.Interface().(float64), info, 64, omitEmpty)
	case reflect.Interface, reflect.Ptr:
		if objValue.IsNil() {
			if omitEmpty {
				return nil
			} else {
				return JSONNull
			}
		}
		return marshalValue(objValue.Elem(), info, false)
	default:
		golog.Errorf("unsupport object %s %s", objValue.Type(), objValue.Interface())
		return JSONNull
	}
}

func MarshalAll(obj interface{}) JSONObject {
	if obj == nil {
		return JSONNull
	}
	val := reflect.ValueOf(obj)
	if kind := val.Kind(); val.IsZero() && kind == reflect.Ptr {
		return JSONNull
	}
	objValue := reflect.Indirect(val)
	mval := marshalValue(objValue, nil, false)
	if mval == nil {
		return JSONNull
	}
	return mval
}
