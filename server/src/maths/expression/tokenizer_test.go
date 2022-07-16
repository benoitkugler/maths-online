package expression

import (
	"reflect"
	"testing"
)

func allTokens(s string) (tokens, peekTokens []tokenData) {
	tk := newTokenizer([]byte(s))

	peekTokens = append(peekTokens, tk.Peek().data)

	for tok := tk.Next(); tok.data != nil; tok = tk.Next() {
		tokens = append(tokens, tok.data)
		if peeked := tk.Peek(); peeked.data != nil {
			peekTokens = append(peekTokens, peeked.data)
		}
	}
	return
}

func TestPeekToken(t *testing.T) {
	for _, test := range []string{
		"7x  + 8",
		"7*x",
		"(2x + 1)(4x)",
		"(x + 3)(x+4)",
		"(x + 3)*(x+4)",
		"7log(10)",
		"7randPrime(1;10)",
		" (1+ 2 ) (x + 3) ",
	} {
		l1, l2 := allTokens(test)
		if !reflect.DeepEqual(l1, l2) {
			t.Fatalf("invalid peek tokens for %s: %v != %v", test, l1, l2)
		}
	}
}

func TestTokens(t *testing.T) {
	for _, test := range []struct {
		expr   string
		tokens []tokenData
	}{
		{"7x  + 8", []tokenData{numberText("7"), mult, Variable{Name: 'x'}, plus, numberText("8")}},
		{"7*x", []tokenData{numberText("7"), mult, Variable{Name: 'x'}}},
		{"(2x + 1)(4x)", []tokenData{openPar, numberText("2"), mult, Variable{Name: 'x'}, plus, numberText("1"), closePar, mult, openPar, numberText("4"), mult, Variable{Name: 'x'}, closePar}},
		{"(x + 3)(x+4)", []tokenData{openPar, Variable{Name: 'x'}, plus, numberText("3"), closePar, mult, openPar, Variable{Name: 'x'}, plus, numberText("4"), closePar}},
		{"(x + 3)*(x+4)", []tokenData{openPar, Variable{Name: 'x'}, plus, numberText("3"), closePar, mult, openPar, Variable{Name: 'x'}, plus, numberText("4"), closePar}},
		{" (1+ 2 ) (x + 3) ", []tokenData{openPar, numberText("1"), plus, numberText("2"), closePar, mult, openPar, Variable{Name: 'x'}, plus, numberText("3"), closePar}},
		{"7log(10)", []tokenData{numberText("7"), mult, logFn, openPar, numberText("10"), closePar}},
		{"7randPrime(1;10)", []tokenData{numberText("7"), mult, randPrime, openPar, numberText("1"), semicolon, numberText("10"), closePar}},
		{"randChoice(1;2)", []tokenData{randChoice, openPar, numberText("1"), semicolon, numberText("2"), closePar}},
		{"randDecDen()", []tokenData{randDenominator, openPar, closePar}},
		{"randSymbol()", []tokenData{randVariable{}, openPar, closePar}},
		{"min(1)", []tokenData{minFn, openPar, numberText("1"), closePar}},
		{"max(1)", []tokenData{maxFn, openPar, numberText("1"), closePar}},
		{"floor(1)", []tokenData{floorFn, openPar, numberText("1"), closePar}},
		{"round(2.12; 5)", []tokenData{roundFn{}, openPar, numberText("2.12"), semicolon, numberText("5"), closePar}},
		{"inf - inf + inf", []tokenData{numberText("inf"), minus, numberText("inf"), plus, numberText("inf")}},
		{"@_myVar", []tokenData{Variable{Name: '@', Indice: "myVar"}}},
		{"2 > 4", []tokenData{numberText("2"), strictlyGreater, numberText("4")}},
		{"2 < 4", []tokenData{numberText("2"), strictlyLesser, numberText("4")}},
		{"2 <= 4", []tokenData{numberText("2"), lesser, numberText("4")}},
		{"2 >= 4", []tokenData{numberText("2"), greater, numberText("4")}},
		{"2 == 4", []tokenData{numberText("2"), equals, numberText("4")}},
	} {
		if got, _ := allTokens(test.expr); !reflect.DeepEqual(got, test.tokens) {
			t.Fatalf("for %s, expected %v, got %v", test.expr, test.tokens, got)
		}
	}
}

func TestVariableIndice(t *testing.T) {
	got, _ := allTokens("x_A + y_B + D_JIE")
	expected := []tokenData{
		Variable{Name: 'x', Indice: "A"},
		plus,
		Variable{Name: 'y', Indice: "B"},
		plus,
		Variable{Name: 'D', Indice: "JIE"},
	}
	if !reflect.DeepEqual(got, expected) {
		t.Fatal()
	}
}
