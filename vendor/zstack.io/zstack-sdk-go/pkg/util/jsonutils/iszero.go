package jsonutils

func (th *JSONValue) IsZero() bool {
	return true
}

func (th *JSONBool) IsZero() bool {
	return !th.data
}

func (th *JSONInt) IsZero() bool {
	return th.data == 0
}

func (th *JSONFloat) IsZero() bool {
	return th.data == 0.0
}

func (th *JSONString) IsZero() bool {
	return len(th.data) == 0
}

func (th *JSONDict) IsZero() bool {
	return th.data == nil || len(th.data) == 0
}

func (th *JSONArray) IsZero() bool {
	return th.data == nil || len(th.data) == 0
}
