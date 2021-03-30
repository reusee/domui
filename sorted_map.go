package domui

import "sort"

type SortedMapItem struct {
	Key   string
	Value any
}

type SortedMap []SortedMapItem

func (s *SortedMap) Set(k string, v any) {
	i := sort.Search(len(*s), func(i int) bool {
		return (*s)[i].Key >= k
	})
	if i < len(*s) {
		if (*s)[i].Key == k {
			// found
			(*s)[i].Value = v
		} else {
			// insert
			*s = append(
				(*s)[:i],
				append([]SortedMapItem{
					{
						Key:   k,
						Value: v,
					},
				}, (*s)[i:]...)...,
			)
		}
	} else {
		// append
		*s = append(*s, SortedMapItem{
			Key:   k,
			Value: v,
		})
	}
}

func (s SortedMap) Get(k string) (any, bool) {
	i := sort.Search(len(s), func(i int) bool {
		return s[i].Key >= k
	})
	if i < len(s) {
		if s[i].Key == k {
			return s[i].Value, true
		}
	}
	return nil, false
}
