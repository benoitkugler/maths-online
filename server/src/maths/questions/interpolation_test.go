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

func text(s string) textChunck { return textChunck{s, iText} }
func form(s string) textChunck { return textChunck{s, iFormula} }
func nb(s string) textChunck   { return textChunck{s, iNumberField} }

func TestInterpolated_parseFormula(t *testing.T) {
	tests := []struct {
		s       Interpolated
		wantOut []textChunck
	}{
		{
			"", []textChunck{text("")},
		},
		{
			"regular\ntwo lines", []textChunck{text("regular\ntwo lines")},
		},
		{
			"false$\n$$", []textChunck{text("false$\n$$")},
		},
		{
			"with formula\n $$ test $$", []textChunck{
				text("with formula"),
				form(" test "),
			},
		},
		{
			"with formula\n\n $$ test $$", []textChunck{
				text("with formula\n"),
				form(" test "),
			},
		},
		{
			"with formula and line\n $$ test $$\nother line", []textChunck{
				text("with formula and line"),
				form(" test "),
				text("other line"),
			},
		},
		{
			"with formula and line\n $$ test $$\nother line\n", []textChunck{
				text("with formula and line"),
				form(" test "),
				text("other line\n"),
			},
		},
		{
			"two formula\n $$ test $$ \n $$ test $$", []textChunck{
				text("two formula"),
				form(" test "),
				form(" test "),
			},
		},
		{
			"$$ # expr # $$\n#f#\na", []textChunck{
				form(" # expr # "),
				nb("f"),
				text("\na"),
			},
		},
	}
	for _, tt := range tests {
		if gotOut := tt.s.parseFormula(); !reflect.DeepEqual(gotOut, tt.wantOut) {
			t.Errorf("Interpolated.parseFormula() = %v, want %v", gotOut, tt.wantOut)
		}
	}
}

func Test_splitNumberField(t *testing.T) {
	tests := []struct {
		args    string
		wantOut []textChunck
	}{
		{"", nil},
		{"regular", []textChunck{text("regular")}},
		{"#expr#", []textChunck{nb("expr")}},
		{"#expr# ", []textChunck{nb("expr"), text(" ")}},
		{" #expr# aa #expr2#", []textChunck{text(" "), nb("expr"), text(" aa "), nb("expr2")}},
		{"invalid #", []textChunck{text("invalid #")}},
		{"#expr# invalid #", []textChunck{nb("expr"), text(" invalid #")}},
	}
	for _, tt := range tests {
		if gotOut := splitNumberField(tt.args); !reflect.DeepEqual(gotOut, tt.wantOut) {
			t.Errorf("splitNumberField() = %v, want %v", gotOut, tt.wantOut)
		}

	}
}
