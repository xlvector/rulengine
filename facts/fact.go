package facts

import (
	"encoding/json"
	"strings"
)

type Fact struct {
	data map[string]interface{}
	keys []string
}

func NewFact(buf string) *Fact {
	ret := Fact{}
	ret.data = make(map[string]interface{})
	json.Unmarshal([]byte(buf), &ret.data)
	ret.keys = ret.extractKeys(ret.data)
	return &ret
}

func (self *Fact) Get(key string) (interface{}, bool) {
	tks := strings.Split(key, ".")
	var data interface{}
	data = self.data
	for _, tk := range tks {
		dict, ok := data.(map[string]interface{})
		if ok {
			sub, exist := dict[tk]
			if exist {
				data = sub
			} else {
				return nil, false
			}
		}
	}
	return data, true
}

func (self *Fact) extractKeys(data map[string]interface{}) []string {
	ret := []string{}
	for key, val := range data {
		dict, ok := val.(map[string]interface{})
		if ok {
			subKeys := self.extractKeys(dict)
			for _, k := range subKeys {
				ret = append(ret, key+"."+k)
			}
			break
		}

		ret = append(ret, key)
	}
	return ret
}

func (self *Fact) Keys() []string {
	return self.keys
}

type FactCollection struct {
	facts map[string]*Fact
}

func NewFactCollection() *FactCollection {
	ret := FactCollection{}
	ret.facts = make(map[string]*Fact)
	return &ret
}

func (self *FactCollection) Get(key string) (interface{}, bool) {
	tks := strings.SplitN(key, ".", 2)
	if len(tks) != 2 {
		return nil, false
	}
	fact, exist := self.facts[tks[0]]
	if exist {
		return fact.Get(tks[1])
	}
	return nil, false
}

func (self *FactCollection) Add(key string, fact *Fact) {
	self.facts["$"+key] = fact
}

func (self *FactCollection) Keys() []string {
	ret := []string{}
	for key, ft := range self.facts {
		for _, sk := range ft.Keys() {
			ret = append(ret, key+"."+sk)
		}
	}
	return ret
}
