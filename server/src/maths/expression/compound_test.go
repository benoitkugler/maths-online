package expression

import (
	"reflect"
	"testing"

	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
)

func TestCompoundUnionTag(t *testing.T) {
	for _, cp := range []Compound{
		Vector{},
		Set{},
		Interval{},
		&Expr{},
	} {
		cp.isCompound()
	}
}

func TestParseCompound(t *testing.T) {
	tests := []struct {
		expr    string
		want    Compound
		wantErr bool
	}{
		{ // invalid expression
			"(2 + a", nil, true,
		},
		{ // simple expression
			"2x + y", mustParse(t, "2x+y"), false,
		},
		{ // simple expression with starting par
			"(2x + y)(a+b)", mustParse(t, "(2x + y)(a+b)"), false,
		},
		{ // invalid vector
			"(2; a", nil, true,
		},
		{ // vector with simple entries
			"(2; a; b)", Vector{mustParse(t, "2"), mustParse(t, "a"), mustParse(t, "b")}, false,
		},
		{ // vector with complex entries
			"(2x+y; (a+b)(a-b); 1.2; (y))", Vector{mustParse(t, "2x+y"), mustParse(t, "(a+b)(a-b)"), mustParse(t, "1.2"), mustParse(t, "y")}, false,
		},
		{ // one-dimentional vector
			"(2;)", Vector{mustParse(t, "2")}, false,
		},
		{ // invalid set
			"{2; a", nil, true,
		},
		{ // set with simple entries
			"{2; a; b}", Set{mustParse(t, "2"), mustParse(t, "a"), mustParse(t, "b")}, false,
		},
		{ // set with complex entries
			"{2x+y; (a+b)(a-b); 1.2; (y)}", Set{mustParse(t, "2x+y"), mustParse(t, "(a+b)(a-b)"), mustParse(t, "1.2"), mustParse(t, "y")}, false,
		},
		{ // invalid interval
			"[2; a", nil, true,
		},
		{ // invalid interval
			"[2; a; c]", nil, true,
		},
		{ // interval with simple entries
			"]2; +inf]", Interval{Left: mustParse(t, "2"), Right: mustParse(t, "+inf"), LeftOpen: true, RightOpen: false}, false,
		},
		{ // interval with simple entries
			"[-inf; a[", Interval{Left: mustParse(t, "-inf"), Right: mustParse(t, "a"), LeftOpen: false, RightOpen: true}, false,
		},
		{ // interval with complex entries
			"]2x+y; (a+b)(a-b)[", Interval{Left: mustParse(t, "2x+y"), Right: mustParse(t, "(a+b)(a-b)"), LeftOpen: true, RightOpen: true}, false,
		},
		{
			// nested sets
			"{ {a; b} ; {c; d}}", nil, true,
		},
		{ // compatilibity between sets and indices
			"{ a_{1}; b_{2} }", Set{mustParse(t, "a_{1}"), mustParse(t, "b_{2}")}, false,
		},
		{ // sets and variables without spaces
			"{a_1;b_2}", Set{mustParse(t, "a_1"), mustParse(t, "b_2")}, false,
		},
	}
	for _, tt := range tests {
		got, err := ParseCompound(tt.expr)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseCompound(%s) error = %v, wantErr %v", tt.expr, err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("ParseCompound(%s) = %v, want %v", tt.expr, got, tt.want)
		}
	}
}

func TestCompoundLatex(t *testing.T) {
	var lines []string
	for _, expr := range []string{
		"{a+b; (a-b) / 3 }",
		"(2x ; y ; z/w)",
		"] -inf; 2x [",
		"] a/b; 3 ]",
	} {
		e, err := ParseCompound(expr)
		if err != nil {
			t.Fatal(err)
		}

		code := e.AsLaTeX()
		lines = append(lines, "$$"+code+"$$")
	}

	generateLatex(t, lines, "compounds.tex")
}

func TestCompoundSubstitute(t *testing.T) {
	vars := Vars{
		NewVar('x'): mustParse(t, "2"),
		NewVar('y'): mustParse(t, "c+d"),
		NewVar('a'): mustParse(t, "-b"),
	}
	for _, expr := range []string{
		"(x; y)",
		"{a+b; a-b}",
		"[0; x]",
		"]0; x[",
	} {
		cp, err := ParseCompound(expr)
		tu.AssertNoErr(t, err)

		cp.Substitute(vars)
		_ = cp.String()
	}
}

func TestAreCompoundsEquivalent(t *testing.T) {
	tests := []struct {
		e1, e2 string
		level  ComparisonLevel
		want   bool
	}{
		{"2+x", "x+2", SimpleSubstitutions, true},
		{"2", "(2;)", ExpandedSubstitutions, false},
		{"2", "{2}", ExpandedSubstitutions, false},
		{"{2}", "(2;)", ExpandedSubstitutions, false},
		{"[2;2]", "2", ExpandedSubstitutions, false},
		{"[2;2]", "(2;)", ExpandedSubstitutions, false},
		{"(2;)", "[2;2]", ExpandedSubstitutions, false},
		{"[2;2]", "{2}", ExpandedSubstitutions, false},
		{"(x+y; y+x)", "(x+y; x+y)", SimpleSubstitutions, true},
		{"(a; b)", "(b; a)", ExpandedSubstitutions, false},
		{"{a; b}", "{b; a}", Strict, true},
		{"{a; b}", "{b}", ExpandedSubstitutions, false},
		{"{b; b}", "{b}", Strict, true},
		{"(x+y; y+x)", "(x+y; x+y; x+y)", ExpandedSubstitutions, false},
		{"{x+y; y+x}", "{x+y; x+y; x+y}", SimpleSubstitutions, true},
		{"[2; 4]", "[2; 4[", ExpandedSubstitutions, false},
		{"[-inf; x+2]", "[-inf; 2+x]", SimpleSubstitutions, true},
	}
	for _, tt := range tests {
		c1, err := ParseCompound(tt.e1)
		tu.AssertNoErr(t, err)
		c2, err := ParseCompound(tt.e2)
		tu.AssertNoErr(t, err)

		if got := AreCompoundsEquivalent(c1, c2, tt.level); got != tt.want {
			t.Errorf("AreCompoundsEquivalent(%s, %s) = %v, want %v", tt.e1, tt.e2, got, tt.want)
		}
	}
}
