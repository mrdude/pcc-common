package ss

import (
	"sort"
)

// StringSet is a sorted set of strings
type StringSet struct {
	slice []string
}

func NewStringSet() *StringSet {
	return &StringSet{}
}

func (ss *StringSet) Add(value string) {
	if ss.slice == nil {
		ss.slice = make([]string, 0, 0)
		ss.slice = append(ss.slice, value)
		return
	}

	idx := sort.SearchStrings(ss.slice, value)
	if idx < len(ss.slice) && ss.slice[idx] == value {
		//value already exists
		return
	}

	// insert value
	if idx == len(ss.slice) {
		ss.slice = append(ss.slice, value)
	} else {
		var newSlice []string = make([]string, 0, len(ss.slice)+1)

		if idx > 0 {
			newSlice = append(newSlice, ss.slice[0:idx]...)
		}

		newSlice = append(newSlice, value)
		newSlice = append(newSlice, ss.slice[idx:]...)
		ss.slice = newSlice
	}
}

func (ss *StringSet) Remove(value string) {
	if ss.slice == nil || len(ss.slice) == 0 {
		return
	}

	if len(ss.slice) == 1 && ss.slice[0] == value {
		ss.slice = nil
		return
	}

	idx := sort.SearchStrings(ss.slice, value)
	if idx < len(ss.slice) && ss.slice[idx] == value {
		// ?value already exists

		if idx == 0 {
			ss.slice = ss.slice[1:]
		} else if idx == len(ss.slice) {
			ss.slice = ss.slice[:len(ss.slice)-1]
		} else {
			var newSlice []string = make([]string, 0, len(ss.slice)-1)
			newSlice = append(newSlice, ss.slice[0:idx]...)
			newSlice = append(newSlice, ss.slice[idx+1:]...)
			ss.slice = newSlice
		}
	}
}

func (ss *StringSet) Contains(value string) bool {
	if ss.slice == nil || len(ss.slice) == 0 {
		return false
	}

	if len(ss.slice) == 1 {
		return ss.slice[0] == value
	}

	idx := sort.SearchStrings(ss.slice, value)
	return idx < len(ss.slice) && ss.slice[idx] == value
}

// Clear removes all elements from the StringSet
func (ss *StringSet) Clear() {
	if len(ss.slice) != 0 {
		ss.slice = nil
	}
}

func (ss *StringSet) AddAll(other *StringSet) {
	// TODO implement this more efficiently?
	for _, s := range other.Slice() {
		ss.Add(s)
	}
}

func (ss *StringSet) RemoveAll(other *StringSet) {
	// TODO implement this more efficiently?
	for _, s := range other.Slice() {
		ss.Remove(s)
	}
}

// do not modify the returned slice
func (ss *StringSet) Slice() []string {
	return ss.slice
}

func (ss *StringSet) Len() int {
	return len(ss.slice)
}
