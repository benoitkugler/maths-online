package expression

import (
	"reflect"
	"testing"
)

func TestTokens(t *testing.T) {
	x := Variable{Name: 'x'}
	for _, test := range []struct {
		expr   string
		tokens []tokenData
	}{
		{"-2", []tokenData{minus, nT("2")}},
		{"7x  + 8", []tokenData{nT("7"), mult, x, plus, nT("8")}},
		{"7*x", []tokenData{nT("7"), mult, x}},
		{"(2x + 1)(4x)", []tokenData{openPar, nT("2"), mult, x, plus, nT("1"), closePar, mult, openPar, nT("4"), mult, x, closePar}},
		{"(x + 3)(x+4)", []tokenData{openPar, x, plus, nT("3"), closePar, mult, openPar, x, plus, nT("4"), closePar}},
		{"(x + 3)*(x+4)", []tokenData{openPar, x, plus, nT("3"), closePar, mult, openPar, x, plus, nT("4"), closePar}},
		{" (1+ 2 ) (x + 3) ", []tokenData{openPar, nT("1"), plus, nT("2"), closePar, mult, openPar, x, plus, nT("3"), closePar}},
		{"7log(10)", []tokenData{nT("7"), mult, logFn, openPar, nT("10"), closePar}},
		{"7randPrime(1;10)", []tokenData{nT("7"), mult, randPrime, openPar, nT("1"), semicolon, nT("10"), closePar}},
		{"randChoice(1;2)", []tokenData{randChoice, openPar, nT("1"), semicolon, nT("2"), closePar}},
		{"randMatrix(1;2)", []tokenData{randMatrixInt, openPar, nT("1"), semicolon, nT("2"), closePar}},
		{"randDecDen()", []tokenData{randDenominator, openPar, closePar}},
		{"choiceFrom()", []tokenData{choiceFrom, openPar, closePar}},
		{"choiceFrom(x ; 2 )", []tokenData{choiceFrom, openPar, x, semicolon, nT("2"), closePar}},
		{"min(1)", []tokenData{minFn, openPar, nT("1"), closePar}},
		{"max(1)", []tokenData{maxFn, openPar, nT("1"), closePar}},
		{"sum(1)", []tokenData{sumFn, openPar, nT("1"), closePar}},
		{"forceDecimal(1/2)", []tokenData{forceDecimalFn, openPar, nT("1"), div, nT("2"), closePar}},
		{"floor(1)", []tokenData{floorFn, openPar, nT("1"), closePar}},
		{"round(2.12; 5)", []tokenData{roundFunc{}, openPar, nT("2.12"), semicolon, nT("5"), closePar}},
		{"inf - inf + inf", []tokenData{nT("inf"), minus, nT("inf"), plus, nT("inf")}},
		{"2 > 4", []tokenData{nT("2"), strictlyGreater, nT("4")}},
		{"2 < 4", []tokenData{nT("2"), strictlyLesser, nT("4")}},
		{"2 <= 4", []tokenData{nT("2"), lesser, nT("4")}},
		{"2 >= 4", []tokenData{nT("2"), greater, nT("4")}},
		{"2 == 4", []tokenData{nT("2"), equals, nT("4")}},
		{"2pi^2", []tokenData{nT("2"), mult, piConstant, pow, nT("2")}},
		// implicit multiplication
		{"2x", []tokenData{nT("2"), mult, x}},
		{"n!n", []tokenData{Variable{Name: 'n'}, factorial, mult, Variable{Name: 'n'}}},
		{"(3;)", []tokenData{openPar, nT("3"), semicolon, closePar}},
		{"{2;4;}", []tokenData{openCurly, nT("2"), semicolon, nT("4"), semicolon, closeCurly}},
		{"]-inf;inf[", []tokenData{closeBracket, minus, nT("inf"), semicolon, nT("inf"), openBracket}},
		// distinction between Variable and indice
		{"a_1", []tokenData{NewVarI('a', "1")}},
		{"a_{1}", []tokenData{NewVar('a'), underscore, openCurly, nT("1"), closeCurly}},
		// variable and }
		{"{s_1 }", []tokenData{openCurly, Variable{Name: 's', Indice: "1"}, closeCurly}},
		{"{s_1}", []tokenData{openCurly, Variable{Name: 's', Indice: "1"}, closeCurly}},
		// custom symbols
		{`"\ge + 5"`, []tokenData{Variable{Name: 0, Indice: `\ge + 5`}}},
		// factorial
		{`2*n!`, []tokenData{nT("2"), mult, Variable{Name: 'n'}, factorial}},
		// matrix
		{`[[1; 4; x]; [2; y; 4]]`, []tokenData{
			openMatrix, nT("1"), semicolon,
			nT("4"), semicolon, NewVar('x'), closeBracket, semicolon,
			openBracket, nT("2"), semicolon, NewVar('y'), semicolon, nT("4"), closeBracket, closeBracket,
		}},
		{`[ [ ] ]`, []tokenData{openMatrix, closeBracket, closeBracket}},
		{`[ 2 ] ]`, []tokenData{openBracket, nT("2"), closeBracket, closeBracket}},
		// replace x² by x^2 (same for x³, etc...)
		{
			"x¹ ; x² ; x³ ; x⁴ ; x⁵ ; x⁶ ; x⁷; x⁸; x⁹",
			[]tokenData{x, pow, nT("1"), semicolon, x, pow, nT("2"), semicolon, x, pow, nT("3"), semicolon, x, pow, nT("4"), semicolon, x, pow, nT("5"), semicolon, x, pow, nT("6"), semicolon, x, pow, nT("7"), semicolon, x, pow, nT("8"), semicolon, x, pow, nT("9")},
		},
	} {
		if got, _ := allTokens(test.expr); !reflect.DeepEqual(got, test.tokens) {
			t.Fatalf("for %s, expected %v, got %v", test.expr, test.tokens, got)
		}
	}
}

func TestTokenUnionTags(t *testing.T) {
	for _, tk := range []tokenData{
		symbol(0),
		nT(""),
		constant(0),
		Variable{},
		function(0),
		roundFunc{},
		specialFunctionKind(0),
		operator(0),
	} {
		tk.isToken()
	}
}

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
