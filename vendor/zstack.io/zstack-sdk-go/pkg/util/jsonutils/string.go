package jsonutils

import (
	"fmt"
	"strconv"
	"strings"
)

func (th *JSONString) String() string {
	return quoteString(th.data)
}

func (th *JSONValue) String() string {
	return "null"
}

func (th *JSONInt) String() string {
	return fmt.Sprintf("%d", th.data)
}

func (th *JSONFloat) String() string {
	if th.bit != 32 && th.bit != 64 {
		th.bit = 64
	}
	return strconv.FormatFloat(th.data, 'f', -1, th.bit)
}

func (th *JSONBool) String() string {
	if th.data {
		return "true"
	} else {
		return "false"
	}
}

func (th *JSONDict) String() string {
	sb := &strings.Builder{}
	th.buildString(sb)
	return sb.String()
}

func (th *JSONArray) String() string {
	sb := &strings.Builder{}
	th.buildString(sb)
	return sb.String()
}
