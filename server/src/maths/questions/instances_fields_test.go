package questions

import (
	"encoding/json"
	"math"
	"os"
	"reflect"
	"testing"

	"github.com/benoitkugler/maths-online/maths/expression"
	"github.com/benoitkugler/maths-online/maths/questions/client"
	"github.com/benoitkugler/maths-online/maths/repere"
	"github.com/benoitkugler/maths-online/utils"
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
		{RadioFieldInstance{Answer: 2, Proposals: []client.TextLine{{}, {}}}, client.RadioAnswer{Index: 0}, true}, // the answer is swapped
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
		Answer: []client.TextLine{ // [12;+infty]
			{{Text: "["}},
			{{Text: "12"}},
			{{Text: ";"}},
			{{Text: "+"}},
			{{Text: `\infty`}},
			{{Text: `]`}},
		},
		AdditionalProposals: []client.TextLine{
			{{Text: "]"}}, // some duplicates
			{{Text: `\infty`}},
			{{Text: "11"}},
			{{Text: "-"}},
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

func TestRadio(t *testing.T) {
	field := RadioFieldInstance{
		Proposals: []client.TextLine{ // [12;+infty]
			{{Text: "["}},
			{{Text: "12"}},
			{{Text: ";"}},
			{{Text: "+"}},
			{{Text: `\infty`}},
			{{Text: `]`}},
		},
		Answer: 6,
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

func TestBug72(t *testing.T) {
	field := DropDownFieldInstance{
		Proposals: []client.TextLine{
			{{Text: "GB", IsMath: true}},
			{{Text: "FA", IsMath: true}},
			{{Text: "FG", IsMath: true}},
		},
		Answer: 1,
	}

	ans := field.correctAnswer()
	if ans.(client.RadioAnswer).Index != 2 { // 0 based
		t.Fatal(ans)
	}

	props := field.toClient().(client.DropDownFieldBlock).Proposals
	if textLineToString(props[2]) != "GB" {
		t.Fatal(props)
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

	field = FigureAffineLineFieldInstance{
		Figure: repere.Figure{
			Bounds: repere.RepereBounds{
				Width: 20, Height: 20,
				Origin: repere.Coord{4, 4},
			},
		},
		AnswerA: math.Inf(1),
		AnswerB: -2,
	}
	if ans := field.correctAnswer(); !field.evaluateAnswer(ans) {
		t.Fatalf("inconsistent answer %v", ans)
	}
}

func Test_shufflingMap(t *testing.T) {
	tests := []struct {
		n int
	}{
		{5}, {6}, {7},
	}
	for _, tt := range tests {
		shuffler := utils.NewDeterministicRand([]byte{'a', 'b', 'c'})
		got := shufflingMap(shuffler, tt.n)
		indices := make([]int, tt.n)
		for i := range indices {
			indices[i] = i
		}
		shuffler = utils.NewDeterministicRand([]byte{'a', 'b', 'c'})
		shuffler.Shuffle(len(indices), reflect.Swapper(indices))

		new4 := got[4]
		if indices[new4] != 4 {
			t.Errorf("shufflingMap() = %v", got)
		}
	}
}

func TestInstantiate01(t *testing.T) {
	by, err := os.ReadFile("examples/bug01.json")
	if err != nil {
		t.Fatal(err)
	}
	var fullQuestion struct {
		Page QuestionPage `json:"page"`
	}
	if err = json.Unmarshal(by, &fullQuestion); err != nil {
		t.Fatal(err)
	}

	qu, err := fullQuestion.Page.instantiate()
	if err != nil {
		t.Fatal(err)
	}

	qu.ToClient() // test there is no crash

	if err = fullQuestion.Page.Validate(); err == nil {
		t.Fatal("expected error because of non integer values")
	}
}
