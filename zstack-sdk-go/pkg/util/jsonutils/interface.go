// Copyright (c) ZStack.io, Inc.

package jsonutils

import (
	"zstack.io/zstack-sdk-go/pkg/util/sortedmap"
)

func (th *JSONValue) Interface() interface{} {
	return nil
}

func (th *JSONBool) Interface() interface{} {
	return th.data
}

func (th *JSONInt) Interface() interface{} {
	return th.data
}

func (th *JSONFloat) Interface() interface{} {
	return th.data
}

func (th *JSONString) Interface() interface{} {
	return th.data
}

func (th *JSONArray) Interface() interface{} {
	ret := make([]interface{}, len(th.data))
	for i := 0; i < len(th.data); i += 1 {
		ret[i] = th.data[i].Interface()
	}
	return ret
}

func (th *JSONDict) Interface() interface{} {
	mapping := make(map[string]interface{})

	for iter := sortedmap.NewIterator(th.data); iter.HasMore(); iter.Next() {
		k, v := iter.Get()
		mapping[k] = v.(JSONObject).Interface()
	}

	return mapping
}
