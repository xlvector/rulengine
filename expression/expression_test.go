package expression

import (
	"rulengine/facts"
	"testing"
)

func TestTokenize(t *testing.T) {
	if !IsOperatorCh('+') {
		t.Error()
	}
	if !IsVariableCh('5') {
		t.Error()
	}
	if !IsVariableCh('r') {
		t.Error()
	}
	if !IsVariableCh('E') {
		t.Error()
	}

	data := facts.NewFactCollection()
	fact := facts.NewFact(`
		{
			"A" : 1,
			"B" : 2,
			"C" : 3,
			"E" : 3,
			"F" : 2,
			"G" : 1,
			"H" : "hello world"
		}
	`)
	data.Add("u", fact)
	/*exprs := []string{"($A + $B) * $C > $E * $F + $G", "$A * 3 + 5 == 8",
	"$A > 0", "0 < $A", "0 > $A - $B", "0 > (($A - $B))", "0 != $C - 1",
	"(((($A)))) >= 1", "$A > 0 && $B == 2", "$H == \"hello world\""}*/
	exprs := []string{"$u.H == \"hello world\""}

	a, b := VariableValue("hello world", "$u.H", data)
	if a != b {
		t.Error(a, b)
	}

	for _, expr := range exprs {
		t.Log(expr)
		t.Log("tokenize")
		tks := Tokenize(expr)
		for _, tk := range tks {
			t.Log(tk)
		}

		pexpr := ToReversePolishNotation(tks)

		t.Log(pexpr)
		v := CalcReversePolishNotation(pexpr, data)
		if v != "true" {
			t.Error(expr, v)
		}
	}
}
