// Copyright (c) ZStack.io, Inc.

package jsonutils

import (
	"reflect"

	"zstack.io/zstack-sdk-go/pkg/errors"
	"zstack.io/zstack-sdk-go/pkg/util/gotypes"
)

var (
	JSONDictType      reflect.Type
	JSONArrayType     reflect.Type
	JSONStringType    reflect.Type
	JSONIntType       reflect.Type
	JSONFloatType     reflect.Type
	JSONBoolType      reflect.Type
	JSONDictPtrType   reflect.Type
	JSONArrayPtrType  reflect.Type
	JSONStringPtrType reflect.Type
	JSONIntPtrType    reflect.Type
	JSONFloatPtrType  reflect.Type
	JSONBoolPtrType   reflect.Type
	JSONObjectType    reflect.Type
)

func init() {
	JSONDictType = reflect.TypeOf(JSONDict{})
	JSONArrayType = reflect.TypeOf(JSONArray{})
	JSONStringType = reflect.TypeOf(JSONString{})
	JSONIntType = reflect.TypeOf(JSONInt{})
	JSONFloatType = reflect.TypeOf(JSONFloat{})
	JSONBoolType = reflect.TypeOf(JSONBool{})
	JSONDictPtrType = reflect.TypeOf(&JSONDict{})
	JSONArrayPtrType = reflect.TypeOf(&JSONArray{})
	JSONStringPtrType = reflect.TypeOf(&JSONString{})
	JSONIntPtrType = reflect.TypeOf(&JSONInt{})
	JSONFloatPtrType = reflect.TypeOf(&JSONFloat{})
	JSONBoolPtrType = reflect.TypeOf(&JSONBool{})
	JSONObjectType = reflect.TypeOf((*JSONObject)(nil)).Elem()

	gotypes.RegisterSerializable(JSONObjectType, func() gotypes.ISerializable {
		return nil
	})

	gotypes.RegisterSerializable(JSONDictPtrType, func() gotypes.ISerializable {
		return NewDict()
	})

	gotypes.RegisterSerializable(JSONArrayPtrType, func() gotypes.ISerializable {
		return NewArray()
	})
}

func JSONDeserialize(objType reflect.Type, strVal string) (gotypes.ISerializable, error) {
	objPtr, err := gotypes.NewSerializable(objType)
	if err != nil {
		return nil, errors.Wrap(err, "gotypes.NewSerializable")
	}
	json, err := ParseString(strVal)
	if err != nil {
		return nil, errors.Wrap(err, "ParseString")
	}
	if objPtr == nil {
		return json, nil
	}
	err = json.Unmarshal(objPtr)
	if err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}
	objPtr = gotypes.Transform(objType, objPtr)
	return objPtr, nil
}
