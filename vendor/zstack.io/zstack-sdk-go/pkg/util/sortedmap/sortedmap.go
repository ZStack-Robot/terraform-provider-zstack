package sortedmap

import (
	"strings"
)

const (
	InitCapa = 10
)

type Item struct {
	key   string
	value interface{}
}

type SortedMap []Item

func NewSortedMap() SortedMap {
	return NewSortedMapWithCapa(InitCapa)
}

func NewSortedMapWithCapa(capa int) SortedMap {
	return make([]Item, 0, capa)
}

func NewSortedMapFromMap(omap map[string]interface{}) SortedMap {
	return NewSortedMapFromMapWithCapa(omap, 0)
}

func NewSortedMapFromMapWithCapa(omap map[string]interface{}, capa int) SortedMap {
	if capa < len(omap) {
		capa = len(omap)
	}
	smap := NewSortedMapWithCapa(capa)
	for k := range omap {
		smap = Add(smap, k, omap[k])
	}
	return smap
}

func (ss SortedMap) Keys() []string {
	ret := make([]string, len(ss))
	for i := range ss {
		ret[i] = ss[i].key
	}
	return ret
}

func Add(ss SortedMap, key string, value interface{}) SortedMap {
	if ss == nil {
		ss = NewSortedMap()
	}
	pos, find := ss.find(key)
	if find {
		ss[pos].value = value
		return ss
	}
	item := Item{
		key:   key,
		value: value,
	}
	ss = append(ss, item)
	copy(ss[pos+1:], ss[pos:])
	ss[pos] = item
	return ss
}

func Delete(ss SortedMap, key string) (SortedMap, bool) {
	newsm, _, exist := deleteInternal(ss, key, false)
	return newsm, exist
}

func deleteInternal(ss SortedMap, key string, ignoreCase bool) (SortedMap, string, bool) {
	if ss == nil {
		return ss, "", false
	}
	var pos int
	var find bool
	if ignoreCase {
		pos, find = ss.findIgnoreCase(key)
	} else {
		pos, find = ss.find(key)
	}
	if !find {
		return ss, "", false
	}
	caseKey := ss[pos].key
	if pos < len(ss)-1 {
		copy(ss[pos:], ss[pos+1:])
	}
	ss = ss[:len(ss)-1]
	return ss, caseKey, true
}

func DeleteIgnoreCase(ss SortedMap, key string) (SortedMap, string, bool) {
	return deleteInternal(ss, key, true)
}

func (ss SortedMap) find(needle string) (int, bool) {
	i := 0
	j := len(ss) - 1
	for i <= j {
		m := (i + j) / 2
		if ss[m].key < needle {
			i = m + 1
		} else if ss[m].key > needle {
			j = m - 1
		} else {
			return m, true
		}
	}
	return j + 1, false
}

func (ss SortedMap) findIgnoreCase(needle string) (int, bool) {
	pos, exist := ss.find(needle)
	if exist {
		return pos, true
	}
	for i := range ss {
		if strings.EqualFold(ss[i].key, needle) {
			return i, true
		}
	}
	return -1, false
}

func (ss SortedMap) Get(key string) (interface{}, bool) {
	pos, find := ss.find(key)
	if find {
		return ss[pos].value, true
	} else {
		return nil, false
	}
}

func (ss SortedMap) GetIgnoreCase(key string) (interface{}, string, bool) {
	pos, find := ss.findIgnoreCase(key)
	if find {
		return ss[pos].value, ss[pos].key, true
	} else {
		return nil, "", false
	}
}

func (ss SortedMap) Contains(needle string) bool {
	_, find := ss.find(needle)
	return find
}

func (ss SortedMap) ContainsAny(needles ...string) bool {
	for i := range needles {
		_, find := ss.find(needles[i])
		if find {
			return true
		}
	}
	return false
}

func (ss SortedMap) ContainsAll(needles ...string) bool {
	for i := range needles {
		_, find := ss.find(needles[i])
		if !find {
			return false
		}
	}
	return true
}

func Split(a, b SortedMap) (a_b, anbA, anbB, b_a SortedMap) {
	minlen := len(a)
	if minlen > len(b) {
		minlen = len(b)
	}
	a_b = NewSortedMapWithCapa(len(a))
	b_a = NewSortedMapWithCapa(len(b))
	anbA = NewSortedMapWithCapa(minlen)
	anbB = NewSortedMapWithCapa(minlen)
	i := 0
	j := 0
	for i < len(a) && j < len(b) {
		if a[i].key == b[j].key {
			anbA = append(anbA, a[i])
			anbB = append(anbB, b[j])
			i += 1
			j += 1
		} else if a[i].key < b[j].key {
			a_b = append(a_b, a[i])
			i += 1
		} else if a[i].key > b[j].key {
			b_a = append(b_a, b[j])
			j += 1
		}
	}
	if i < len(a) {
		a_b = append(a_b, a[i:]...)
	}
	if j < len(b) {
		b_a = append(b_a, b[j:]...)
	}
	return
}

// order matters, values of the latter override the former
func Merge(a, b SortedMap) SortedMap {
	ret := NewSortedMapWithCapa(len(a) + len(b))
	i := 0
	j := 0
	for i < len(a) && j < len(b) {
		if a[i].key == b[j].key {
			ret = append(ret, b[j])
			i += 1
			j += 1
		} else if a[i].key < b[j].key {
			ret = append(ret, a[i])
			i += 1
		} else if a[i].key > b[j].key {
			ret = append(ret, b[j])
			j += 1
		}
	}
	if i < len(a) {
		ret = append(ret, a[i:]...)
	}
	if j < len(b) {
		ret = append(ret, b[j:]...)
	}
	return ret
}
