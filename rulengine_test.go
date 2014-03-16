package rulengine

import (
	"rulengine/logic"
	"testing"
)

func Test(t *testing.T) {
	engine := NewRuleEngine()
	engine.AddExpression("$age > 35", "age_larger_than_35")
	engine.AddExpression("$gender == male", "male")
	engine.AddExpression("$gender == female", "female")
	engine.AddRule(&logic.Rule{Expression: "age_larger_than_35 * male", Action: "Pass"})

	data := make(map[string]interface{})
	data["$age"] = 47
	data["$gender"] = "male"

	exprNames := engine.GetFiredExpressions(data)
	t.Error(exprNames)

	actions := engine.GetAction(data)
	t.Error(actions)
}
