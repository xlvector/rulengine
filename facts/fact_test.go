package facts

import (
	"testing"
)

func Test(t *testing.T) {
	buf := `
		{
			"name" : "liangxiang",
			"age" : 30,
			"boss" : {
				"name" : "huazheng"
			}
		}
	`
	fact := NewFact(buf)
	name, _ := fact.Get("name")
	if name != "liangxiang" {
		t.Error(name)
	}

	bossName, _ := fact.Get("boss.name")
	if bossName != "huazheng" {
		t.Error(bossName)
	}

	age, _ := fact.Get("age")
	if age.(float64) != 30 {
		t.Error(age)
	}

	factCollection := NewFactCollection()
	factCollection.Add("behavior", fact)
	behaviorBossName, _ := factCollection.Get("behavior.boss.name")
	if behaviorBossName != "huazheng" {
		t.Error(behaviorBossName)
	}

	if len(factCollection.Keys()) == 0 {
		t.Error("fact collection keys empty")
	}
}
