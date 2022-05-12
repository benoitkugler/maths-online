package exercice

import (
	"reflect"
	"testing"

	"github.com/benoitkugler/maths-online/maths/exercice/client"
	"github.com/benoitkugler/maths-online/maths/expression"
	"github.com/benoitkugler/maths-online/maths/repere"
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
		{RadioFieldInstance{Answer: 2}, client.RadioAnswer{Index: 1}, true},
		{ExpressionFieldInstance{Answer: expression.MustParse("x+2"), ComparisonLevel: expression.SimpleSubstitutions}, client.ExpressionAnswer{Expression: "x + 2"}, true},
		{ExpressionFieldInstance{Answer: expression.MustParse("x+2"), ComparisonLevel: expression.SimpleSubstitutions}, client.ExpressionAnswer{Expression: "2+x "}, true},
		{ExpressionFieldInstance{Answer: expression.MustParse("x+2"), ComparisonLevel: expression.ExpandedSubstitutions}, client.ExpressionAnswer{Expression: "2+ 1*x "}, true},
	}
	for _, tt := range tests {
		f := tt.field
		gotIsCorrect := f.evaluateAnswer(tt.args)
		if gotIsCorrect != tt.wantIsCorrect {
			t.Errorf("NumberFieldInstance.evaluateAnswer(%v) = %v, want %v", tt.args, gotIsCorrect, tt.wantIsCorrect)
		}
	}
}

func TestOrderedList(t *testing.T) {
	field := OrderedListFieldInstance{
		Answer: []string{ // [12;+infty]
			"[",
			"12",
			";",
			"+",
			`\infty`,
			`]`,
		},
		AdditionalProposals: []string{
			"]", // some duplicates
			`\infty`,
			"11",
			"-",
		},
	}

	p1 := field.proposals()
	p2 := field.proposals()
	if !reflect.DeepEqual(p1, p2) {
		t.Fatal(p1, p2)
	}

	// check that the correct answer is indeed correct
	ans := field.correctAnswer()

	if !field.evaluateAnswer(ans) {
		t.Fatalf("invalid answer %v (proposals : %v)", ans, p1)
	}
}

func TestFigureAffineLineField(t *testing.T) {
	field := FigureAffineLineFieldInstance{
		Figure: repere.Figure{
			Bounds: repere.RepereBounds{
				Width: 20, Height: 20,
				Origin: repere.Coord{4, 4},
			},
		},
		AnswerA: 1 / 3,
		AnswerB: +3,
	}
	if ans := field.correctAnswer(); !field.evaluateAnswer(ans) {
		t.Fatalf("inconsistent answer %v", ans)
	}
}
