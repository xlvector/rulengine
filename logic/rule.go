package logic

import (
	"sort"
	"strings"
)

type StringArray []string

func NewStringArray(m map[string]bool) StringArray {
	ms := make(StringArray, 0, len(m))

	for k, _ := range m {
		ms = append(ms, k)
	}

	return ms
}

func (self StringArray) Len() int {
	return len(self)
}

func (self StringArray) Less(i, j int) bool {
	return self[i] < self[j]
}

func (self StringArray) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

type AndSet struct {
	data map[string]bool
}

func NewAndSet() *AndSet {
	ret := AndSet{}
	ret.data = make(map[string]bool)
	return &ret
}

func (self *AndSet) ToString() string {
	keys := NewStringArray(self.data)
	sort.Sort(keys)
	ret := ""
	for _, key := range keys {
		ret = ret + key
		ret += "\t"
	}
	return strings.Trim(ret, "\t")
}

func (self *AndSet) Add(buf string) {
	self.data[buf] = true
}

func Union(a *AndSet, b *AndSet) *AndSet {
	ret := NewAndSet()
	for k, _ := range a.data {
		ret.data[k] = true
	}
	for k, _ := range b.data {
		ret.data[k] = true
	}
	return ret
}

type UnionedAndSet struct {
	Index map[string]bool
	Sets  []*AndSet
}

func NewUnionedAndSet() *UnionedAndSet {
	ret := UnionedAndSet{}
	ret.Index = make(map[string]bool)
	ret.Sets = []*AndSet{}
	return &ret
}

func (self *UnionedAndSet) Add(set *AndSet) {
	h := set.ToString()
	_, ok := self.Index[h]
	if !ok {
		self.Index[h] = true
		self.Sets = append(self.Sets, set)
	}
}

func TrimExp(exp string) string {
	exp = strings.Replace(exp, " ", "", -1)
	if exp[0] != '(' {
		return exp
	}
	num := 0
	for i, ch := range exp {
		if ch == '(' {
			num += 1
		} else if ch == ')' {
			num -= 1
		}
		if num == 0 {
			if i == len(exp)-1 {
				return exp[1 : len(exp)-1]
			} else {
				return exp
			}
		}
	}
	return exp
}

func AndOrFormat(exp string) *UnionedAndSet {
	exp = TrimExp(exp)
	num := 0
	for i, ch := range exp {
		if ch == '(' {
			num += 1
		} else if ch == ')' {
			num -= 1
		} else if ch == '&' {
			if num == 0 {
				left := AndOrFormat(exp[0:i])
				right := AndOrFormat(exp[i+1 : len(exp)])
				ret := NewUnionedAndSet()
				for _, le := range left.Sets {
					for _, re := range right.Sets {
						ret.Add(Union(le, re))
					}
				}
				return ret
			}
		} else if ch == '|' {
			if num == 0 {
				left := AndOrFormat(exp[0:i])
				right := AndOrFormat(exp[i+1 : len(exp)])
				for _, e := range right.Sets {
					left.Add(e)
				}
				return left
			}
		}
	}
	andSet := NewAndSet()
	andSet.Add(exp)
	ret := NewUnionedAndSet()
	ret.Add(andSet)
	return ret
}

type Rule struct {
	Expression string
	Action     string
}
