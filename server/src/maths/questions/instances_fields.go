package questions

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/benoitkugler/maths-online/server/src/maths/expression"
	"github.com/benoitkugler/maths-online/server/src/maths/functiongrapher"
	"github.com/benoitkugler/maths-online/server/src/maths/questions/client"
	"github.com/benoitkugler/maths-online/server/src/maths/repere"
	"github.com/benoitkugler/maths-online/server/src/utils"
)

// InvalidFieldAnswer is returned for syntactically incorrect answers
type InvalidFieldAnswer struct {
	ID     int
	Reason string
}

func (ifa InvalidFieldAnswer) Error() string {
	return fmt.Sprintf("field %d: %s", ifa.ID, ifa.Reason)
}

// fieldInstance is an answer field, identified with an integer ID
type fieldInstance interface {
	instance

	fieldID() int

	// evaluateAnswer evaluate the given answer against the reference
	// validateAnswerSyntax is assumed to have already been called on `answer`
	// so that is has a valid format.
	evaluateAnswer(answer client.Answer) (isCorrect bool)

	// correctAnswer returns the expected answer for this field
	// it may not always be unique, in such case the returned value
	// is one of the possible solutions
	correctAnswer() client.Answer

	// validateAnswerSyntax is called during editing for complex fields,
	// to catch syntax mistake before validating the answer
	// an error may also be returned against malicious query
	// if non nil, the error is of type InvalidFieldAnswer
	validateAnswerSyntax(answer client.Answer) error
}

var (
	_ fieldInstance = NumberFieldInstance{}
	_ fieldInstance = ExpressionFieldInstance{}
	_ fieldInstance = RadioFieldInstance{}
	_ fieldInstance = DropDownFieldInstance{}
	_ fieldInstance = OrderedListFieldInstance{}
	_ fieldInstance = GeometricConstructionFieldInstance{}
	_ fieldInstance = VariationTableFieldInstance{}
	_ fieldInstance = SignTableFieldInstance{}
	_ fieldInstance = FunctionPointsFieldInstance{}
	_ fieldInstance = TreeFieldInstance{}
	_ fieldInstance = ProofFieldInstance{}
	_ fieldInstance = TableFieldInstance{}
	_ fieldInstance = VectorFieldInstance{}
)

// NumberFieldInstance is an answer field where only
// numbers are allowed.
// Answers are compared as float values, with a fixed
// precision.
type NumberFieldInstance struct {
	ID     int
	Answer float64 // expected answer
}

func (f NumberFieldInstance) fieldID() int { return f.ID }

func (f NumberFieldInstance) toClient() client.Block {
	s := expression.Number(f.Answer).String()
	return client.NumberFieldBlock{
		ID:       f.ID,
		SizeHint: len([]rune(s)),
	}
}

func (f NumberFieldInstance) validateAnswerSyntax(answer client.Answer) error {
	_, ok := answer.(client.NumberAnswer)
	if !ok {
		return InvalidFieldAnswer{
			ID:     f.ID,
			Reason: fmt.Sprintf("expected NumberAnswer, got %T", answer),
		}
	}
	return nil
}

func (f NumberFieldInstance) evaluateAnswer(answer client.Answer) (isCorrect bool) {
	return expression.AreFloatEqual(f.Answer, answer.(client.NumberAnswer).Value)
}

func (f NumberFieldInstance) correctAnswer() client.Answer {
	return client.NumberAnswer{Value: expression.RoundFloat(f.Answer)}
}

// ExpressionFieldInstance is an answer field where a single mathematical expression
// if expected
type ExpressionFieldInstance struct {
	// if not empty, the field is displayed on a new line
	LabelLaTeX string

	Answer          expression.Compound
	ComparisonLevel ComparisonLevel

	// If true an hint for fraction is displayed
	ShowFractionHelp bool

	ID int
}

func (f ExpressionFieldInstance) fieldID() int { return f.ID }

// add some random padding to avoid leaking to much info about
// the correct answer
func (f ExpressionFieldInstance) sizeHint() int { return len([]rune(f.Answer.String())) + rand.Intn(3) }

func (f ExpressionFieldInstance) toClient() client.Block {
	out := client.ExpressionFieldBlock{
		ID:               f.ID,
		Label:            f.LabelLaTeX,
		SizeHint:         f.sizeHint(),
		ShowFractionHelp: f.ShowFractionHelp,
	}
	if f.ComparisonLevel == AsLinearEquation {
		out.Suffix = " = 0"
	}
	return out
}

