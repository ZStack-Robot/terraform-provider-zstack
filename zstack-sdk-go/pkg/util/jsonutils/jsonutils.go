// Copyright (c) ZStack.io, Inc.

package jsonutils

import (
	"bytes"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"zstack.io/zstack-sdk-go/pkg/errors"
	"zstack.io/zstack-sdk-go/pkg/util/gotypes"
	"zstack.io/zstack-sdk-go/pkg/util/sortedmap"
)

type JSONObject interface {
	gotypes.ISerializable

	parse(str []byte, offset int) (int, error)
	writeSource

	// String() string
	PrettyString() string
	prettyString(level int) string
	YAMLString() string
	QueryString() string
	_queryString(key string) string
	Contains(keys ...string) bool
	ContainsIgnoreCases(keys ...string) bool
	Get(keys ...string) (JSONObject, error)
	GetIgnoreCases(keys ...string) (JSONObject, error)
	GetAt(i int, keys ...string) (JSONObject, error)
	Int(keys ...string) (int64, error)
	Float(keys ...string) (float64, error)
	Bool(keys ...string) (bool, error)
	GetMap(keys ...string) (map[string]JSONObject, error)
	GetArray(keys ...string) ([]JSONObject, error)
	GetTime(keys ...string) (time.Time, error)
	GetString(keys ...string) (string, error)
	Unmarshal(obj interface{}, keys ...string) error
	Equals(obj JSONObject) bool
	unmarshalValue(val reflect.Value) error
	// IsZero() bool
	Interface() interface{}
	isCompond() bool
}

type JSONValue struct {
}

var (
	JSONNull  = &JSONValue{}
	JSONTrue  = &JSONBool{data: true}
	JSONFalse = &JSONBool{data: false}
)

type JSONDict struct {
	JSONValue
	data sortedmap.SortedMap
}

type JSONArray struct {
	JSONValue
	data []JSONObject
}

type JSONString struct {
	JSONValue
	data string
}

type JSONInt struct {
	JSONValue
	data int64
}

type JSONFloat struct {
	JSONValue
	data float64
	bit  int
}

type JSONBool struct {
	JSONValue
	data bool
}

func skipEmpty(str []byte, offset int) int {
	i := offset
	for i < len(str) {
		switch str[i] {
		case ' ', '\t', '\n', '\r':
			i++
		default:
			return i
		}
	}
	return i
}

func hexchar2num(v byte) (byte, error) {
	switch {
	case v >= '0' && v <= '9':
		return v - '0', nil
	case v >= 'a' && v <= 'f':
		return v - 'a' + 10, nil
	case v >= 'A' && v <= 'F':
		return v - 'A' + 10, nil
	default:
		return 0, ErrInvalidChar // fmt.Errorf("Illegal char %c", v)
	}
}

func hexstr2byte(str []byte) (byte, error) {
	if len(str) < 2 {
		return 0, ErrInvalidHex // fmt.Errorf("Input must be 2 hex chars")
	}
	v1, e := hexchar2num(str[0])
	if e != nil {
		return 0, e
	}
	v2, e := hexchar2num(str[1])
	if e != nil {
		return 0, e
	}
	return v1*16 + v2, nil
}

func hexstr2rune(str []byte) (rune, error) {
	if len(str) < 4 {
		return 0, ErrInvalidRune // fmt.Errorf("Input must be 4 hex chars")
	}
	v1, e := hexstr2byte(str[0:2])
	if e != nil {
		return 0, e
	}
	v2, e := hexstr2byte(str[2:4])
	if e != nil {
		return 0, e
	}
	return rune(v1)*256 + rune(v2), nil
}

