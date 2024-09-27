package expression

import (
	"math"
	"reflect"
	"testing"

	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

func TestRandomVariables_range(t *testing.T) {
	for range [10]int{} {
		rv := RandomParameters{defs: map[Variable]*Expr{
			NewVar('a'): mustParse(t, "3*randInt(1; 10)"),
			NewVar('b'): mustParse(t, "-a"),
		}}
		values, err := rv.Instantiate()
		if err != nil {
			t.Fatal(err)
		}
		if values[NewVar('a')].mustEvaluate(nil) != -values[NewVar('b')].mustEvaluate(nil) {
			t.Fatal(values)
		}

		if a := values[NewVar('a')].mustEvaluate(nil); a < 3 || a > 30 {
			t.Fatal(a)
		}

		rv = RandomParameters{defs: map[Variable]*Expr{
			NewVar('a'): mustParse(t, "randInt(1; 10)"),
			NewVar('b'): mustParse(t, "sgn(2*randInt(0;1)-1) * a"),
		}}
		values, err = rv.Instantiate()
		if err != nil {
			t.Fatal(err)
		}
		if a := values[NewVar('a')].mustEvaluate(nil); a < 1 || a > 10 {
			t.Fatal(a)
		}
		if a, b := values[NewVar('a')].mustEvaluate(nil), values[NewVar('b')].mustEvaluate(nil); math.Abs(a) != math.Abs(b) {
			t.Fatal(a, b)
		}
	}
}

func Test_sieveOfEratosthenes(t *testing.T) {
	tests := []struct {
		min, max   int
		wantPrimes []int
	}{
		{4, 4, nil},
		{0, 10, []int{2, 3, 5, 7}},
		{0, 11, []int{2, 3, 5, 7, 11}},
		{3, 10, []int{3, 5, 7}},
		{4, 11, []int{5, 7, 11}},
	}
	for _, tt := range tests {
		if gotPrimes := sieveOfEratosthenes(tt.min, tt.max); !reflect.DeepEqual(gotPrimes, tt.wantPrimes) {
			t.Errorf("sieveOfEratosthenes() = %v, want %v", gotPrimes, tt.wantPrimes)
		}
	}
}

func TestExpression_IsValidNumber(t *testing.T) {
	tests := []struct {
		expr           string
		parameters     RandomParameters
		checkPrecision bool
		wantErr        bool
	}{
		{
			"2a - sin(a) + exp(1 + a)", RandomParameters{defs: map[Variable]*Expr{NewVar('a'): mustParse(t, "2")}}, false, false,
		},
		{
			"2a + b", RandomParameters{defs: map[Variable]*Expr{NewVar('a'): mustParse(t, "2")}}, false, true,
		},
		{
			"1/0", RandomParameters{}, false, true,
		},
		{
			"1/a", RandomParameters{defs: map[Variable]*Expr{NewVar('a'): mustParse(t, "randInt(0;4)")}}, false, true,
		},
		{
			"1/a", RandomParameters{defs: map[Variable]*Expr{NewVar('a'): mustParse(t, "randInt(1;4)")}}, false, false,
		},
		{
			"1/a", RandomParameters{defs: map[Variable]*Expr{NewVar('a'): mustParse(t, "randDecDen()")}}, true, false,
		},
		{
			"(v_f - v_i) / v_i", RandomParameters{defs: map[Variable]*Expr{NewVarI('v', "f"): mustParse(t, "randint(1;10)"), NewVarI('v', "i"): mustParse(t, "randDecDen()")}}, true, false,
		},
		{
			"round(1/3; 3)", RandomParameters{}, true, false,
		},
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)

		allValid := true
		for range [100]int{} {
			vars, err := tt.parameters.Instantiate()
			if err != nil {
				t.Fatal(err)
			}

			err = expr.IsValidNumber(vars, tt.checkPrecision, true)
			allValid = allValid && err == nil
		}
		if !allValid != tt.wantErr {
			t.Errorf("Expression.IsValidNumber(%s) got = %v", tt.expr, !allValid)
		}
	}
}

func TestExpression_IsValidProba(t *testing.T) {
	tests := []struct {
		expr string
		vars Vars
		want bool
	}{
		{
			"1.1", nil, false,
		},
		{
			"1/a", Vars{NewVar('a'): mustParse(t, "4")}, true,
		},
		{
			"1/a", Vars{NewVar('a'): mustParse(t, "randDecDen()")}, true,
		},
		{
			"0.2 + 1/a", Vars{NewVar('a'): mustParse(t, "3")}, false,
		},
		{
			"round(1/3; 3)", nil, true,
		},
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)
		err := expr.IsValidProba(tt.vars)
		if (err == nil) != tt.want {
			t.Errorf("Expression.IsValidProba(%s) got = %v, want %v", tt.expr, err, tt.want)
		}
	}
}