func (f ExpressionFieldInstance) validateAnswerSyntax(answer client.Answer) error {
	expr, ok := answer.(client.ExpressionAnswer)
	if !ok {
		return InvalidFieldAnswer{
			ID:     f.ID,
			Reason: fmt.Sprintf("expected ExpressionAnswer, got %T", answer),
		}
	}

	_, err := expression.ParseCompound(expr.Expression)
	if err != nil {
		err := err.(expression.ErrInvalidExpr)
		return InvalidFieldAnswer{
			ID:     f.ID,
			Reason: fmt.Sprintf(`L'expression %s est invalide: %s (à "%s")`, err.Input, err.Reason, err.Portion()),
		}
	}
	return nil
}

func (f ExpressionFieldInstance) evaluateAnswer(answer client.Answer) (isCorrect bool) {
	expr, _ := expression.ParseCompound(answer.(client.ExpressionAnswer).Expression)
	if f.ComparisonLevel == AsLinearEquation {
		return expression.AreLinearEquationsEquivalent(f.Answer, expr)
	}
	return expression.AreCompoundsEquivalent(f.Answer, expr, expression.ComparisonLevel(f.ComparisonLevel))
}

func (f ExpressionFieldInstance) correctAnswer() client.Answer {
	return client.ExpressionAnswer{Expression: f.Answer.String()}
}

// RadioFieldInstance is an answer field where one choice
// is to be made against a fixed list
type RadioFieldInstance struct {
	Proposals []client.TextLine
	ID        int
	Answer    int // index into Proposals, starting at 1
}

func (rf RadioFieldInstance) fieldID() int {
	return rf.ID
}

func (rf RadioFieldInstance) shuffler() utils.Shuffler {
	var hash []byte
	for _, a := range rf.Proposals {
		hash = append(hash, []byte(textLineToString(a))...)
	}
	return utils.NewDeterministicShuffler(hash, len(rf.Proposals))
}

// returns the shuffled proposals
func (rf RadioFieldInstance) proposals() []client.TextLine {
	out := make([]client.TextLine, len(rf.Proposals))
	rf.shuffler().Shuffle(func(dst, src int) { out[dst] = rf.Proposals[src] })
	return out
}

func (rf RadioFieldInstance) toClient() client.Block {
	return client.RadioFieldBlock{
		ID:        rf.ID,
		Proposals: rf.proposals(),
	}
}

func (f RadioFieldInstance) validateAnswerSyntax(answer client.Answer) error {
	_, ok := answer.(client.RadioAnswer)
	if !ok {
		return InvalidFieldAnswer{
			ID:     f.ID,
			Reason: fmt.Sprintf("expected RadioAnswer, got %T", answer),
		}
	}
	return nil
}

func (f RadioFieldInstance) evaluateAnswer(answer client.Answer) (isCorrect bool) {
	ma := f.shuffler().OriginalToShuffled()
	expected := ma[f.Answer-1]
	return expected == answer.(client.RadioAnswer).Index
}

func (f RadioFieldInstance) correctAnswer() client.Answer {
	ma := f.shuffler().OriginalToShuffled()
	expected := ma[f.Answer-1]
	return client.RadioAnswer{Index: expected}
}

type DropDownFieldInstance RadioFieldInstance

func (rf DropDownFieldInstance) fieldID() int { return rf.ID }

func (rf DropDownFieldInstance) toClient() client.Block {
	v := RadioFieldInstance(rf).toClient().(client.RadioFieldBlock)
	return client.DropDownFieldBlock(v)
}

func (f DropDownFieldInstance) validateAnswerSyntax(answer client.Answer) error {
	return RadioFieldInstance(f).validateAnswerSyntax(answer)
}

func (f DropDownFieldInstance) evaluateAnswer(answer client.Answer) (isCorrect bool) {
	return RadioFieldInstance(f).evaluateAnswer(answer)
}

func (f DropDownFieldInstance) correctAnswer() client.Answer {
	return RadioFieldInstance(f).correctAnswer()
}

// OrderedListFieldInstance asks the student to reorder part of the
// given symbols
type OrderedListFieldInstance struct {
	Label               string            // optionnal, LaTeX code displayed in front of the anwser field
	Answer              []client.TextLine // LaTeX code
	AdditionalProposals []client.TextLine // added to Answer when displaying the field
	ID                  int
}

func (olf OrderedListFieldInstance) fieldID() int { return olf.ID }

// proposals groups Answer and AdditionalProposals and shuffle the list
// in a random way, which only depends on the field content though
func (olf OrderedListFieldInstance) proposals() (out []client.TextLine) {
	input := append(append(out, olf.Answer...), olf.AdditionalProposals...)
	out = make([]client.TextLine, len(input))
	// shuffle in a deterministic way
	rd := olf.shuffler()
	rd.Shuffle(func(dst, src int) { out[dst] = input[src] })
	return out
}

func (olf OrderedListFieldInstance) toClient() client.Block {
	out := client.OrderedListFieldBlock{
		ID:           olf.ID,
		Label:        olf.Label,
		AnswerLength: len(olf.Answer),
		Proposals:    olf.proposals(),
	}

	return out
}

