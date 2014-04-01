package rulengine

import (
	"rulengine/facts"
	"rulengine/logic"
	"testing"
)

func Test(t *testing.T) {
	engine := NewRuleEngine()
	engine.AddExpression("$user.age > 35", "age_larger_than_35")
	engine.AddExpression("$user.gender == \"male\"", "male")
	engine.AddExpression("$user.gender == \"female\"", "female")
	engine.AddRule(&logic.Rule{Expression: "age_larger_than_35 & male", Action: "Pass"})

	data := facts.NewFactCollection()
	user := `
		{
			"age" : 47,
			"gender" : "male"
		}
	`
	fact := facts.NewFact(user)
	data.Add("user", fact)

	exprNames := engine.GetFiredExpressions(data)
	if len(exprNames) == 0 {
		t.Error()
	}

	actions := engine.GetAction(data)
	if actions[0].Name != "Pass" {
		t.Error(actions)
	}
}

func TestCreditPay(t *testing.T) {
	engine := NewRuleEngine()
	engine.AddExpression("($behavior.money > $behavior.user.limit * 0.9) && $behavior.type == \"consume\"", "creditpay_consume_larger_than_90_percent_limit")
	engine.AddRule(&logic.Rule{Expression: "creditpay_consume_larger_than_90_percent_limit", Action: "alert"})

	data := facts.NewFactCollection()
	user := `
		{
			"type" : "consume",
			"money" : 91,
			"user" : {
				"limit" : 100
			}
		}
	`
	fact := facts.NewFact(user)
	data.Add("behavior", fact)

	exprNames := engine.GetFiredExpressions(data)
	if len(exprNames) == 0 {
		t.Error()
	}

	actions := engine.GetAction(data)
	if len(actions) == 0 {
		t.Error()
	}
	if actions[0].Name != "alert" {
		t.Error()
	}
}
