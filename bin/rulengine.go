package main

import (
	"flag"
	"fmt"
	"rulengine"
	"strconv"
	"strings"
)

func GuessType(str string) interface{} {
	if str == "true" {
		return true
	}
	if str == "false" {
		return false
	}
	{
		n, err := strconv.ParseInt(str, 10, 64)
		if err == nil {
			return int(n)
		}
	}
	{
		n, err := strconv.ParseFloat(str, 64)
		if err == nil {
			return n
		}
	}
	return str
}

func main() {
	data := flag.String("data", "", "data")
	rules := flag.String("rules", "", "rules file path")
	flag.Parse()

	tks := strings.Split(*data, "&")
	obj := make(map[string]interface{})
	for _, tk := range tks {
		kv := strings.Split(tk, "=")
		obj["$"+kv[0]] = GuessType(kv[1])
	}
	fmt.Println(obj)

	ruleEngine := rulengine.NewRuleEngine()
	ruleEngine.Load(*rules)
	actions := ruleEngine.GetAction(obj)
	fmt.Println("--------------")
	for _, action := range actions {
		fmt.Println(action)
	}
}
