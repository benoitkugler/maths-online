package expression

import (
	"reflect"
	"testing"
)

func mustParse(t *testing.T, s string) *Expression {
	t.Helper()

	want, _, err := Parse(s)
	if err != nil {
		t.Fatal(s, err)
	}
	return want
}

func Test_Expression_extractOperator(t *testing.T) {
	tests := []struct {
		expr string
		op   operator
		want []*Expression
	}{
		{"1 +  2 + ( 3 + 4)", plus, []*Expression{
			newNumber(1),
			newNumber(2),
			newNumber(3),
			newNumber(4),
		}},
		{"1 + (2 * x) + 3", plus, []*Expression{
			newNumber(1),
			{atom: mult, left: newNumber(2), right: newVariable('x')},
			newNumber(3),
		}},
		{"1 + 2 * (x + y) + 3", plus, []*Expression{
			newNumber(1),
			{
				atom:  mult,
				left:  newNumber(2),
				right: &Expression{atom: plus, left: newVariable('x'), right: newVariable('y')},
			},
			newNumber(3),
		}},
		{"1 *  2 * ( 3 * 4)", mult, []*Expression{
			newNumber(1),
			newNumber(2),
			newNumber(3),
			newNumber(4),
		}},
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)
		if got := expr.extractOperator(tt.op); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("node.extractPlus() = %v, want %v", got, tt.want)
		}
	}
}

func Test_Expression_sortPlusAndMultOperands(t *testing.T) {
	tests := []struct {
		expr string
		want string
	}{
		{"1 +  2 +3", "1+2+3"}, // no op
		{"1 + 3 + 2", "1+2+3"},
		{" + 3 + 2", "2+3"},
		{"1 + x + 4", "1+4+x"},
		{"(-1) * (2-4)", "(-1) * (2-4)"},
		{"(2-4)  * (-1)", "(-1) * (2-4)"},
		{"(-4)  * (-1)", "(-1) * (-4)"},
		{"1 + x + 4 + a", "1+4+a+x"},
		{"1 + x + (2 * y) + (2 - y)", "1+x + (2 - y) + (2*y)"},
		{"1 + x + (3 * y) + (2 * y)", "1+x + (2 * y) + (3*y)"},
		{"exp(5) + 1 + log(10) + x + sin(2) + sin(1)", "1+x+log(10)+exp(5)+sin(1)+sin(2)"},
		{"1 + 2 * (x+1) + y", "1+y+2*(1+x)"},
		{"1 + e + \u03C0", "1+\u03C0+e"},
		{"1 + \u03C0 + e + e ", "1+\u03C0+e+e"},
		{"1 + (2*x) + (2/x)", "1 + (2*x) + (2/x)"},
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)
		expr.sortPlusAndMultOperands()

		want := mustParse(t, tt.want)
		if !reflect.DeepEqual(expr, want) {
			t.Fatalf("sortPlusAndMultOperands(%s) = %v, expected %v", tt.expr, expr, tt.want)
		}
	}
}

func Test_Expression_expandMult(t *testing.T) {
	tests := []struct {
		expr string
		want string
	}{
		{"1 + 2", "1+2"}, // no op
		{"2*(3+2)", "2*3 + 2*2"},
		{"(3+2)*2", "3 * 2+ 2*2"},
		{"2*( 3 * (x + y) +2)", "2*(3*x + 3*y) + 2*2"}, // only one level of expansion
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)
		expr.expandMult()

		want := mustParse(t, tt.want)
		if !reflect.DeepEqual(expr, want) {
			t.Fatalf("expand(%s) = %v, expected %v", tt.expr, expr, tt.want)
		}
	}
}

func Test_Expression_expandPow(t *testing.T) {
	tests := []struct {
		expr string
		want string
	}{
		{"1 + 2", "1+2"}, // no op
		{"2^3", "2*2*2"},
		{"2^3.1", "2^3.1"},
		{"2^x", "2^x"},
		{"2^(-1)", "2^(-1)"},
		{"2^1", "2"},
		{"(a+b)^2", "(a+b)*(a+b)"},
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)
		expr.expandPow()

		want := mustParse(t, tt.want)
		if !reflect.DeepEqual(expr, want) {
			t.Fatalf("expand(%s) = %v, expected %v", tt.expr, expr, tt.want)
		}
	}
}

func Test_Expression_groupAdditions(t *testing.T) {
	tests := []struct {
		expr string
		want string
	}{
		{"1+x+y", "1+x+y"}, // no-op
		{"x+x", "2*x"},
		{"x+x+x", "3*x"},
		{"2*(x+x) + 2*(x+x) + 2*(x+x)", "3*(2*(2*x))"},
		{"x+x+x + y + y", "3*x + 2*y"},
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)
		expr.groupAdditions()

		want := mustParse(t, tt.want)
		if !reflect.DeepEqual(expr, want) {
			t.Fatalf("groupAdditions(%s) = %v, expected %v", tt.expr, expr, tt.want)
		}
	}
}

