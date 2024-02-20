package jsonutils

func (th *JSONDict) Size() int {
	return len(th.data)
}

func (th *JSONArray) Size() int {
	return len(th.data)
}
