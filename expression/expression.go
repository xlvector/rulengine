package expression

import (
	"container/list"
	"rulengine/facts"
	"strconv"
	"strings"
)

func IsOperatorCh(ch rune) bool {
	if ch == '(' || ch == ')' || ch == '+' || ch == '-' || ch == '=' || ch == '!' || ch == '*' || ch == '/' || ch == '>' || ch == '<' || ch == '&' || ch == '|' {
		return true
	} else {
		return false
	}
}

func ShouldSplit(a rune, b rune) bool {
	if a == '!' && b == '=' {
		return false
	}
	if a == '=' && b == '=' {
		return false
	}
	if a == '>' && b == '=' {
		return false
	}
	if a == '<' && b == '=' {
		return false
	}
	if a == '&' && b == '&' {
		return false
	}
	if a == '|' && b == '|' {
		return false
	}
	if a == ')' || a == '(' {
		return true
	}

	if b == '(' || b == ')' {
		return true
	}
	if IsOperatorCh(a) && !IsOperatorCh(b) {
		return true
	}
	if !IsOperatorCh(a) && IsOperatorCh(b) {
		return true
	}
	return false
}

func IsVariableCh(ch rune) bool {
	if ch == '_' || ch == '$' || ch == '"' || ch == ' ' || ch == '.' {
		return true
	}
	if ch >= 'a' && ch <= 'z' {
		return true
	}
	if ch >= 'A' && ch <= 'Z' {
		return true
	}
	if ch >= '0' && ch <= '9' {
		return true
	}
	return false
}

func Tokenize(expr string) []string {
	ret := []string{}
	tmp := []byte{}
	prevCh := ' '
	nQuot := 0
	for i, ch := range expr {
		if ch == '"' {
			nQuot += 1
		}
		if i == 0 {
			tmp = append(tmp, byte(ch))
			prevCh = ch
		} else {
			if ch == '\t' || ch == '\n' {
				continue
			}
			if ch == ' ' && nQuot%2 == 0 {
				continue
			}
			if ShouldSplit(prevCh, ch) {
				str := strings.Trim(string(tmp), "\"")
				if len(str) > 0 {
					ret = append(ret, str)
				}
				tmp = []byte{}
			}

			if IsVariableCh(ch) || IsOperatorCh(ch) {
				tmp = append(tmp, byte(ch))
			}
			prevCh = ch
		}
	}
	str := strings.Trim(string(tmp), "\"")
	if len(str) > 0 {
		ret = append(ret, str)
	}
	if nQuot%2 != 0 {
		panic("synax error : quot number error")
	}
	return ret
}

func IsOperator(str string) bool {
	if str == "+" || str == "-" || str == "==" || str == "!=" || str == "*" || str == "/" || str == ">" || str == "<" || str == ">=" || str == "<=" || str == "&&" || str == "||" {
		return true
	} else {
		return false
	}
}

func PriorityHigherThan(oa string, ob string) bool {
	ops := []string{oa, ob}
	ps := []int{0, 0}

	for i := 0; i < len(ops); i++ {
		if ops[i] == "+" || ops[i] == "-" {
			ps[i] = 1
		} else if ops[i] == "*" || ops[i] == "/" {
			ps[i] = 2
		} else if ops[i] == "||" {
			ps[i] = -2
		} else if ops[i] == "&&" {
			ps[i] = -1
		}
	}
	return ps[0] > ps[1]
}

func ToReversePolishNotation(expr []string) []string {
	s1 := list.New()
	s2 := list.New()

	for _, tk := range expr {
		if IsOperator(tk) {
			for {
				if s1.Len() == 0 {
					s1.PushBack(tk)
					break
				} else {
					top := s1.Back().Value.(string)
					if top == "(" {
						s1.PushBack(tk)
						break
					} else {
						if PriorityHigherThan(tk, top) {
							s1.PushBack(tk)
							break
						} else {
							s2.PushBack(top)
							s1.Remove(s1.Back())
						}
					}
				}
			}
		} else if tk == "(" {
			s1.PushBack(tk)
		} else if tk == ")" {
			for {
				if s1.Len() == 0 {
					panic("Error expression")
				}
				top := s1.Back().Value.(string)
				s1.Remove(s1.Back())
				if top == "(" {
					break
				} else {
					s2.PushBack(top)
				}
			}
		} else {
			s2.PushBack(tk)
		}
	}

	for {
		if s1.Len() == 0 {
			break
		}
		top := s1.Back().Value.(string)
		s2.PushBack(top)
		s1.Remove(s1.Back())
	}

	ret := []string{}
	head := s2.Front()
	for {
		if head == nil {
			break
		}
		ret = append(ret, head.Value.(string))
		head = head.Next()
	}
	return ret
}