func mustParseMany(t *testing.T, exprs []string) []*Expr {
	out := make([]*Expr, len(exprs))
	for i, s := range exprs {
		out[i] = mustParse(t, s)
	}
	return out
}

func TestExpression_AreSortedNumbers(t *testing.T) {
	tests := []struct {
		exprs []string
		vars  Vars
		want  bool
	}{
		{
			[]string{"1", "2", "a"}, Vars{NewVar('a'): mustParse(t, "3")}, true,
		},
		{
			[]string{"1", "2", "a"}, Vars{NewVar('a'): mustParse(t, "1")}, false,
		},
		{
			[]string{"1", "2", "b"}, Vars{NewVar('a'): mustParse(t, "3")}, false,
		},
	}
	for _, tt := range tests {
		exprs := mustParseMany(t, tt.exprs)

		err := AreSortedNumbers(exprs, tt.vars)
		if (err == nil) != tt.want {
			t.Errorf("AreSortedNumbers(%s) got = %v, want %v", tt.exprs, err, tt.want)
		}
	}
}

func TestExpression_IsValidIndex(t *testing.T) {
	tests := []struct {
		expr   string
		vars   Vars
		length int
		want   bool
	}{
		{
			"+1 + 1 * (a==1) + 2 * (a==2) + 3*(a==3)", Vars{NewVar('a'): mustParse(t, "2")}, 4, true,
		},
		{
			"+1 + 1 * (a==1) + 2 * (a==2) + 3*(a==3)", Vars{NewVar('a'): mustParse(t, "3")}, 4, true,
		},
		{
			"+1 + 1 * (a==1) + 2 * (a==2) + 2.5*(a==3)", Vars{NewVar('a'): mustParse(t, "3")}, 4, false,
		},
		{
			"+1 + 1 * (a==1) + 2 * (a==2) + 4*(a==3)", Vars{NewVar('a'): mustParse(t, "3")}, 4, false,
		},
		{
			"+1 + 1 * (a^2 - b^2 - c^2 == 0) + 2*(b^2 - a^2 - c^2 == 0) + 3*(c^2 - a^2 - b^2 == 0)", Vars{
				NewVar('a'): mustParse(t, "4"), // BC
				NewVar('b'): mustParse(t, "5"), // AC
				NewVar('c'): mustParse(t, "6"), // AB
			}, 4, true,
		},
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)
		err := expr.IsValidIndex(tt.vars, tt.length)
		if (err == nil) != tt.want {
			t.Errorf("Expression.IsValidIndex() got = %v, want %v", err, tt.want)
		}
	}
}

func TestFunctionDefinition_IsValid(t *testing.T) {
	tests := []struct {
		expr     string
		variable rune
		vars     Vars
		from, to string
		bound    float64 // expected extrema
		want     bool
	}{
		{"2x + 1", 'x', nil, "-10", "10", 25, true},
		{"2x + 1", 'x', nil, "-10", "10", 10, false},
		{"2x + 1", 'x', nil, "2", "2", 10, false},
		{"2x + a", 'x', nil, "-10", "10", 10, false},
		{"1/x", 'x', nil, "-10", "10", 100, false},
		{"exp(x)", 'x', nil, "-10", "10", 100, false},
		{"ax + b", 'x', Vars{
			NewVar('a'): mustParse(t, "3"),
			NewVar('b'): mustParse(t, "7"),
		}, "-10", "10", 100, true},
		{"ax + b", 'x', Vars{
			NewVar('a'): mustParse(t, "90"),
			NewVar('b'): mustParse(t, "7"),
		}, "-10", "10", 100, false},
	}
	for _, tt := range tests {
		expr := mustParse(t, tt.expr)
		fn := FunctionExpr{
			Function: expr,
			Variable: NewVar(tt.variable),
		}
		err := fn.IsValidAsFunction(Domain{mustParse(t, tt.from), mustParse(t, tt.to)}, tt.vars, tt.bound)
		if (err == nil) != tt.want {
			t.Errorf("Expression.AreFxsIntegers() got = %v, want %v", err, tt.want)
		}
	}
}

func mustParseDomains(t *testing.T, domains [][2]string) []Domain {
	out := make([]Domain, len(domains))
	for i, d := range domains {
		out[i] = Domain{mustParse(t, d[0]), mustParse(t, d[1])}
	}
	return out
}