func (olf OrderedListFieldInstance) validateAnswerSyntax(answer client.Answer) error {
	list, ok := answer.(client.OrderedListAnswer)
	if !ok {
		return InvalidFieldAnswer{
			ID:     olf.ID,
			Reason: fmt.Sprintf("expected OrderedListAnswer, got %T", answer),
		}
	}

	props := olf.proposals()
	for _, v := range list.Indices {
		if v >= len(props) {
			return InvalidFieldAnswer{
				ID:     olf.ID,
				Reason: fmt.Sprintf("invalid indice %d for length %d", v, len(props)),
			}
		}
	}

	return nil
}

func textLineToString(l client.TextLine) string {
	var s strings.Builder
	for _, c := range l {
		s.WriteString(c.Text)
	}
	return s.String()
}

func areLineEquals(l1, l2 client.TextLine) bool {
	// to avoid suprising errors, we compare the values
	// by concatenating the text
	// for instance x and $x$ are thus equal
	return textLineToString(l1) == textLineToString(l2)
}

func (olf OrderedListFieldInstance) evaluateAnswer(answer client.Answer) (isCorrect bool) {
	list := answer.(client.OrderedListAnswer).Indices

	if len(list) != len(olf.Answer) {
		return false
	}

	// reference and student answer have now the same length
	proposals := olf.proposals()
	for i, ref := range olf.Answer {
		got := proposals[list[i]] // check in `validateAnswerSyntax`

		// we compare by value, not indices, since two different indices may have the same
		// value and then not be distinguable by the student,
		// and also, the indices has been shuffled
		if !areLineEquals(got, ref) {
			return false
		}
	}

	return true
}

func (olf OrderedListFieldInstance) shuffler() utils.Shuffler {
	var hash []byte
	for _, a := range olf.Answer {
		hash = append(hash, []byte(textLineToString(a))...)
	}
	return utils.NewDeterministicShuffler(hash, len(olf.Answer)+len(olf.AdditionalProposals))
}

func (olf OrderedListFieldInstance) correctAnswer() client.Answer {
	rd := olf.shuffler()

	answer := rd.OriginalToShuffled()
	answer = answer[0:len(olf.Answer)] // restrict to answer

	return client.OrderedListAnswer{Indices: answer}
}

type GeometricConstructionFieldInstance struct {
	ID         int
	Field      geoFieldInstance
	Background client.FigureOrGraph
}

func (f GeometricConstructionFieldInstance) fieldID() int { return f.ID }

func (f GeometricConstructionFieldInstance) toClient() client.Block {
	return client.GeometricConstructionFieldBlock{
		ID:         f.ID,
		Field:      f.Field.toClient(),
		Background: f.Background,
	}
}

func (f GeometricConstructionFieldInstance) validateAnswerSyntax(answer client.Answer) error {
	switch f.Field.(type) {
	case gfPoint:
		_, ok := answer.(client.PointAnswer)
		if !ok {
			return InvalidFieldAnswer{
				ID:     f.ID,
				Reason: fmt.Sprintf("expected PointAnswer, got %T", answer),
			}
		}
	case gfVector, gfAffineLine:
		_, ok := answer.(client.DoublePointAnswer)
		if !ok {
			return InvalidFieldAnswer{
				ID:     f.ID,
				Reason: fmt.Sprintf("expected DoublePointAnswer, got %T", answer),
			}
		}
	case gfVectorPair:
		_, ok := answer.(client.DoublePointPairAnswer)
		if !ok {
			return InvalidFieldAnswer{
				ID:     f.ID,
				Reason: fmt.Sprintf("expected DoublePointPairAnswer, got %T", answer),
			}
		}
	}
	return nil
}

func (f GeometricConstructionFieldInstance) evaluateAnswer(answer client.Answer) (isCorrect bool) {
	return f.Field.evaluateAnswer(answer)
}

func (f GeometricConstructionFieldInstance) correctAnswer() client.Answer {
	bounds := f.Background.FigBounds()
	return f.Field.correctAnswer(bounds)
}

type geoFieldInstance interface {
	// evaluateAnswer evaluate the given answer against the reference
	// validateAnswerSyntax is assumed to have already been called on `answer`
	// so that is has a valid format.
	evaluateAnswer(answer client.Answer) (isCorrect bool)

	// correctAnswer returns the expected answer for this field
	// it may not always be unique, in such case the returned value
	// is one of the possible solutions
	correctAnswer(repere.RepereBounds) client.Answer

	toClient() client.GeoField
}

type gfPoint repere.IntCoord

func (f gfPoint) toClient() client.GeoField { return client.GFPoint{} }

func (f gfPoint) evaluateAnswer(answer client.Answer) (isCorrect bool) {
	return repere.IntCoord(f) == answer.(client.PointAnswer).Point
}

