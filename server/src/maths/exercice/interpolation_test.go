package exercice

import (
	"reflect"
	"testing"
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
