package exercice

import (
	"reflect"
	"testing"
)

const regularText = "This is a regular text"

func TestParseInterpolatedString(t *testing.T) {
	tests := []struct {
		args    string
		want    TextParts
		wantErr bool
	}{
		{
			regularText, TextParts{{Content: regularText}}, false,
		},
		{
			"Regular $Latex$ #{2x +1} regular end", TextParts{
				{Content: "Regular ", Kind: Text},
				{Content: "Latex", Kind: StaticMath},
				{Content: " ", Kind: Text},
				{Content: "2x +1", Kind: Expression},
				{Content: " regular end", Kind: Text},
			}, false,
		},
		{ // we accept unterminated sequences
			"$ 45", TextParts{{Content: " 45", Kind: StaticMath}}, false,
		},
		{ // we accept unterminated sequences
			"#{45x ", TextParts{{Content: "45x ", Kind: Expression}}, false,
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