func (f gfPoint) correctAnswer(repere.RepereBounds) client.Answer {
	return client.PointAnswer{Point: repere.IntCoord(f)}
}

type gfVector struct {
	Answer repere.IntCoord
	// It true, the vector must be anchored at `AnswerOrigin`
	MustHaveOrigin bool
	AnswerOrigin   repere.IntCoord
}

func (f gfVector) toClient() client.GeoField { return client.GFVector{} }

func (f gfVector) evaluateAnswer(answer client.Answer) (isCorrect bool) {
	ans := answer.(client.DoublePointAnswer)
	vector := repere.IntCoord{
		X: ans.To.X - ans.From.X,
		Y: ans.To.Y - ans.From.Y,
	}
	if f.MustHaveOrigin { // compare vector and origin
		return f.Answer == vector && f.AnswerOrigin == ans.From
	}
	// only compare the vectors
	return f.Answer == vector
}

func (f gfVector) correctAnswer(repere.RepereBounds) client.Answer {
	to := repere.IntCoord{
		X: f.AnswerOrigin.X + f.Answer.X,
		Y: f.AnswerOrigin.Y + f.Answer.Y,
	}
	return client.DoublePointAnswer{From: f.AnswerOrigin, To: to}
}

type gfAffineLine struct {
	Label   string // of the expected affine function
	AnswerA float64
	AnswerB int
}

func (f gfAffineLine) toClient() client.GeoField {
	return client.GFVector{AsLine: true, LineLabel: f.Label}
}

func (f gfAffineLine) isAnswerVertical() bool { return math.IsInf(f.AnswerA, 0) }

func (f gfAffineLine) evaluateAnswer(answer client.Answer) (isCorrect bool) {
	ans := answer.(client.DoublePointAnswer)

	if f.isAnswerVertical() {
		return ans.From.X == f.AnswerB && ans.To.X == f.AnswerB
	}

	a := float64(ans.To.Y-ans.From.Y) / float64(ans.To.X-ans.From.X)
	b := int(float64(ans.From.Y) - a*float64(ans.From.X))
	return f.AnswerA == a && f.AnswerB == b
}

func (f gfAffineLine) correctAnswer(bounds repere.RepereBounds) client.Answer {
	origin := bounds.Origin.Round()

	if f.isAnswerVertical() { // vertical line
		return client.DoublePointAnswer{
			From: repere.IntCoord{X: f.AnswerB, Y: 0},
			To:   repere.IntCoord{X: f.AnswerB, Y: 1},
		}
	}

	// try to get an integer point
	x := -origin.X
	for ; x < bounds.Width-origin.X; x++ {
		y := f.AnswerA * float64(x)
		if math.Trunc(y) == y {
			break
		}
	}
	return client.DoublePointAnswer{
		From: repere.IntCoord{X: 0, Y: f.AnswerB},
		To:   repere.IntCoord{X: x, Y: int(f.AnswerA*float64(x)) + f.AnswerB},
	}
}

type gfVectorPair VectorPairCriterion

func (f gfVectorPair) toClient() client.GeoField { return client.GFVectorPair{} }

func (f gfVectorPair) evaluateAnswer(answer client.Answer) (isCorrect bool) {
	ans := answer.(client.DoublePointPairAnswer)
	vector1 := repere.IntCoord{
		X: ans.To1.X - ans.From1.X,
		Y: ans.To1.Y - ans.From1.Y,
	}
	vector2 := repere.IntCoord{
		X: ans.To2.X - ans.From2.X,
		Y: ans.To2.Y - ans.From2.Y,
	}
	switch VectorPairCriterion(f) {
	case VectorEquals:
		return vector1 == vector2
	case VectorColinear: // check if det(v1, v2) = 0
		return vector1.X*vector2.Y-vector1.Y*vector2.X == 0
	case VectorOrthogonal: // check if v1.v2 = 0
		return vector1.X*vector2.X+vector1.Y*vector2.Y == 0
	default:
		panic("exhaustive switch")
	}
}

func (f gfVectorPair) correctAnswer(repere.RepereBounds) client.Answer {
	switch VectorPairCriterion(f) {
	case VectorEquals:
		return client.DoublePointPairAnswer{
			From1: repere.IntCoord{X: 0, Y: 0},
			To1:   repere.IntCoord{X: 3, Y: 3},
			From2: repere.IntCoord{X: 0, Y: 1},
			To2:   repere.IntCoord{X: 3, Y: 4},
		}
	case VectorColinear:
		return client.DoublePointPairAnswer{
			From1: repere.IntCoord{X: 0, Y: 0},
			To1:   repere.IntCoord{X: 3, Y: 3},
			From2: repere.IntCoord{X: 3, Y: 4},
			To2:   repere.IntCoord{X: -1, Y: 0},
		}
	case VectorOrthogonal:
		return client.DoublePointPairAnswer{
			From1: repere.IntCoord{X: 0, Y: 0},
			To1:   repere.IntCoord{X: 4, Y: 0},
			From2: repere.IntCoord{X: 0, Y: -2},
			To2:   repere.IntCoord{X: 0, Y: 2},
		}
	default:
		panic("exhaustive switch")
	}
}

