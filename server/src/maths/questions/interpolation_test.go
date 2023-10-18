package questions

import (
	"reflect"
	"testing"

	"github.com/benoitkugler/maths-online/server/src/maths/expression"
)

const regularText = "This is a regular text"

func TestParseInterpolatedString(t *testing.T) {
	tests := []struct {
		args    Interpolated
		want    TextParts
		wantErr bool
	}{
		{
			regularText, TextParts{{Content: regularText}}, false,
		},
		{
			"Regular $Latex$ &2x +1& regular end", TextParts{
				{Content: "Regular ", Kind: Text},
				{Content: "Latex", Kind: StaticMath},
				{Content: " ", Kind: Text},
				{Content: "2x +1", Kind: Expression},
				{Content: " regular end", Kind: Text},
			}, false,
		},
		{ // expression inside latex
			"Regular $A = &a&$", TextParts{
				{Content: "Regular ", Kind: Text},
				{Content: "A = ", Kind: StaticMath},
				{Content: "a", Kind: Expression},
			}, false,
		},
		{
			"&45x - +&", nil, true,
		},
		// expressions with compound
		{
			"&{x; y}&&]2; 10[&", TextParts{
				{Content: "{x; y}", Kind: Expression},
				{Content: "]2; 10[", Kind: Expression},
			}, false,
		},
	}
	for _, tt := range tests {
		got, err := tt.args.parse()
		if (err != nil) != tt.wantErr {
			t.Errorf("Interpolated.Parse() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Interpolated.Parse() = %v, want %v", got, tt.want)
		}
	}
}

func Test_splitByLaTeX(t *testing.T) {
	tests := []struct {
		args    string
		wantOut []TextPart
	}{
		{"Plain text", []TextPart{{Content: "Plain text"}}},
		{"$L$$L$", []TextPart{{Content: "L", Kind: StaticMath}, {Content: "L", Kind: StaticMath}}},
		{"Plain text $2x +3$", []TextPart{{Content: "Plain text "}, {Content: "2x +3", Kind: StaticMath}}},
		{"Plain text $2x +3$", []TextPart{{Content: "Plain text "}, {Content: "2x +3", Kind: StaticMath}}},
		{"Plain text $&2x +3&$ end", []TextPart{{Content: "Plain text "}, {Content: "&2x +3&", Kind: StaticMath}, {Content: " end"}}},
	}
	for _, tt := range tests {
		if gotOut := splitByLaTeX(tt.args); !reflect.DeepEqual(gotOut, tt.wantOut) {
			t.Errorf("splitByLaTeX() = %v, want %v", gotOut, tt.wantOut)
		}
	}
}

func Test_splitByExpression(t *testing.T) {
	tests := []struct {
		args TextPart
		want []TextPart
	}{
		{
			TextPart{Content: "mlqk "}, []TextPart{{Content: "mlqk "}},
		},
		{
			TextPart{Content: "mlqk &lsd& smdl"}, []TextPart{{Content: "mlqk "}, {Content: "lsd", Kind: Expression}, {Content: " smdl"}},
		},
		{
			TextPart{Content: "mlqk &lsd&", Kind: StaticMath}, []TextPart{{Content: "mlqk ", Kind: StaticMath}, {Content: "lsd", Kind: Expression}},
		},
	}
	for _, tt := range tests {
		if got := splitByExpression(tt.args); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("splitByExpression() = %v, want %v", got, tt.want)
		}
	}
}

func TestLatexOutput(t *testing.T) {
	s := Interpolated(`$\overset{\rightarrow}{ &P& &P& }$`)
	parts, err := s.parse()
	if err != nil {
		t.Fatal(err)
	}
	sample, err := parts.instantiate(expression.Vars{expression.NewVar('P'): expression.NewVarExpr(expression.NewVar('B'))})
	if err != nil {
		t.Fatal(err)
	}
	if len(sample) != 1 {
		t.Fatal(len(sample))
	}
}

func TestInterpolated_parseFormula(t *testing.T) {
	tests := []struct {
		s       Interpolated
		wantOut []textOrFormula
	}{
		{
			"", []textOrFormula{{isFormula: false}},
		},
		{
			"regular\ntwo lines", []textOrFormula{{"regular\ntwo lines", false}},
		},
		{
			"false$\n$$", []textOrFormula{{"false$\n$$", false}},
		},
		{
			"with formula\n $$ test $$", []textOrFormula{
				{"with formula", false},
				{" test ", true},
			},
		},
		{
			"with formula\n\n $$ test $$", []textOrFormula{
				{"with formula\n", false},
				{" test ", true},
			},
		},
		{
			"with formula and line\n $$ test $$\nother line", []textOrFormula{
				{"with formula and line", false},
				{" test ", true},
				{"other line", false},
			},
		},
		{
			"with formula and line\n $$ test $$\nother line\n", []textOrFormula{
				{"with formula and line", false},
				{" test ", true},
				{"other line\n", false},
			},
		},
		{
			"two formula\n $$ test $$ \n $$ test $$", []textOrFormula{
				{"two formula", false},
				{" test ", true},
				{" test ", true},
			},
		},
	}
	for _, tt := range tests {
		if gotOut := tt.s.parseFormula(); !reflect.DeepEqual(gotOut, tt.wantOut) {
			t.Errorf("Interpolated.parseFormula() = %v, want %v", gotOut, tt.wantOut)
		}
	}
}
