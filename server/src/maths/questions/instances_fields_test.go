package questions

import (
	"math"
	"reflect"
	"testing"

	"github.com/benoitkugler/maths-online/server/src/maths/expression"
	"github.com/benoitkugler/maths-online/server/src/maths/questions/client"
	"github.com/benoitkugler/maths-online/server/src/maths/repere"
	"github.com/benoitkugler/maths-online/server/src/utils"
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
		{ExpressionFieldInstance{Answer: expression.MustParse("x+2"), ComparisonLevel: SimpleSubstitutions}, client.ExpressionAnswer{Expression: "x + 2"}, true},
		{ExpressionFieldInstance{Answer: expression.MustParse("x+2"), ComparisonLevel: SimpleSubstitutions}, client.ExpressionAnswer{Expression: "2+x "}, true},
		{ExpressionFieldInstance{Answer: expression.MustParse("x+2"), ComparisonLevel: ExpandedSubstitutions}, client.ExpressionAnswer{Expression: "2+ 1*x "}, true},
		{ExpressionFieldInstance{Answer: expression.MustParse("4x+2y+2"), ComparisonLevel: AsLinearEquation}, client.ExpressionAnswer{Expression: "2x + 1 + y"}, true},
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

	props := field.toClient().(client.DropDownFieldBlock).Proposals
	if textLineToString(props[0]) != "GB" {
		t.Fatal(props)
	}
	ans := field.correctAnswer()
	if ans.(client.RadioAnswer).Index != 0 { // 0 based
		t.Fatal(ans)
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
		AnswerA: 1. / 3,
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
		shuffler := utils.NewDeterministicShuffler([]byte{'a', 'b', 'c'}, tt.n)
		got := shuffler.OriginalToShuffled()

		indices := make([]int, tt.n)
		for i := range indices {
			indices[i] = i
		}
		out := make([]int, len(indices))

		shuffler = utils.NewDeterministicShuffler([]byte{'a', 'b', 'c'}, tt.n)
		shuffler.Shuffle(func(dst, src int) { out[dst] = indices[src] })

		new4 := got[4]
		if out[new4] != 4 {
			t.Errorf("shufflingMap() = %v", got)
		}
	}
}

func TestInstantiate01(t *testing.T) {
	bug01 := QuestionPage{
		// Construire la courbe représentative d'une fonction
		Enonce: []Block{
			TextBlock{
				Parts:   "Construire $C_&f&$ la courbe représentative de la fonction : $&f&(x)=&ax^2+bx+c&$",
				Bold:    false,
				Italic:  false,
				Smaller: false,
			},
			FunctionPointsFieldBlock{
				Function: "ax",
				Label:    "C_f",
				Variable: expression.NewVar('x'),
				XGrid: []string{
					"-2",
					"-1",
					"0",
					"1",
					"2",
				},
			},
		},
		Parameters: Parameters{
			Variables: []RandomParameter{
				{
					Expression: "randChoice(f;g;h;k)",
					Variable:   expression.NewVar('f'),
				},
				{
					Expression: "randChoice(-0,5;-0,25;0,25;0,5;1)",
					Variable:   expression.NewVar('a'),
				},
				{
					Expression: "randint(1;3)",
					Variable:   expression.NewVar('b'),
				},
				{
					Expression: "randint(1;2)",
					Variable:   expression.NewVar('c'),
				},
			},
		},
	}

	qu, err := bug01.instantiate()
	if err != nil {
		t.Fatal(err)
	}

	qu.ToClient() // test there is no crash

	if err = bug01.Validate(); err == nil {
		t.Fatal("expected error because of non integer values")
	}
}