type VariationTableFieldInstance struct {
	Answer VariationTableInstance
	ID     int
}

func (f VariationTableFieldInstance) fieldID() int { return f.ID }

// lengthProposals returns randomized proposals around the correct value `L`
// the returned value is truely random, but contains L
func lengthProposals(L int) []int {
	var tmp []int

	seed := time.Now().Unix()
	rd := rand.New(rand.NewSource(seed))
	if L <= 1 {
		if rd.Intn(2) == 1 {
			tmp = []int{L, L + 1}
		} else {
			tmp = []int{L, L + 1, L + 2}
		}
	} else {
		tmp = []int{L - 1, L, L + 1}
		// add some random noise to prevent the
		// right solution (L) to be in the middle of the proposals
		// note that we need to ensure L - 1 + r >= 1
		r := rd.Intn(2)
		for i := range tmp {
			tmp[i] += r
		}
	}

	suffler := utils.NewDeterministicShuffler([]byte{byte(seed & 0xff)}, len(tmp))
	out := make([]int, len(tmp))
	suffler.Shuffle(func(dst, src int) { out[dst] = tmp[src] })
	return tmp
}

// lengthProposals returns proposals for the number of arrows to fill,
// depending on the answer and randomized
func (vtf VariationTableFieldInstance) lengthProposals() []int {
	L := len(vtf.Answer.Xs) - 1
	return lengthProposals(L)
}

func (f VariationTableFieldInstance) toClient() client.Block {
	return client.VariationTableFieldBlock{
		Label:           f.Answer.Label,
		LengthProposals: f.lengthProposals(),
		ID:              f.ID,
	}
}

func parseVariationTableAnswer(answer client.VariationTableAnswer) (xs, fxs []*expression.Expr, err error) {
	xs = make([]*expression.Expr, len(answer.Xs))
	fxs = make([]*expression.Expr, len(answer.Fxs))
	for i, x := range answer.Xs {
		xs[i], err = expression.Parse(x)
		if err != nil {
			return nil, nil, err
		}
	}
	for i, fx := range answer.Fxs {
		fxs[i], err = expression.Parse(fx)
		if err != nil {
			return nil, nil, err
		}
	}
	return xs, fxs, nil
}

func (f VariationTableFieldInstance) validateAnswerSyntax(answer client.Answer) error {
	ans, ok := answer.(client.VariationTableAnswer)
	if !ok {
		return InvalidFieldAnswer{
			ID:     f.ID,
			Reason: fmt.Sprintf("expected DoublePointPairAnswer, got %T", answer),
		}
	}

	if L := len(ans.Xs); len(ans.Fxs) != L || len(ans.Arrows) != L-1 {
		return InvalidFieldAnswer{
			ID:     f.ID,
			Reason: fmt.Sprintf("invalid lengths Xs: %d Fxs: %d Arrows: %d", len(ans.Xs), len(ans.Fxs), len(ans.Arrows)),
		}
	}

	_, _, err := parseVariationTableAnswer(ans)
	return err
}

func areNumbersEqual(s1, s2 []float64) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, v := range s1 {
		if !expression.AreFloatEqual(s2[i], v) {
			return false
		}
	}
	return true
}

func areExpressionsEquals(got, exp []*expression.Expr) bool {
	if len(got) != len(exp) {
		return false
	}
	for i, v := range got {
		if !expression.AreExpressionsEquivalent(v, exp[i], expression.SimpleSubstitutions) {
			return false
		}
	}
	return true
}

func areEvExpressionsEquals(got []*expression.Expr, exp []evaluatedExpression) bool {
	tmp := make([]*expression.Expr, len(exp))
	for i, e := range exp {
		tmp[i] = e.Expr
	}
	return areExpressionsEquals(got, tmp)
}

func (f VariationTableFieldInstance) evaluateAnswer(answer client.Answer) (isCorrect bool) {
	ans := answer.(client.VariationTableAnswer)
	xs, fxs, _ := parseVariationTableAnswer(ans)
	if !(areEvExpressionsEquals(xs, f.Answer.Xs) && areEvExpressionsEquals(fxs, f.Answer.Fxs)) {
		return false
	}

	for i, arrow := range ans.Arrows {
		arrowExp := !f.Answer.inferNumberAlignment(i)
		if arrowExp != arrow {
			return false
		}
	}

	return true
}