func parseQuoteString(str []byte, offset int, quotec byte) (string, int, error) {
	var (
		buffer    []byte
		runebytes = make([]byte, 4)
		runen     int
		i         = offset
	)
ret:
	for i < len(str) {
		switch str[i] {
		case '\\':
			if i+1 < len(str) {
				i++
				switch str[i] {
				case 'u':
					i++
					if i+4 >= len(str) {
						return "", i, NewJSONError(str, i, "Incomplete unicode")
					}
					r, e := hexstr2rune(str[i : i+4])
					if e != nil {
						return "", i, NewJSONError(str, i, e.Error())
					}
					runen = utf8.EncodeRune(runebytes, r)
					buffer = append(buffer, runebytes[0:runen]...)
					i += 4
				case 'x':
					i++
					if i+2 >= len(str) {
						return "", i, NewJSONError(str, i, "Incomplete hex")
					}
					b, e := hexstr2byte(str[i : i+2])
					if e != nil {
						return "", i, NewJSONError(str, i, e.Error())
					}
					buffer = append(buffer, b)
					i += 2
				case 'n':
					buffer = append(buffer, '\n')
					i++
				case 'r':
					buffer = append(buffer, '\r')
					i++
				case 't':
					buffer = append(buffer, '\t')
					i++
				case 'b':
					buffer = append(buffer, '\b')
					i++
				case 'f':
					buffer = append(buffer, '\f')
					i++
				case '\\':
					buffer = append(buffer, '\\')
					i++
				default:
					buffer = append(buffer, str[i])
					i++
				}
			} else {
				return "", i, NewJSONError(str, i, "Incomplete escape")
			}
		case quotec:
			i++
			break ret
		default:
			buffer = append(buffer, str[i])
			i++
		}
	}
	return string(buffer), i, nil
}

func parseString(str []byte, offset int) (string, bool, int, error) {
	var (
		i = offset
	)
	if c := str[i]; c == '"' || c == '\'' {
		r, newOfs, err := parseQuoteString(str, i+1, c)
		return r, true, newOfs, err
	}
ret2:
	for i < len(str) {
		switch str[i] {
		case ' ', ':', ',', '\t', '\r', '\n', '}', ']':
			break ret2
		default:
			i++
		}
	}
	return string(str[offset:i]), false, i, nil
}

func parseJSONValue(str []byte, offset int) (JSONObject, int, error) {
	val, quote, i, e := parseString(str, offset)
	if e != nil {
		return nil, i, errors.Wrap(e, "parseString")
	} else if quote {
		return &JSONString{data: val}, i, nil
	} else {
		lval := strings.ToLower(val)
		if len(lval) == 0 || lval == "null" || lval == "none" {
			return JSONNull, i, nil
		}
		if lval == "true" || lval == "yes" {
			return JSONTrue, i, nil
		}
		if lval == "false" || lval == "no" {
			return JSONFalse, i, nil
		}
		ival, err := strconv.ParseInt(val, 10, 64)
		if err == nil {
			return &JSONInt{data: ival}, i, nil
		}
		fval, err := strconv.ParseFloat(val, 64)
		if err == nil {
			return &JSONFloat{data: fval}, i, nil
		}
		return &JSONString{data: val}, i, nil
	}
}

// https://www.ietf.org/rfc/rfc4627.txt
//
//	string = quotation-mark *char quotation-mark
//
//	char = unescaped /
//	       escape (
//	           %x22 /          ; "    quotation mark  U+0022
//	           %x5C /          ; \    reverse solidus U+005C
//	           %x2F /          ; /    solidus         U+002F
//	           %x62 /          ; b    backspace       U+0008
//	           %x66 /          ; f    form feed       U+000C
//	           %x6E /          ; n    line feed       U+000A
//	           %x72 /          ; r    carriage return U+000D
//	           %x74 /          ; t    tab             U+0009
//	           %x75 4HEXDIG )  ; uXXXX                U+XXXX
//
//	escape = %x5C              ; \
//
//	quotation-mark = %x22      ; "
//
//	unescaped = %x20-21 / %x23-5B / %x5D-10FFFF
func escapeJsonChar(sb *strings.Builder, ch byte) {
	switch ch {
	case '"':
		sb.Write([]byte{'\\', '"'})
	case '\\':
		sb.Write([]byte{'\\', '\\'})
	case '\b':
		sb.Write([]byte{'\\', 'b'})
	case '\f':
		sb.Write([]byte{'\\', 'f'})
	case '\n':
		sb.Write([]byte{'\\', 'n'})
	case '\r':
		sb.Write([]byte{'\\', 'r'})
	case '\t':
		sb.Write([]byte{'\\', 't'})
	default:
		sb.WriteByte(ch)
		/*if ((ch >= 0x20 && ch <= 0x21) || (ch >= 0x23 || ch <= 0x5B) || (ch >= 0x5D && ch <= 0x10FFFF)) && ch != 0x81 && ch != 0x8d && ch != 0x8f && ch != 0x90 && ch != 0x9d {
			sb.WriteRune(ch)
		} else if ch <= 0xff {
			sb.Write([]byte{'\\', 'x'})
			sb.WriteString(fmt.Sprintf("%02x", ch))
		} else if ch <= 0xffff {
			sb.Write([]byte{'\\', 'u'})
			sb.WriteString(fmt.Sprintf("%04x", ch))
		} else {
			sb.Write([]byte{'\\', 'u'})
			sb.WriteString(fmt.Sprintf("%04x", ch>>16))
			sb.Write([]byte{'\\', 'u'})
			sb.WriteString(fmt.Sprintf("%04x", (ch & 0xffff)))
		}*/
	}
}

