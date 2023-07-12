package questions

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/benoitkugler/maths-online/server/src/maths/expression"
	"github.com/benoitkugler/maths-online/server/src/maths/questions/client"
	"github.com/benoitkugler/maths-online/server/src/maths/repere"
	"github.com/benoitkugler/maths-online/server/src/utils"
	tu "github.com/benoitkugler/maths-online/server/src/utils/testutils"
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
		{SignTableFieldInstance{
			Answer: SignTableInstance{Xs: mustParseMany([]string{"1"}), Functions: []client.FunctionSign{{}}},
		}, client.SignTableAnswer{Xs: []string{"1"}, Functions: []client.FunctionSign{{FxSymbols: []client.SignSymbol{client.Zero}}}}, false},
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

func TestExpressionCompound(t *testing.T) {
	fi := ExpressionFieldBlock{
		Expression:      "{a; 2}",
		ComparisonLevel: Strict,
	}
	inst, err := fi.instantiate(nil, 0)
	tu.AssertNoErr(t, err)
	tu.Assert(t, !inst.(fieldInstance).evaluateAnswer(client.ExpressionAnswer{Expression: "a"}))
	tu.Assert(t, inst.(fieldInstance).evaluateAnswer(client.ExpressionAnswer{Expression: "{2; a}"}))

	fi = ExpressionFieldBlock{
		Expression:       "{a; 2}",
		ComparisonLevel:  Strict,
		ShowFractionHelp: true,
	}
	_, err = fi.setupValidator(nil)
	tu.Assert(t, err != nil)

	fi = ExpressionFieldBlock{
		Expression:      "{a; 2}",
		ComparisonLevel: AsLinearEquation,
	}
	_, err = fi.setupValidator(nil)
	tu.Assert(t, err != nil)

	fi = ExpressionFieldBlock{
		Expression:      "2x + 3y",
		ComparisonLevel: AsLinearEquation,
	}
	_, err = fi.setupValidator(nil)
	tu.AssertNoErr(t, err)
	inst, err = fi.instantiate(nil, 0)
	tu.AssertNoErr(t, err)
	tu.Assert(t, !inst.(fieldInstance).evaluateAnswer(client.ExpressionAnswer{Expression: "2x + y"}))
	tu.Assert(t, inst.(fieldInstance).evaluateAnswer(client.ExpressionAnswer{Expression: "4x + 6y"}))
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
			Rp{
				Expression: "randChoice(f;g;h;k)",
				Variable:   expression.NewVar('f'),
			},
			Rp{
				Expression: "randChoice(-0,5;-0,25;0,25;0,5;1)",
				Variable:   expression.NewVar('a'),
			},
			Rp{
				Expression: "randint(1;3)",
				Variable:   expression.NewVar('b'),
			},
			Rp{
				Expression: "randint(1;2)",
				Variable:   expression.NewVar('c'),
			},
		},
	}

	qu, _, err := bug01.InstantiateErr()
	tu.AssertNoErr(t, err)

	qu.ToClient() // test there is no crash

	if err = bug01.Validate(); err == nil {
		t.Fatal("expected error because of non integer values")
	}
}

func TestSyntaxHint(t *testing.T) {
	b := ExpressionFieldBlock{Expression: "{2x / 4; -inf; x^4 +9}"}
	hint, err := b.SyntaxHint(nil)
	tu.AssertNoErr(t, err)

	fmt.Println(hint.Parts)
}

func TestSignTableField(t *testing.T) {
	field := SignTableFieldInstance{
		Answer: SignTableInstance{
			Xs: mustParseMany([]string{"-2", "1", "inf"}),
			Functions: []client.FunctionSign{
				{
					FxSymbols: []client.SignSymbol{client.ForbiddenValue, client.Zero, client.Nothing},
					Signs:     []bool{true, false},
				},
				{
					FxSymbols: []client.SignSymbol{client.ForbiddenValue, client.Zero, client.Nothing},
					Signs:     []bool{false, true},
				},
			},
		},
	}

	_ = field.toClient()
	_ = field.correctAnswer()

	ans := client.SignTableAnswer{
		Xs:        []string{"-2", "1", "inf"},
		Functions: field.Answer.Functions,
	}

	tu.Assert(t, field.evaluateAnswer(ans))
}

func TestTree(t *testing.T) {
	level := client.TreeNodeAnswer{
		Value:         0,
		Probabilities: []float64{0.2, 0.8},
		Children: []client.TreeNodeAnswer{
			{Value: 1},
			{Value: 2},
		},
	}
	levelRev := client.TreeNodeAnswer{
		Value:         0,
		Probabilities: []float64{0.8, 0.2},
		Children: []client.TreeNodeAnswer{
			{Value: 2},
			{Value: 1},
		},
	}

	level3 := client.TreeNodeAnswer{
		Value:         0,
		Probabilities: []float64{0.2, 0.7, 0.1},
		Children: []client.TreeNodeAnswer{
			{Value: 1},
			{Value: 2},
			{Value: 3},
		},
	}
	level3Rev := client.TreeNodeAnswer{
		Value:         0,
		Probabilities: []float64{0.7, 0.2, 0.1},
		Children: []client.TreeNodeAnswer{
			{Value: 2},
			{Value: 1},
			{Value: 3},
		},
	}

	levelOther := client.TreeNodeAnswer{
		Value:         0,
		Probabilities: []float64{0.8, 0.2},
		Children: []client.TreeNodeAnswer{
			{Value: 2},
			{Value: 4},
		},
	}
	levelOtherOther := client.TreeNodeAnswer{
		Value:         0,
		Probabilities: []float64{0.7, 0.2, 0.1},
		Children: []client.TreeNodeAnswer{
			{Value: 2},
			{Value: 4},
			{Value: 4},
		},
	}

	tu.Assert(t, !areTreeEquivalent(level, levelOther))
	tu.Assert(t, !areTreeEquivalent(levelOther, levelOtherOther))
	tu.Assert(t, areTreeEquivalent(level, levelRev))
	tu.Assert(t, areTreeEquivalent(level3, level3Rev))

	// test for nested invertion
	level2 := client.TreeNodeAnswer{
		Value:         0,
		Probabilities: []float64{0.2, 0.8},
		Children: []client.TreeNodeAnswer{
			levelOther,
			level,
		},
	}
	level2bis := client.TreeNodeAnswer{
		Value:         0,
		Probabilities: []float64{0.2, 0.8},
		Children: []client.TreeNodeAnswer{
			levelOther,
			levelRev,
		},
	}
	tu.Assert(t, areTreeEquivalent(level2, level2bis))
}