func (f VariationTableFieldInstance) correctAnswer() client.Answer {
	out := client.VariationTableAnswer{
		Xs:     make([]string, len(f.Answer.Xs)),
		Fxs:    make([]string, len(f.Answer.Fxs)),
		Arrows: make([]bool, len(f.Answer.Xs)-1),
	}
	for i, x := range f.Answer.Xs {
		out.Xs[i] = x.Expr.String()
	}
	for i, x := range f.Answer.Fxs {
		out.Fxs[i] = x.Expr.String()
	}
	for i := range out.Arrows {
		out.Arrows[i] = !f.Answer.inferNumberAlignment(i)
	}
	return out
}

type SignTableFieldInstance struct {
	Answer SignTableInstance
	ID     int
}

func (f SignTableFieldInstance) fieldID() int { return f.ID }

// lengthProposals returns proposals for the number of signs to fill
func (vtf SignTableFieldInstance) lengthProposals() []int {
	L := len(vtf.Answer.Xs) - 1
	return lengthProposals(L)
}

func (f SignTableFieldInstance) toClient() client.Block {
	out := client.SignTableFieldBlock{
		Labels:          make([]string, len(f.Answer.Functions)),
		LengthProposals: f.lengthProposals(),
		ID:              f.ID,
	}
	for i, fo := range f.Answer.Functions {
		out.Labels[i] = fo.Label
	}
	return out
}

func parseSignTableAnswer(answer client.SignTableAnswer) (xs []*expression.Expr, err error) {
	xs = make([]*expression.Expr, len(answer.Xs))
	for i, x := range answer.Xs {
		xs[i], err = expression.Parse(x)
		if err != nil {
			return nil, err
		}
	}
	return xs, nil
}

func (f SignTableFieldInstance) validateAnswerSyntax(answer client.Answer) error {
	ans, ok := answer.(client.SignTableAnswer)
	if !ok {
		return InvalidFieldAnswer{
			ID:     f.ID,
			Reason: fmt.Sprintf("expected DoublePointPairAnswer, got %T", answer),
		}
	}

	if len(ans.Functions) != len(f.Answer.Functions) {
		return InvalidFieldAnswer{
			ID:     f.ID,
			Reason: fmt.Sprintf("invalid lengths for Functions: %d != %d", len(ans.Functions), len(f.Answer.Functions)),
		}
	}

	L := len(ans.Xs)
	for _, function := range ans.Functions {
		if len(function.FxSymbols) != L || len(function.Signs) != L-1 {
			return InvalidFieldAnswer{
				ID:     f.ID,
				Reason: fmt.Sprintf("invalid lengths : Xs = %d Fxs = %d Arrows = %d", len(ans.Xs), len(function.FxSymbols), len(function.Signs)),
			}
		}
	}

	_, err := parseSignTableAnswer(ans)
	return err
}

func (f SignTableFieldInstance) evaluateAnswer(answer client.Answer) (isCorrect bool) {
	ans := answer.(client.SignTableAnswer)
	xs, _ := parseSignTableAnswer(ans)
	if !areExpressionsEquals(xs, f.Answer.Xs) {
		return false
	}

	// here we know the lengths are corrects (areExpressionsEquals validating Xs length)
	for i, exp := range f.Answer.Functions {
		got := ans.Functions[i]
		for j, symbol := range got.FxSymbols {
			if symbol != exp.FxSymbols[j] {
				return false
			}
		}
		for j, sign := range got.Signs {
			if sign != exp.Signs[j] {
				return false
			}
		}

	}

	return true
}

func (f SignTableFieldInstance) correctAnswer() client.Answer {
	out := client.SignTableAnswer{
		Xs:        make([]string, len(f.Answer.Xs)),
		Functions: append([]client.FunctionSign(nil), f.Answer.Functions...),
	}
	for i, x := range f.Answer.Xs {
		out.Xs[i] = x.String()
	}

	return out
}

type FunctionPointsFieldInstance struct {
	IsDiscrete   bool
	Function     expression.FunctionExpr
	Label        string
	XGrid        []int
	ID           int
	offsetHeight int // added to the natural repere height
}

func (f FunctionPointsFieldInstance) fieldID() int { return f.ID }

func (f FunctionPointsFieldInstance) toClient() client.Block {
	bounds, _, dfxs := functiongrapher.PointsFromExpression(f.Function, f.XGrid)
	bounds.Height += f.offsetHeight
	return client.FunctionPointsFieldBlock{
		IsDiscrete: f.IsDiscrete,
		Label:      f.Label,
		Xs:         f.XGrid, ID: f.ID,
		Bounds: bounds,
		Dfxs:   dfxs,
	}
}

