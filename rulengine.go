package rulengine

import (
	"bufio"
	"fmt"
	"os"
	"rulengine/expression"
	"rulengine/logic"
	"strings"
)

type RuleEngine struct {
	expressions           [][]string
	expressionNames       []string
	varCountInExpressions []int
	expressionIndex       map[string][]int
	rules                 [][]string
	actions               []string
	ruleIndex             map[string][]int
}

func NewRuleEngine() *RuleEngine {
	ret := RuleEngine{}
	ret.expressions = [][]string{}
	ret.expressionNames = []string{}
	ret.varCountInExpressions = []int{}
	ret.expressionIndex = make(map[string][]int)
	ret.ruleIndex = make(map[string][]int)
	ret.rules = [][]string{}
	ret.actions = []string{}
	return &ret
}

func (self *RuleEngine) AddExpression(expr string, name string) {
	tks := expression.Tokenize(expr)
	reversePolishExpr := expression.ToReversePolishNotation(tks)
	size := len(self.expressions)
	nv := 0
	for _, tk := range reversePolishExpr {
		if tk[0] == '$' {
			_, ok := self.expressionIndex[tk]
			if !ok {
				self.expressionIndex[tk] = []int{}
			}
			self.expressionIndex[tk] = append(self.expressionIndex[tk], size)
			nv += 1
		}
	}
	self.varCountInExpressions = append(self.varCountInExpressions, nv)
	self.expressionNames = append(self.expressionNames, name)
	self.expressions = append(self.expressions, reversePolishExpr)
}

func (self *RuleEngine) GetFiredExpressions(data map[string]interface{}) []string {
	counter := make(map[int]int)
	for k, _ := range data {
		ids, ok := self.expressionIndex[k]
		if ok {
			for _, id := range ids {
				v, exist := counter[id]
				if !exist {
					counter[id] = 1
				} else {
					counter[id] = v + 1
				}
			}
		}
	}

	firedExpressions := []string{}
	for k, v := range counter {
		if v == self.varCountInExpressions[k] {
			ret := expression.CalcReversePolishNotation(self.expressions[k], data)
			if ret == "true" {
				fmt.Println("fire expression", self.expressionNames[k], self.expressions[k])
				firedExpressions = append(firedExpressions, self.expressionNames[k])
			}
		}
	}
	return firedExpressions
}

func (self *RuleEngine) AddRule(rule *logic.Rule) {
	andOrRule := logic.AndOrFormat(rule.Expression)
	for _, e := range andOrRule.Sets {
		keys := strings.Split(e.ToString(), "\t")
		size := len(self.rules)
		self.rules = append(self.rules, keys)
		self.actions = append(self.actions, rule.Action)
		for _, key := range keys {
			_, ok := self.ruleIndex[key]
			if ok {
				self.ruleIndex[key] = append(self.ruleIndex[key], size)
			} else {
				self.ruleIndex[key] = []int{size}
			}
		}
	}
}

func (self *RuleEngine) GetAction(data map[string]interface{}) []string {
	exprs := self.GetFiredExpressions(data)
	ret := []string{}
	keys := make(map[string]bool)
	for _, key := range exprs {
		keys[key] = true
	}
	checked := make(map[string]bool)
	checkedRules := make(map[int]bool)
	counter := make(map[int]int)
	for {
		if len(keys) == 0 {
			break
		}
		for key, _ := range keys {
			checked[key] = true
			ids, ok := self.ruleIndex[key]
			if ok {
				for _, id := range ids {
					_, ok2 := counter[id]
					if ok2 {
						counter[id] += 1
					} else {
						counter[id] = 1
					}
				}
			}
		}

		for id, c := range counter {
			if len(self.rules[id]) == c {
				_, ok := checkedRules[id]
				if ok {
					continue
				}
				checkedRules[id] = true
				ret = append(ret, strings.Join(self.rules[id], "*")+" -> "+self.actions[id])
				_, ok = checked[self.actions[id]]
				if ok {
					continue
				}
				keys[self.actions[id]] = true
			}
		}

		for key, _ := range checked {
			delete(keys, key)
		}
	}

	return ret
}

func (self *RuleEngine) Load(fname string) {
	f, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.Trim(line, "\n")
		if strings.Contains(line, ":=") {
			tks := strings.Split(line, ":=")
			self.AddExpression(strings.Trim(tks[1], " "), strings.Trim(tks[0], " "))
			continue
		}

		if strings.Contains(line, "->") {
			tks := strings.Split(line, "->")

			rule := logic.Rule{Expression: tks[0], Action: strings.Trim(tks[1], " ")}
			self.AddRule(&rule)
		}
	}
}
