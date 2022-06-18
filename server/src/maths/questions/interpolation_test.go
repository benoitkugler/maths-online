package questions

import (
	"reflect"
	"testing"

	"github.com/benoitkugler/maths-online/maths/expression"
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
	}
	for _, tt := range tests {
		got, err := tt.args.Parse()
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
	parts, err := s.Parse()
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
