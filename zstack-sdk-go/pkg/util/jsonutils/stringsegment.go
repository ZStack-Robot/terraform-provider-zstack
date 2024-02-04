package jsonutils

import (
	"strconv"
	"strings"
)

type TextNumber struct {
	text     string
	number   int64
	isNumber bool
}

func (tn TextNumber) compare(tn2 TextNumber) int {
	if tn.isNumber && !tn2.isNumber {
		return -1
	} else if !tn.isNumber && tn2.isNumber {
		return 1
	} else if tn.isNumber && tn2.isNumber {
		if tn.number < tn2.number {
			return -1
		} else if tn.number > tn2.number {
			return 1
		} else {
			return 0
		}
	} else {
		if tn.text < tn2.text {
			return -1
		} else if tn.text > tn2.text {
			return 1
		} else {
			return 0
		}
	}
}

func (tn TextNumber) String() string {
	if tn.isNumber {
		return strconv.FormatInt(tn.number, 10)
	} else {
		return tn.text
	}
}

func string2TextNumber(str string) TextNumber {
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return TextNumber{text: str, isNumber: false}
	} else {
		return TextNumber{number: num, isNumber: true}
	}
}

func string2Segments(str string) []TextNumber {
	segs := strings.Split(str, ".")
	ret := make([]TextNumber, len(segs))
	for i := range segs {
		ret[i] = string2TextNumber(segs[i])
	}
	return ret
}

func segments2string(segs []TextNumber) string {
	segStrs := make([]string, len(segs))
	for i := range segs {
		segStrs[i] = segs[i].String()
	}
	return strings.Join(segStrs, ".")
}

type StringSegments [][]TextNumber

func (ss StringSegments) Len() int      { return len(ss) }
func (ss StringSegments) Swap(i, j int) { ss[i], ss[j] = ss[j], ss[i] }
func (ss StringSegments) Less(i, j int) bool {
	if len(ss[i]) < len(ss[j]) {
		return true
	} else if len(ss[i]) > len(ss[j]) {
		return false
	}
	for ii := range ss[i] {
		ret := ss[i][ii].compare(ss[j][ii])
		if ret < 0 {
			return true
		} else if ret > 0 {
			return false
		}
	}
	return false
}

func strings2stringSegments(strs []string) StringSegments {
	ret := make([][]TextNumber, len(strs))
	for i := range strs {
		ret[i] = string2Segments(strs[i])
	}
	return ret
}
