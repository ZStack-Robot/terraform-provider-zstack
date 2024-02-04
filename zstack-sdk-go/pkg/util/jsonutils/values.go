package jsonutils

// Deprecated
func (th *JSONString) Value() string {
	return th.data
}

// Deprecated
func (th *JSONInt) Value() int64 {
	return th.data
}

// Deprecated
func (th *JSONFloat) Value() float64 {
	return th.data
}

// Deprecated
func (th *JSONBool) Value() bool {
	return th.data
}

// Deprecated
func (th *JSONDict) Value() map[string]JSONObject {
	mapJson, _ := th.GetMap()
	return mapJson
}

// Deprecated
func (th *JSONArray) Value() []JSONObject {
	arrJson, _ := th.GetArray()
	return arrJson
}
