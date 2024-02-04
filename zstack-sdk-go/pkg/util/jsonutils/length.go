package jsonutils

func (th *JSONString) Length() int {
	return len(th.data)
}

func (th *JSONDict) Length() int {
	return len(th.data)
}

func (th *JSONArray) Length() int {
	return len(th.data)
}
