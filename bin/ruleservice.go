package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rulengine"
	"rulengine/facts"
	"time"
)

type RuleService struct {
	ruleEngine *rulengine.RuleEngine
	keys       []string
}

func NewRuleService(keys []string) *RuleService {
	ret := RuleService{}
	ret.ruleEngine = rulengine.NewRuleEngine()
	ret.ruleEngine.Load("rules.txt")
	ret.keys = keys
	return &ret
}

func (self *RuleService) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	factCollection := facts.NewFactCollection()
	for _, key := range self.keys {
		fact := facts.NewFact(req.PostFormValue(key))
		factCollection.Add(key, fact)
	}
	actions := self.ruleEngine.GetAction(factCollection)
	output := rulengine.ConverActionListToActionRecords(actions)
	results, _ := json.Marshal(output)
	fmt.Fprint(w, string(results))
}

func main() {
	http.Handle("/risk/api/v1/creditpay/behavior/fraud", NewRuleService([]string{"behavior"}))
	http.Handle("/risk/api/v1/creditpay/application/fraud", NewRuleService([]string{"application"}))
	s := &http.Server{
		Addr:           ":8700",
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
