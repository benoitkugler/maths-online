package exercice

import (
	"reflect"
	"testing"

	"github.com/benoitkugler/maths-online/maths/exercice/client"
	"github.com/benoitkugler/maths-online/maths/expression"
)

func TestFieldInstance_validateAnswerSyntax(t *testing.T) {
	tests := []struct {
		field   fieldInstance
		args    client.Answer
		wantErr bool
	}{
		{NumberFieldInstance{}, client.NumberAnswer{}, false},
		{NumberFieldInstance{}, client.RadioAnswer{}, true},
		{ExpressionFieldInstance{}, client.RadioAnswer{}, true},
		{ExpressionFieldInstance{}, client.ExpressionAnswer{Expression: ""}, true},
		{ExpressionFieldInstance{}, client.ExpressionAnswer{Expression: "2+4"}, false},
	}
	for _, tt := range tests {
		err := tt.field.validateAnswerSyntax(tt.args)
		if (err != nil) != tt.wantErr {
			t.Fatalf("NumberFieldInstance.validateAnswerSyntax() error = %v, wantErr %v", err, tt.wantErr)
		}
		if err != nil {
			_ = err.(InvalidFieldAnswer).Error()
		}
	}
}

func TestNumberFieldInstance_evaluateAnswer(t *testing.T) {
	tests := []struct {
		field         fieldInstance
		args          client.Answer
		wantIsCorrect bool
	}{
		{NumberFieldInstance{Answer: 1}, client.NumberAnswer{Value: 1}, true},
		{NumberFieldInstance{Answer: 1}, client.NumberAnswer{Value: 1.1}, false},
		{ExpressionFieldInstance{Answer: mustParse("x+2"), ComparisonLevel: expression.SimpleSubstitutions}, client.ExpressionAnswer{Expression: "x + 2"}, true},
		{ExpressionFieldInstance{Answer: mustParse("x+2"), ComparisonLevel: expression.SimpleSubstitutions}, client.ExpressionAnswer{Expression: "2+x "}, true},
		{ExpressionFieldInstance{Answer: mustParse("x+2"), ComparisonLevel: expression.SimpleSubstitutions}, client.ExpressionAnswer{Expression: "2+ 1*x "}, true},
	}
	for _, tt := range tests {
		f := tt.field
		gotIsCorrect := f.evaluateAnswer(tt.args)
		if gotIsCorrect != tt.wantIsCorrect {
			t.Errorf("NumberFieldInstance.evaluateAnswer() = %v, want %v", gotIsCorrect, tt.wantIsCorrect)
		}
	}
}

func TestOrderedList(t *testing.T) {
	field := OrderedListFieldInstance{
		Answer: []StringOrExpression{ // [12;+infty]
			{String: "["},
			{Expression: mustParse("12")},
			{String: ";"},
			{String: "+"},
			{String: `\infty`},
			{String: `]`},
		},
		AdditionalProposals: []StringOrExpression{
			{String: "]"}, // some duplicates
			{String: `\infty`},
			{Expression: mustParse("11")},
			{String: "-"},
		},
	}

	p1 := field.proposals()
	p2 := field.proposals()
	if !reflect.DeepEqual(p1, p2) {
		t.Fatal(p1, p2)
	}
}