func quoteString(str string) string {
	sb := &strings.Builder{}
	sb.Grow(len(str) + 2)
	sb.WriteByte('"')
	for i := 0; i < len(str); i += 1 {
		escapeJsonChar(sb, str[i])
	}
	sb.WriteByte('"')
	return sb.String()
}

func jsonPrettyString(o JSONObject, level int) string {
	var buffer bytes.Buffer
	for i := 0; i < level; i++ {
		buffer.WriteString("  ")
	}
	buffer.WriteString(o.String())
	return buffer.String()
}

func (th *JSONString) PrettyString() string {
	return th.String()
}

func (th *JSONString) prettyString(level int) string {
	return jsonPrettyString(th, level)
}

func (th *JSONValue) parse(str []byte, offset int) (int, error) {
	return 0, nil
}

func (th *JSONValue) PrettyString() string {
	return th.String()
}

func (th *JSONValue) prettyString(level int) string {
	return jsonPrettyString(th, level)
}

func (th *JSONInt) PrettyString() string {
	return th.String()
}

func (th *JSONInt) prettyString(level int) string {
	return jsonPrettyString(th, level)
}

func (th *JSONFloat) PrettyString() string {
	return th.String()
}

func (th *JSONFloat) prettyString(level int) string {
	return jsonPrettyString(th, level)
}

func (th *JSONBool) PrettyString() string {
	return th.String()
}

func (th *JSONBool) prettyString(level int) string {
	return jsonPrettyString(th, level)
}

func parseDict(str []byte, offset int) (sortedmap.SortedMap, int, error) {
	smap := sortedmap.NewSortedMap()
	if str[offset] != '{' {
		return smap, offset, NewJSONError(str, offset, "{ not found")
	}
	var i = offset + 1
	var e error = nil
	var key string
	var stop = false
	for !stop && i < len(str) {
		i = skipEmpty(str, i)
		if i >= len(str) {
			return smap, i, NewJSONError(str, i, "Truncated")
		}
		if str[i] == '}' {
			stop = true
			i++
			continue
		}
		key, _, i, e = parseString(str, i)
		if e != nil {
			return smap, i, errors.Wrap(e, "parseString")
		}
		if i >= len(str) {
			return smap, i, NewJSONError(str, i, "Truncated")
		}
		i = skipEmpty(str, i)
		if i >= len(str) {
			return smap, i, NewJSONError(str, i, "Truncated")
		}
		if str[i] != ':' {
			return smap, i, NewJSONError(str, i, ": not found")
		}
		i++
		i = skipEmpty(str, i)
		if i >= len(str) {
			return smap, i, NewJSONError(str, i, "Truncated")
		}
		var val JSONObject = nil
		switch str[i] {
		case '[':
			val = &JSONArray{}
			i, e = val.parse(str, i)
		case '{':
			val = &JSONDict{}
			i, e = val.parse(str, i)
		default:
			val, i, e = parseJSONValue(str, i)
		}
		if e != nil {
			return smap, i, errors.Wrap(e, "parse misc")
		}
		smap = sortedmap.Add(smap, key, val)
		i = skipEmpty(str, i)
		if i >= len(str) {
			return smap, i, NewJSONError(str, i, "Truncated")
		}
		switch str[i] {
		case ',':
			i++
		case '}':
			i++
			stop = true
		default:
			return smap, i, NewJSONError(str, i, "Unexpected char")
		}
	}
	return smap, i, nil
}

