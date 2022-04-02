package exercice

import (
	"github.com/benoitkugler/maths-online/maths/expression"
)

const ExhaustiveTextKind = "exhaustiveTextKind"

var (
	_ block = TextBlock{}
	_ block = FormulaBlock{}
	_ block = VariationTableBlock{}
	_ block = SignTableBlock{}
)

// randomParameters is a serialized form of expression.RandomParameters
type randomParameters []randomParameter

type randomParameter struct {
	Expression string `json:"expression"` // as typed by the user, but validated
	Variable   rune   `json:"variable"`
}

// toMap assumes `rp` only contains valid expressions
func (rp randomParameters) toMap() expression.RandomParameters {
	out := make(expression.RandomParameters, len(rp))
	for _, item := range rp {
		out[expression.Variable(item.Variable)], _, _ = expression.Parse(item.Expression)
	}
	return out
}

// Exercice is a sequence of questions
type ExerciceQuestions struct {
	Exercice
	Questions Questions
}

// instantiate returns a deep copy of `eq`, where all random parameters
// have been resolved
// It assumes that the expressions and random parameters definitions are valid.
func (eq ExerciceQuestions) instantiate() ExerciceInstance {
	rp := eq.RandomParameters.toMap()
	// generate random params
	params, _ := rp.Instantiate()

	out := ExerciceInstance{
		Id:          eq.Id,
		Title:       eq.Title,
		Description: eq.Description,
	}
	out.Questions = make([]QuestionInstance, len(eq.Questions))

	for i, qu := range eq.Questions {
		out.Questions[i] = qu.instantiate(params)
	}

	return out
}

func (qu Question) instantiate(params expression.Variables) QuestionInstance {
	enonce := make(EnonceInstance, len(qu.Enonce))
	var currentID int
	for j, bl := range qu.Enonce {
		enonce[j] = bl.instantiate(params, currentID)
		if _, isField := enonce[j].(fieldInstance); isField {
			currentID++
		}
	}
	return QuestionInstance{Title: qu.Title, Enonce: enonce}
}

type TextPart struct {
	Content string
	Kind    TextKind
}

func (tp TextPart) instantiate(params expression.Variables) TextOrMaths {
	switch tp.Kind {
	case Text:
		return TextOrMaths{StringOrExpression: StringOrExpression{String: tp.Content}}
	case StaticMath:
		return TextOrMaths{StringOrExpression: StringOrExpression{String: tp.Content}, IsMath: true}
	case Expression:
		expr, _, _ := expression.Parse(tp.Content)
		expr.Substitute(params)
		return TextOrMaths{StringOrExpression: StringOrExpression{Expression: expr}, IsMath: true}
	default:
		panic(ExhaustiveTextKind)
	}
}

type TextParts []TextPart

func (t TextBlock) instantiate(params expression.Variables, _ int) instance {
	parts := make([]TextOrMaths, len(t.Parts))
	for i, p := range t.Parts {
		parts[i] = p.instantiate(params)
	}
	return TextInstance{
		IsHint: t.IsHint,
		Parts:  parts,
	}
}

// FormulaContent is a list of chunks, either
//	- static math symbols, such as f(x) =
//	- valid expression, such as a*x - b, which will be instantiated
// when rendering the question
//
// For instance, the formula "f(x) = a*(x + 2)"
// is represented by two FormulaPart elements:
// 	{ f(x) = } and { a*(x + 2) }
type FormulaContent []FormulaPart

// FormulaPart forms a logic chunk of a formula.
type FormulaPart struct {
	Content      string
	IsExpression bool // when true, Content is interpreted as an expression.Expression
}

// assume the expression is valid
func (fp FormulaPart) instantiate(params expression.Variables) StringOrExpression {
	if !fp.IsExpression { // nothing to do
		return StringOrExpression{String: fp.Content}
	}

	expr, _, _ := expression.Parse(fp.Content)
	expr.Substitute(params)
	return StringOrExpression{Expression: expr}
}

func (f FormulaBlock) instantiate(params expression.Variables, _ int) instance {
	out := FormulaDisplayInstance{}
	out.Parts = make([]StringOrExpression, len(f.Parts))
	for i, c := range f.Parts {
		out.Parts[i] = c.instantiate(params)
	}
	return FormulaDisplayInstance{}
}

func (vt VariationTableBlock) instantiate(params expression.Variables, _ int) instance {
	out := VariationTableInstance{
		Xs:  make([]float64, len(vt.Xs)),
		Fxs: make([]float64, len(vt.Fxs)),
	}
	for i, c := range vt.Xs {
		out.Xs[i] = mustEvaluate(c, params)
	}
	for i, c := range vt.Fxs {
		out.Fxs[i] = mustEvaluate(c, params)
	}
	return out
}

func (st SignTableBlock) instantiate(params expression.Variables, _ int) instance {
	out := SignTableInstance{
		Xs: make([]string, len(st.Xs)),
	}
	for i, c := range st.Xs {
		out.Xs[i] = c.instantiate(params).asLaTeX()
	}
	out.FxSymbols = append([]SignSymbol(nil), st.FxSymbols...)
	out.Signs = append([]bool(nil), st.Signs...)
	return out
}

func (n NumberField) instantiate(params expression.Variables, ID int) instance {
	expr, _, _ := expression.Parse(n.Expression)
	answer, _ := expr.Evaluate(params)
	return NumberFieldInstance{ID: ID, Answer: answer}
}

// TODO
func (f FormulaField) instantiate(params expression.Variables, ID int) instance {
	return ExpressionFieldInstance{}
}

// TODO
func (l ListField) instantiate(params expression.Variables, ID int) instance {
	return RadioFieldInstance{}
}
