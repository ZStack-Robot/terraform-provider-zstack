package sortedmap

type SortedMapIterator struct {
	smap  SortedMap
	index int
}

func (i *SortedMapIterator) Init(smap SortedMap) {
	i.smap = smap
	i.index = 0
}

func (i SortedMapIterator) HasMore() bool {
	return i.index < len(i.smap)
}

func (i *SortedMapIterator) Next() {
	i.index += 1
}

func (i SortedMapIterator) Get() (string, interface{}) {
	return i.smap[i.index].key, i.smap[i.index].value
}

func NewIterator(smap SortedMap) *SortedMapIterator {
	iter := &SortedMapIterator{}
	iter.Init(smap)
	return iter
}