func (f FunctionPointsFieldInstance) validateAnswerSyntax(answer client.Answer) error {
	ans, ok := answer.(client.FunctionPointsAnswer)
	if !ok {
		return InvalidFieldAnswer{
			ID:     f.ID,
			Reason: fmt.Sprintf("expected DoublePointPairAnswer, got %T", answer),
		}
	}

	if L := len(ans.Fxs); L != len(f.XGrid) {
		return InvalidFieldAnswer{
			ID:     f.ID,
			Reason: fmt.Sprintf("invalid length %d", L),
		}
	}

	return nil
}

func (f FunctionPointsFieldInstance) evaluateAnswer(answer client.Answer) (isCorrect bool) {
	ans := answer.(client.FunctionPointsAnswer).Fxs
	_, ys, _ := functiongrapher.PointsFromExpression(f.Function, f.XGrid)
	for i := range ys {
		if ans[i] != ys[i] {
			return false
		}
	}
	return true
}

func (f FunctionPointsFieldInstance) correctAnswer() client.Answer {
	_, ys, _ := functiongrapher.PointsFromExpression(f.Function, f.XGrid)
	return client.FunctionPointsAnswer{Fxs: ys}
}

type TreeFieldInstance struct {
	Answer TreeInstance
	ID     int
}

// compute the shape of the given tree
// it assumes the tree is regular, that is the number of branches
// is constant on a given level, but may changes over levels.
func shape(tree TreeNodeInstance) (out client.TreeShape) {
	if len(tree.Children) == 0 {
		return nil
	}
	levelWidth := len(tree.Children)
	return append(client.TreeShape{levelWidth}, shape(tree.Children[0])...)
}

func (f TreeFieldInstance) shapeProposals() []client.TreeShape {
	realShape := shape(f.Answer.AnswerRoot)
	alternative1 := append(client.TreeShape(nil), realShape...)
	alternative1[0] += 1
	alternative2 := append(client.TreeShape(nil), realShape...)
	alternative2[len(alternative2)-1] += 1
	tmp := []client.TreeShape{
		realShape,
		append(realShape, realShape[0]),
		alternative1,
		alternative2,
	}

	var content strings.Builder
	for _, event := range f.Answer.EventsProposals {
		content.WriteString(textLineToString(event))
	}

	rd := utils.NewDeterministicShuffler([]byte(content.String()), len(tmp))
	out := make([]client.TreeShape, len(tmp))
	rd.Shuffle(func(dst, src int) { out[dst] = tmp[src] })
	return tmp
}

func (f TreeFieldInstance) fieldID() int { return f.ID }

func (f TreeFieldInstance) toClient() client.Block {
	return client.TreeFieldBlock{
		ID:              f.ID,
		ShapeProposals:  f.shapeProposals(),
		EventsProposals: f.Answer.EventsProposals,
	}
}

func (f TreeFieldInstance) validateAnswerSyntax(answer client.Answer) error {
	ans, ok := answer.(client.TreeAnswer)
	if !ok {
		return InvalidFieldAnswer{
			ID:     f.ID,
			Reason: fmt.Sprintf("expected TreeAnswer, got %T", answer),
		}
	}

	var isCorrect func(node client.TreeNodeAnswer) error
	isCorrect = func(node client.TreeNodeAnswer) error {
		if len(node.Children) != len(node.Probabilities) {
			return InvalidFieldAnswer{
				ID:     f.ID,
				Reason: "mismatch between Children and Probabilities length",
			}
		}
		for _, expr := range node.Probabilities {
			_, err := expression.Parse(expr)
			if err != nil {
				return InvalidFieldAnswer{
					ID:     f.ID,
					Reason: fmt.Sprintf("invalid Tree probability expression: %s", err),
				}
			}
		}
		// recurse
		for _, child := range node.Children {
			if err := isCorrect(child); err != nil {
				return err
			}
		}
		return nil
	}

	return isCorrect(ans.Root)
}

type treeItem struct {
	proba *expression.Expr
	child TreeNodeInstance
}

func sliceFromNode(node TreeNodeInstance) []treeItem {
	out := make([]treeItem, len(node.Children))
	for i := range out {
		out[i] = treeItem{node.Probabilities[i], node.Children[i]}
	}
	return out
}

// Permute the values at index i to len(a)-1.
//
// perm(a, f, 0) calls f with each permutation of a.
func perm(a []treeItem, f func([]treeItem), i int) {
	if i > len(a) {
		f(a)
		return
	}
	perm(a, f, i+1)
	for j := i + 1; j < len(a); j++ {
		a[i], a[j] = a[j], a[i]
		perm(a, f, i+1)
		a[i], a[j] = a[j], a[i]
	}
}