func TestAreDisjointsDomains(t *testing.T) {
	tests := []struct {
		domains [][2]string
		vars    Vars
		wantErr bool
	}{
		{
			[][2]string{{"0", "1"}}, nil, false,
		},
		{
			[][2]string{{"0", "1"}, {"1", "2"}}, nil, false,
		},
		{
			[][2]string{{"0", "1"}, {"0", "2"}}, nil, true,
		},
		{
			[][2]string{{"0", "1"}, {"x", "0.5"}}, Vars{NewVar('x'): newNb(0)}, true,
		},
		{ // evaluation error
			[][2]string{{"0", "1"}, {"x", "0.5"}}, nil, true,
		},
	}
	for _, tt := range tests {
		domains := mustParseDomains(t, tt.domains)
		if err := AreDisjointsDomains(domains, tt.vars); (err != nil) != tt.wantErr {
			t.Errorf("AreDisjointsDomains() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
}

func TestDomain_IsIncludedIntoOne(t *testing.T) {
	tests := []struct {
		From    string
		To      string
		domains [][2]string
		vars    Vars
		wantErr bool
	}{
		{
			"1", "2", [][2]string{{"0", "1"}, {"1", "2"}}, nil, false,
		},
		{ // infinity
			"1", "2", [][2]string{{"", ""}, {"2", "2"}}, nil, false,
		},
		{
			"1", "2", [][2]string{{"0", "1"}, {"1.5", "2"}}, nil, true,
		},
		{
			"x", "2", [][2]string{{"0", "1"}, {"1.5", "2"}}, Vars{NewVar('x'): newNb(0)}, true,
		},
		{
			"x", "2", [][2]string{{"0", "1"}, {"x", "3"}}, Vars{NewVar('x'): newNb(-1)}, false,
		},
		{ // evaluation error
			"x", "2", [][2]string{{"0", "1"}, {"1", "2"}}, nil, true,
		},
		{ // evaluation error
			"1", "2", [][2]string{{"x", "1"}, {"x", "2"}}, nil, true,
		},
	}
	for _, tt := range tests {
		d := Domain{
			From: mustParse(t, tt.From),
			To:   mustParse(t, tt.To),
		}
		domains := mustParseDomains(t, tt.domains)
		if err := d.IsIncludedIntoOne(domains, tt.vars); (err != nil) != tt.wantErr {
			t.Errorf("Domain.IsIncludedIntoOne() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
}

func TestExpr_IsIncludedIntoOne(t *testing.T) {
	tests := []struct {
		expr    string
		domains [][2]string
		vars    Vars
		wantErr bool
	}{
		{
			"1", [][2]string{{"0", "1"}, {"1", "2"}}, nil, false,
		},
		{ // infinity
			"1", [][2]string{{"-2", "-1"}, {"0", ""}}, nil, false,
		},
		{
			"1", [][2]string{{"0", "0.8"}, {"1.5", "2"}}, nil, true,
		},
		{
			"x", [][2]string{{"0", "1"}, {"1.5", "2"}}, Vars{NewVar('x'): newNb(-1)}, true,
		},
		{
			"x", [][2]string{{"0", "1"}, {"x-1", "3"}}, Vars{NewVar('x'): newNb(-1)}, false,
		},
		{ // evaluation error
			"x", [][2]string{{"0", "1"}, {"1", "2"}}, nil, true,
		},
		{ // evaluation error
			"1", [][2]string{{"x", "1"}, {"x", "2"}}, nil, true,
		},
	}
	for _, tt := range tests {
		e := mustParse(t, tt.expr)
		domains := mustParseDomains(t, tt.domains)
		if err := e.IsIncludedIntoOne(domains, tt.vars); (err != nil) != tt.wantErr {
			t.Errorf("Domain.IsIncludedIntoOne() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
}

func Test_binomialCoefficient(t *testing.T) {
	tests := []struct {
		k    int
		n    int
		want int
	}{
		{0, 4, 1},
		{4, 4, 1},
		{1, 4, 4},
		{3, 4, 4},
		{-3, 4, 0},
		{2, 4, 6},
		{2, 3, 3},
		{2, 5, 10},
	}
	for _, tt := range tests {
		if got := binomialCoefficient(tt.k, tt.n); got != tt.want {
			t.Errorf("binomialCoefficient() = %v, want %v", got, tt.want)
		}
	}
}

func TestDecDenominator(t *testing.T) {
	tu.Assert(t, reflect.DeepEqual(generateDecDenominator(33, 33), []int{}))
	tu.Assert(t, reflect.DeepEqual(generateDecDenominator(32, 32), []int{32}))
	tu.Assert(t, reflect.DeepEqual(generateDecDenominator(30, 40), []int{40, 32}))
	tu.Assert(t, reflect.DeepEqual(generateDecDenominator(20, 24), []int{20}))
	tu.Assert(t, reflect.DeepEqual(generateDecDenominator(1, 100), []int{1, 5, 25, 2, 10, 50, 4, 20, 100, 8, 40, 16, 80, 32, 64}))
}
