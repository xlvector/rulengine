package expression

import (
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

	data := make(map[string]interface{})
	data["$A"] = 1
	data["$B"] = 2
	data["$C"] = 3
	data["$E"] = 3
	data["$F"] = 2
	data["$G"] = 1
	data["$H"] = "hello"
	exprs := []string{"($A + $B) * $C > $E * $F + $G", "$A * 3 + 5 == 8",
		"$A > 0", "0 < $A", "0 > $A - $B", "0 > (($A - $B))", "0 != $C - 1",
		"(((($A)))) >= 1", "$A > 0 && $B == 2", "$H == hello"}
	for _, expr := range exprs {
		t.Log(expr)
		tks := Tokenize(expr)
		for _, tk := range tks {
			t.Log(tk)
		}

		pexpr := ToReversePolishNotation(tks)

		v := CalcReversePolishNotation(pexpr, data)
		if v != "true" {
			t.Error(expr, v)
		}
	}
}