func Test_Expression_expandMinus(t *testing.T) {
	tests := []struct {
		expr string
		want string
	}{
		{"x - y", "x + (-y)"},
		{"x + -y", "x + (-y)"},
		{"x - (y - z)", "x + (- (y + (-z))) "},
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)
		expr.expandMinus()

		want := mustParse(t, tt.want)
		if !reflect.DeepEqual(expr, want) {
			t.Fatalf("expandMinus(%s) = %v, expected %v", tt.expr, expr, tt.want)
		}
	}
}

func Test_Expression_fullSimplification(t *testing.T) {
	tests := []struct {
		expr string
		want string
	}{
		{"1+x+y", "1+x+y"}, // no-op
		{"2*(2 + x)", "4 + 2*x"},
		{"(2+z)*(2 + x)", "4 + 2*x + 2*z + x*z"},
		{"(2+x)*(2 + z)", "4 + 2*x + 2*z + x*z"},
		{"(a+b)^2", "a*a + b*b + 2*a*b"},
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)
		nbPasses := expr.fullSimplification()

		t.Logf("Required %d passes", nbPasses)

		want := mustParse(t, tt.want)
		if !reflect.DeepEqual(expr, want) {
			t.Fatalf("fullSimplification(%s) = %v, expected %v", tt.expr, expr, tt.want)
		}
	}
}

func Test_Expression_basicSimplification(t *testing.T) {
	tests := []struct {
		expr string
		want string
	}{
		{"1+x+y", "1+x+y"}, // no-op
		{"x+2", "2+x"},
		{"2+3+x", "5+x"},
		{"2*(x+ 2)", "2*(2+x)"},
		{"2 - y", "2 + (-y)"},
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)
		expr.basicSimplification()

		want := mustParse(t, tt.want)
		if !reflect.DeepEqual(expr, want) {
			t.Fatalf("basicSimplification(%s) = %v, expected %v", tt.expr, expr, tt.want)
		}
	}
}

func Test_AreExpressionEquivalent(t *testing.T) {
	tests := []struct {
		e1, e2 string
		level  ComparisonLevel
		want   bool
	}{
		{"1+x", "x+1", Strict, false},
		{"1+x", "x+1", SimpleSubstitutions, true},
		{"1+x", "x+1", ExpandedSubstitutions, true},
		{"(a+b)^2", "a^2 + 2*a*b + b^2", SimpleSubstitutions, false},
		{"(a+b)^2", "a^2 + 2*a*b + b^2", ExpandedSubstitutions, true},
		{"(a+b)^2", "a^2 + b^2 + 2*a*b ", ExpandedSubstitutions, true},
		{"a + y", "a", Strict, false},
		{"a + y", "a", SimpleSubstitutions, false},
		{"a + y", "a", ExpandedSubstitutions, false},
		{"a + y", "a", ExpandedSubstitutions, false},
		{"x + x", "2*x", ExpandedSubstitutions, true},
		{"x + (- y) + x", "2*x - y", ExpandedSubstitutions, true},
		{"x + (- y) + x", "2*x + (- y)", ExpandedSubstitutions, true},
		// {"x + x - x", "x", ExpandedSubstitutions, true}, // TODO ?
		{"e^x * (2*x +1)", "2*x*e^x + e^x", ExpandedSubstitutions, true},
		{"e^x * (2*x +1)", "2*x*e^x + e^x", SimpleSubstitutions, false},
		{"2 + randInt(1,3)", "randInt(1,3) + 2", SimpleSubstitutions, true},
		{"2 + randInt(1,3) + randInt(2,4)", "randInt(1,3) + 2 + randInt(2,4)", SimpleSubstitutions, true},
		{"2 + randInt(1,3) + randInt(1,4)", "randInt(1,3) + 2 + randInt(1,4)", SimpleSubstitutions, true},
		{"2 + randPrime(1,3) + randPrime(2,4)", "randPrime(1,3) + 2 + randPrime(2,4)", SimpleSubstitutions, true},
		{"2 + randPrime(1,3) + randPrime(1,4)", "randPrime(1,3) + 2 + randPrime(1,4)", SimpleSubstitutions, true},
		{"randPrime(1,4) + randInt(1,3)", "randInt(1,3) + randPrime(1,4)", SimpleSubstitutions, true},
		{"-1.2", "-(1.2)", SimpleSubstitutions, true},
	}
	for _, tt := range tests {
		e1, e2 := mustParse(t, tt.e1), mustParse(t, tt.e2)
		if got := AreExpressionsEquivalent(e1, e2, tt.level); got != tt.want {
			t.Fatalf("AreExpressionEquivalent(%s, %s) = %v, want %v", tt.e1, tt.e2, got, tt.want)
		}
	}
}