package main

import (
	"log"
	"rulengine"
	"rulengine/facts"
	"rulengine/logic"
)

func main() {
	engine := rulengine.NewRuleEngine()
	engine.AddExpression("$user.age > 35", "age_larger_than_35")
	engine.AddExpression("$user.gender == male", "male")
	engine.AddExpression("$user.gender == female", "female")
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
	log.Println(exprNames)

	actions := engine.GetAction(data)
	log.Println(actions)
}
