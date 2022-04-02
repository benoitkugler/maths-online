package editor

import (
	"reflect"
	"testing"

	ex "github.com/benoitkugler/maths-online/maths/exercice"
)

const regularText = "This is a regular text"

func TestParseInterpolatedString(t *testing.T) {
	tests := []struct {
		args    string
		want    ex.TextParts
		wantErr bool
	}{
		{
			regularText, ex.TextParts{{Content: regularText}}, false,
		},
		{
			"Regular $Latex$ #{2x +1} regular end", ex.TextParts{
				{Content: "Regular ", Kind: ex.Text},
				{Content: "Latex", Kind: ex.StaticMath},
				{Content: " ", Kind: ex.Text},
				{Content: "2x +1", Kind: ex.Expression},
				{Content: " regular end", Kind: ex.Text},
			}, false,
		},
		{ // we accept unterminated sequences
			"$ 45", ex.TextParts{{Content: " 45", Kind: ex.StaticMath}}, false,
		},
		{ // we accept unterminated sequences
			"#{45x ", ex.TextParts{{Content: "45x ", Kind: ex.Expression}}, false,
		},
		{
			"#{45x - +}", nil, true,
		},
	}
	for _, tt := range tests {
		got, err := ParseInterpolatedString(tt.args)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseInterpolatedString() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("ParseInterpolatedString() = %v, want %v", got, tt.want)
		}
	}
}

func TestInterpolationRoundtrip(t *testing.T) {
	for _, s := range []string{
		regularText,
		"Regular $Latex$ #{2x +1} regular end",
		"$45$$45687$ #{78}",
	} {
		got, err := ParseInterpolatedString(s)
		if err != nil {
			t.Fatal(s)
		}
		if s2 := NewTextBlock(ex.TextBlock{Parts: got}).Parts; s2 != s {
			t.Fatalf("invalid roundtrip: %s != %s", s2, s)
		}
	}
}
