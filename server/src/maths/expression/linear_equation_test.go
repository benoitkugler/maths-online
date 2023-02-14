package expression

import (
	"math"
	"reflect"
	"testing"
)

func TestExpr_isLinearEquation(t *testing.T) {
	tests := []struct {
		expr       string
		wantCoeffs linearCoefficients
		wantErr    bool
	}{
		{"x", linearCoefficients{NewVar('x'): 1}, false},
		{"2x - 2y +t -4", linearCoefficients{NewVar('x'): 2, NewVar('y'): -2, NewVar('t'): 1, {}: -4}, false},
		{"2x + y -1 * t -4", linearCoefficients{NewVar('x'): 2, NewVar('y'): 1, NewVar('t'): -1, {}: -4}, false},
		{"x/2 + (1/2)y", linearCoefficients{NewVar('x'): 0.5, NewVar('y'): 0.5}, false},
		{"x +0y", linearCoefficients{NewVar('x'): 1}, false},            // 0 handling
		{"5 -x + 4", linearCoefficients{NewVar('x'): -1, {}: 9}, false}, // duplicate constant term are accepted
		{"x*y", nil, true},
		{"x*y + z^2", nil, true},
		{"x + 2x", nil, true},
		{"sin(x)", nil, true},
	}

	for _, tt := range tests {
		e := mustParse(t, tt.expr)
		coeffs, err := e.isLinearEquation()
		if (err != nil) != tt.wantErr {
			t.Fatal(tt.expr)
		}
		if !reflect.DeepEqual(coeffs, tt.wantCoeffs) {
			t.Fatal(coeffs, tt.wantCoeffs)
		}
		err = e.IsValidLinearEquation(nil)
		if (err != nil) != tt.wantErr {
			t.Fatal()
		}
	}
}

func TestExpr_isLinearTerm(t *testing.T) {
	tests := []struct {
		expr         string
		wantCoeff    float64
		wantVariable Variable
		wantOk       bool
	}{
		{"2x", 2, NewVar('x'), true},
		{"2x_2", 2, NewVarI('x', "2"), true},
		{"-3x", -3, NewVar('x'), true},
		{"x*0.2", 0.2, NewVar('x'), true},
		{"y*0.2", 0.2, NewVar('y'), true},
		{"(-2)*(-z)", 2, NewVar('z'), true},
		{"-x/2", -0.5, NewVar('x'), true},
		{"x/(-2)", -0.5, NewVar('x'), true},
		{"x*sin(2)", math.Sin(2), NewVar('x'), true},
		{"2/x", 0, Variable{}, false},
		{"xy", 0, Variable{}, false},
		{"2xy", 0, Variable{}, false},
		{"x^2", 0, Variable{}, false},
		{"x*sin(y)", 0, Variable{}, false},
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)
		gotCoeff, gotVariable, gotOk := expr.isLinearTerm()
		if gotCoeff != tt.wantCoeff {
			t.Errorf("Expr.isLinearTerm() gotCoeff = %v, want %v", gotCoeff, tt.wantCoeff)
		}
		if gotVariable != tt.wantVariable {
			t.Errorf("Expr.isLinearTerm() gotVariable = %v, want %v", gotVariable, tt.wantVariable)
		}
		if gotOk != tt.wantOk {
			t.Errorf("Expr.isLinearTerm() gotOk = %v, want %v", gotOk, tt.wantOk)
		}
	}
}

func Test_linearCoefficients_isEquivalent(t *testing.T) {
	tests := []struct {
		c1, c2 linearCoefficients
		want   bool
	}{
		{
			linearCoefficients{NewVar('x'): 1, NewVar('y'): 1},
			linearCoefficients{NewVar('x'): 1, NewVar('y'): 1},
			true,
		},
		{
			linearCoefficients{NewVar('x'): 1, NewVar('y'): 1},
			linearCoefficients{NewVar('x'): 2, NewVar('y'): 2},
			true,
		},
		{
			linearCoefficients{NewVar('x'): 1, NewVar('y'): -1, Variable{}: 0.4},
			linearCoefficients{NewVar('x'): 2, NewVar('y'): -2, Variable{}: 0.8},
			true,
		},
		{
			linearCoefficients{NewVar('x'): 1, NewVar('y'): -1, Variable{}: 0.4},
			linearCoefficients{NewVar('x'): 2, NewVar('y'): -2},
			false,
		},
		{
			linearCoefficients{NewVar('x'): 1, NewVar('y'): -1},
			linearCoefficients{NewVar('x'): 2, NewVar('y'): 0.8},
			false,
		},
	}
	for _, tt := range tests {
		if got := tt.c1.isEquivalent(tt.c2); got != tt.want {
			t.Errorf("linearCoefficients.isEquivalent() = %v, want %v", got, tt.want)
		}
	}
}

func TestAreLinearEquationEquivalent(t *testing.T) {
	tests := []struct {
		e1, e2 string
		want   bool
	}{
		{"2x + y", "x+0.5y", true},
		{"3x - 5", "6x - 10", true},
		{"3x - 5", "6x - 4", false},
	}
	for _, tt := range tests {
		e1 := mustParse(t, tt.e1)
		e2 := mustParse(t, tt.e2)
		if got := AreLinearEquationsEquivalent(e1, e2); got != tt.want {
			t.Errorf("AreLinearEquationEquivalent() = %v, want %v", got, tt.want)
		}
	}
}