func VariableValue(a, b string, data *facts.FactCollection) (interface{}, interface{}) {
	va, oka := data.Get(a)
	vb, okb := data.Get(b)
	if oka && okb {
		return va, vb
	} else if oka && !okb {
		{
			na, ok := va.(int)
			if ok {
				nb, err := strconv.Atoi(b)
				if err != nil {
					panic(err)
				}
				return na, nb
			}
		}
		{
			na, ok := va.(float64)
			if ok {
				nb, err := strconv.ParseFloat(b, 64)
				if err != nil {
					panic(err)
				}
				return na, nb
			}
		}
		{
			na, ok := va.(string)
			if ok {
				return na, b
			}
		}
	} else if !oka && okb {
		{
			nb, ok := vb.(int)
			if ok {
				na, err := strconv.Atoi(a)
				if err != nil {
					panic(err)
				}
				return na, nb
			}
		}
		{
			nb, ok := vb.(float64)
			if ok {
				na, err := strconv.ParseFloat(a, 64)
				if err != nil {
					panic(err)
				}
				return na, nb
			}
		}
		{
			nb, ok := vb.(string)
			if ok {
				return a, nb
			}
		}
	} else {
		{
			na, err_na := strconv.Atoi(a)
			nb, err_nb := strconv.Atoi(b)

			if err_na == nil && err_nb == nil {
				return na, nb
			}
		}
		{
			na, err_na := strconv.ParseFloat(a, 64)
			nb, err_nb := strconv.ParseFloat(b, 64)

			if err_na == nil && err_nb == nil {
				return na, nb
			}
		}
	}
	return a, b
}

func IntNumberOp(na, nb int, op string) int {
	if op == "+" {
		return na + nb
	} else if op == "-" {
		return na - nb
	} else if op == "*" {
		return na * nb
	} else if op == "/" {
		return na / nb
	} else {
		panic("Does not support op " + op)
		return 0.0
	}
}

func FloatNumberOp(na, nb float64, op string) float64 {
	if op == "+" {
		return na + nb
	} else if op == "-" {
		return na - nb
	} else if op == "*" {
		return na * nb
	} else if op == "/" {
		return na / nb
	} else {
		panic("Does not support op " + op)
		return 0.0
	}
}

func NumberOp(a, b interface{}, op string) string {
	va, ok := a.(int)
	if ok {
		vb, okb := b.(int)
		if !okb {
			panic("b is not an integer")
		}
		return strconv.Itoa(IntNumberOp(va, vb, op))
	}

	fa, ok := a.(float64)
	if ok {
		vb, okb := b.(float64)
		if !okb {
			panic("b is not an float")
		}
		return strconv.FormatFloat(FloatNumberOp(fa, vb, op), 'g', 10, 64)
	}

	panic("val is not int or float64")
	return ""
}

func IntBoolOp(na, nb int, op string) bool {
	if op == "==" {
		return na == nb
	} else if op == "!=" {
		return na != nb
	} else if op == ">" {
		return na > nb
	} else if op == "<" {
		return na < nb
	} else if op == ">=" {
		return na >= nb
	} else if op == "<=" {
		return na <= nb
	} else {
		panic("Does not support operator " + op)
		return false
	}
}

func FloatBoolOp(na, nb float64, op string) bool {
	if op == "==" {
		return na == nb
	} else if op == "!=" {
		return na != nb
	} else if op == ">" {
		return na > nb
	} else if op == "<" {
		return na < nb
	} else if op == ">=" {
		return na >= nb
	} else if op == "<=" {
		return na <= nb
	} else {
		panic("Does not support operator " + op)
		return false
	}
}

func StringBoolOp(na, nb string, op string) bool {
	if op == "==" {
		return na == nb
	} else if op == "!=" {
		return na != nb
	} else if op == ">" {
		return na > nb
	} else if op == "<" {
		return na < nb
	} else if op == ">=" {
		return na >= nb
	} else if op == "<=" {
		return na <= nb
	} else {
		panic("Does not support operator " + op)
		return false
	}
}

func BoolOp(a, b interface{}, op string) string {
	va, ok := a.(int)
	if ok {
		vb, okb := b.(int)
		if !okb {
			panic("b is not an integer")
		}
		return strconv.FormatBool(IntBoolOp(va, vb, op))
	}

	fa, ok := a.(float64)
	if ok {
		vb, okb := b.(float64)
		if !okb {
			panic("b is not an float")
		}
		return strconv.FormatBool(FloatBoolOp(fa, vb, op))
	}

	sa, ok := a.(string)
	if ok {
		vb, okb := b.(string)
		if !okb {
			panic("b is not an float")
		}
		return strconv.FormatBool(StringBoolOp(sa, vb, op))
	}

	panic("val is not int or float64 or string")
	return ""
}

func LogicOp(a, b interface{}, op string) string {
	va := true
	if a.(string) == "false" {
		va = false
	}
	vb := true
	if b.(string) == "false" {
		vb = false
	}
	if op == "&&" {
		return strconv.FormatBool(va && vb)
	} else if op == "||" {
		return strconv.FormatBool(va || vb)
	}
	panic("operator " + op + " is not support")
	return ""
}

func Calc(a, b string, op string, data *facts.FactCollection) string {
	va, vb := VariableValue(a, b, data)
	if op == "+" || op == "-" || op == "*" || op == "/" {
		return NumberOp(va, vb, op)
	} else if op == "==" || op == "!=" || op == ">" || op == "<" || op == ">=" || op == "<=" {
		return BoolOp(va, vb, op)
	} else if op == "&&" || op == "||" {
		return LogicOp(va, vb, op)
	} else {
		panic("Does not support operator " + op)
		return ""
	}
}

func CalcReversePolishNotation(expr []string, data *facts.FactCollection) interface{} {
	s := list.New()
	for _, tk := range expr {
		if IsOperator(tk) {
			a := s.Back().Value.(string)
			s.Remove(s.Back())
			b := s.Back().Value.(string)
			s.Remove(s.Back())
			val := Calc(b, a, tk, data)
			s.PushBack(val)
		} else {
			s.PushBack(tk)
		}
	}
	return s.Front().Value
}