func areTreeEquivalent(exp, got TreeNodeInstance) bool {
	if exp.Value != got.Value {
		return false
	}

	if len(exp.Children) != len(got.Children) {
		return false
	}

	// compare the probabilities and the associated subtree, up
	// to permutations
	expL := sliceFromNode(exp)
	gotL := sliceFromNode(got)
	onePermCorrect := false
	perm(gotL, func(gotPermutedL []treeItem) {
		allCorrect := true
		// check if expL gotPermutedL match
		for i := range expL {
			expI, gotI := expL[i], gotPermutedL[i]
			if !expression.AreExpressionsEquivalent(expI.proba, gotI.proba, expression.SimpleSubstitutions) {
				allCorrect = false
				break
			}
			// recurse on children (also accepting permutations)
			if !areTreeEquivalent(expI.child, gotI.child) {
				allCorrect = false
				break
			}
		}
		if allCorrect { // we found at least one correct permutation
			onePermCorrect = true
		}
	}, 0)

	return onePermCorrect
}

// parse the client exprs
func treeNodeAnswerToInstance(tna client.TreeNodeAnswer) TreeNodeInstance {
	out := TreeNodeInstance{
		Value:         tna.Value,
		Probabilities: make([]*expression.Expr, len(tna.Probabilities)),
		Children:      make([]TreeNodeInstance, len(tna.Children)),
	}
	for i, expr := range tna.Probabilities {
		out.Probabilities[i], _ = expression.Parse(expr) // checked in validateAnswerSyntax
	}
	for i, child := range tna.Children {
		out.Children[i] = treeNodeAnswerToInstance(child)
	}
	return out
}

func (f TreeFieldInstance) evaluateAnswer(answer client.Answer) (isCorrect bool) {
	ans := answer.(client.TreeAnswer)
	asInstance := treeNodeAnswerToInstance(ans.Root)
	return areTreeEquivalent(f.Answer.AnswerRoot, asInstance)
}

func (f TreeFieldInstance) correctAnswer() client.Answer {
	return client.TreeAnswer{Root: f.Answer.AnswerRoot.toClient()}
}

type TableFieldInstance struct {
	HorizontalHeaders []client.TextOrMath
	VerticalHeaders   []client.TextOrMath
	Answer            client.TableAnswer
	ID                int
}

func (f TableFieldInstance) fieldID() int { return f.ID }

func (f TableFieldInstance) toClient() client.Block {
	return client.TableFieldBlock{
		ID:                f.ID,
		HorizontalHeaders: f.HorizontalHeaders,
		VerticalHeaders:   f.VerticalHeaders,
	}
}

func (f TableFieldInstance) validateAnswerSyntax(answer client.Answer) error {
	_, ok := answer.(client.TableAnswer)
	if !ok {
		return InvalidFieldAnswer{
			ID:     f.ID,
			Reason: fmt.Sprintf("expected TableAnswer, got %T", answer),
		}
	}

	return nil
}

func (f TableFieldInstance) evaluateAnswer(answer client.Answer) (isCorrect bool) {
	ans := answer.(client.TableAnswer)

	if len(ans.Rows) != len(f.Answer.Rows) {
		return false
	}
	for i := range f.Answer.Rows {
		if !areNumbersEqual(f.Answer.Rows[i], ans.Rows[i]) {
			return false
		}
	}
	return true
}

func (f TableFieldInstance) correctAnswer() client.Answer {
	return f.Answer
}

type VectorFieldInstance struct {
	ID             int
	Answer         repere.Coord
	AcceptColinear bool
	DisplayColumn  bool
}

func (v VectorFieldInstance) fieldID() int { return v.ID }

func (v VectorFieldInstance) toClient() client.Block {
	sX := expression.Number(v.Answer.X).String()
	sY := expression.Number(v.Answer.Y).String()
	return client.VectorFieldBlock{
		ID:            v.ID,
		DisplayColumn: v.DisplayColumn,
		SizeHintX:     len([]rune(sX)),
		SizeHintY:     len([]rune(sY)),
	}
}

func (v VectorFieldInstance) validateAnswerSyntax(answer client.Answer) error {
	_, ok := answer.(client.VectorNumberAnswer)
	if !ok {
		return InvalidFieldAnswer{
			ID:     v.ID,
			Reason: fmt.Sprintf("expected VectorNumberAnswer, got %T", answer),
		}
	}
	return nil
}

func (v VectorFieldInstance) evaluateAnswer(answer client.Answer) (isCorrect bool) {
	ans := answer.(client.VectorNumberAnswer)
	if v.AcceptColinear { // check if det(f.Answer, ans) = 0 and ans not 0
		if ans.X == 0 && ans.Y == 0 {
			return false
		}
		return expression.AreFloatEqual(v.Answer.X*ans.Y-v.Answer.Y*ans.X, 0)
	}
	return expression.AreFloatEqual(v.Answer.X, ans.X) && expression.AreFloatEqual(v.Answer.Y, ans.Y)
}

func (v VectorFieldInstance) correctAnswer() client.Answer {
	return client.VectorNumberAnswer{X: expression.RoundFloat(v.Answer.X), Y: expression.RoundFloat(v.Answer.Y)}
}