func parseArray(str []byte, offset int) ([]JSONObject, int, error) {
	if str[offset] != '[' {
		return nil, offset, NewJSONError(str, offset, "[ not found")
	}
	var (
		list []JSONObject
		i    = offset + 1
		val  JSONObject
		e    error
		stop bool
	)
	for !stop && i < len(str) {
		i = skipEmpty(str, i)
		if i >= len(str) {
			return list, i, NewJSONError(str, i, "Truncated")
		}
		switch str[i] {
		case ']':
			i++
			stop = true
			continue
		case '[':
			val = &JSONArray{}
			i, e = val.parse(str, i)
		case '{':
			val = &JSONDict{}
			i, e = val.parse(str, i)
		default:
			val, i, e = parseJSONValue(str, i)
		}
		if e != nil {
			return list, i, errors.Wrap(e, "parse misc")
		}
		if i >= len(str) {
			return list, i, NewJSONError(str, i, "Truncated")
		}
		list = append(list, val)
		i = skipEmpty(str, i)
		if i >= len(str) {
			return list, i, NewJSONError(str, i, "Truncated")
		}
		switch str[i] {
		case ',':
			i++
		case ']':
			i++
			stop = true
		default:
			return list, i, NewJSONError(str, i, "Unexpected char")
		}
	}
	return list, i, nil
}

func (th *JSONDict) parse(str []byte, offset int) (int, error) {
	smap, i, e := parseDict(str, offset)
	if e == nil {
		th.data = smap
		return i, nil
	}
	return i, errors.Wrap(e, "parseDict")
}

func (th *JSONDict) SortedKeys() []string {
	return th.data.Keys()
}

func (th *JSONDict) PrettyString() string {
	return th.prettyString(0)
}

func (th *JSONDict) prettyString(level int) string {
	var buffer bytes.Buffer
	var linebuf bytes.Buffer
	for i := 0; i < level; i++ {
		linebuf.WriteString("  ")
	}
	var tab = linebuf.String()
	buffer.WriteString(tab)
	buffer.WriteByte('{')
	var idx = 0
	for iter := sortedmap.NewIterator(th.data); iter.HasMore(); iter.Next() {
		k, vInf := iter.Get()
		v := vInf.(JSONObject)
		if idx > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteByte('\n')
		buffer.WriteString(tab)
		buffer.WriteString("  ")
		buffer.WriteByte('"')
		buffer.WriteString(k)
		buffer.WriteString("\":")
		_, okdict := v.(*JSONDict)
		_, okarray := v.(*JSONArray)
		if okdict || okarray {
			buffer.WriteByte('\n')
			buffer.WriteString(v.prettyString(level + 2))
		} else {
			buffer.WriteByte(' ')
			buffer.WriteString(v.String())
		}
		idx++
	}
	if len(th.data) > 0 {
		buffer.WriteByte('\n')
		buffer.WriteString(tab)
	}
	buffer.WriteByte('}')
	return buffer.String()
}

func (th *JSONArray) parse(str []byte, offset int) (int, error) {
	val, i, e := parseArray(str, offset)
	if e == nil {
		th.data = val
	}
	return i, errors.Wrap(e, "parseArray")
}

func (th *JSONArray) PrettyString() string {
	return th.prettyString(0)
}

func (th *JSONArray) prettyString(level int) string {
	var buffer bytes.Buffer
	var linebuf bytes.Buffer
	for i := 0; i < level; i++ {
		linebuf.WriteString("  ")
	}
	var tab = linebuf.String()
	buffer.WriteString(tab)
	buffer.WriteByte('[')
	for idx, v := range th.data {
		if idx > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteByte('\n')
		buffer.WriteString(v.prettyString(level + 1))
	}
	if len(th.data) > 0 {
		buffer.WriteByte('\n')
		buffer.WriteString(tab)
	}
	buffer.WriteByte(']')
	return buffer.String()
}

func ParseString(str string) (JSONObject, error) {
	return Parse([]byte(str))
}

func Parse(str []byte) (JSONObject, error) {
	var i = 0
	i = skipEmpty(str, i)
	var val JSONObject = nil
	var e error = nil
	if i < len(str) {
		switch str[i] {
		case '{':
			val = &JSONDict{}
			_, e = val.parse(str, i)
		case '[':
			val = &JSONArray{}
			_, e = val.parse(str, i)
		default:
			val, _, e = parseJSONValue(str, i)
		}
		if e != nil {
			return nil, errors.Wrap(e, "parse misc")
		} else {
			return val, nil
		}
	} else {
		return nil, NewJSONError(str, i, "Empty string")
	}
}
