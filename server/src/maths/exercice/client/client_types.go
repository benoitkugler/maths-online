package client

// TODO: autogenerate JSON wrappers for type using interfaces

import "encoding/json"

//go:generate ../../../../../../structgen/structgen -source=client_types.go -mode=dart:../../../../../eleve/lib/exercices/types.gen.dart  -mode=itfs-json:gen_itfs_client.go

type Question struct {
	Title  string
	Enonce Enonce
}

type Enonce []Block

func (evs Enonce) MarshalJSON() ([]byte, error) {
	tmp := make([]BlockWrapper, len(evs))
	for i, v := range evs {
		tmp[i] = BlockWrapper{Data: v}
	}
	return json.Marshal(tmp)
}

func (evs *Enonce) UnmarshalJSON(data []byte) error {
	var tmp []BlockWrapper
	err := json.Unmarshal(data, &tmp)
	*evs = make(Enonce, len(tmp))
	for i, v := range tmp {
		(*evs)[i] = v.Data
	}
	return err
}

type Block interface {
	isBlock()
}

func (TextBlock) isBlock()             {}
func (FormulaBlock) isBlock()          {}
func (VariationTableBlock) isBlock()   {}
func (SignTableBlock) isBlock()        {}
func (NumberFieldBlock) isBlock()      {}
func (ListFieldBlock) isBlock()        {}
func (ExpressionFieldBlock) isBlock()  {}
func (RadioFieldBlock) isBlock()       {}
func (OrderedListFieldBlock) isBlock() {}

// TextOrMath is a part of a text line, rendered
// either as plain text or using LaTeX in text mode.
type TextOrMath struct {
	Text   string
	IsMath bool
}

type TextBlock struct {
	Parts []TextOrMath
}

// FormulaBlock is whole line, rendered as LaTeX in display mode
type FormulaBlock struct {
	Formula string // as latex
}

// VariationColumn is a column in a variation table,
// either displaying (x, f(x)) values, or an arrow
// between two local extrema.
type VariationColumn struct {
	X, Y    string // as LaTeX code
	IsArrow bool
	IsUp    bool
}

type VariationTableBlock struct {
	Columns []VariationColumn
}

// SignColumn is a column in a sign table,
// either displaying (x, f(x)) values, or an arrow
// between two local extrema.
type SignColumn struct {
	X                 string // as LaTeX code
	IsYForbiddenValue bool   // if true, Y is ignored and a double bar is displayed instead
	IsSign            bool
	IsPositive        bool // for signs, displays a +, for numbers displays a 0 (else nothing)
}

type SignTableBlock struct {
	Columns []SignColumn
}

// NumberFieldBlock is an answer field where only
// numbers are allowed
// answers are compared as float values
type NumberFieldBlock struct {
	ID int
}
type ExpressionFieldBlock struct {
	Label string // as LaTeX, optional
	ID    int
}

type ListFieldProposal struct {
	Content []TextOrMath
}

type RadioFieldBlock struct {
	Proposals []ListFieldProposal
	ID        int
}

type OrderedListFieldBlock struct {
	// Proposals is a shuffled version of the list,
	// displayed as math text
	Proposals    []string
	AnswerLength int
	ID           int
}

// TODO:
type ListFieldBlock struct {
	Choices []string
	ID      int
}

// Answer is a sum type for the possible answers
// of question fields
type Answer interface {
	isAnswer()
}

func (NumberAnswer) isAnswer()      {}
func (RadioAnswer) isAnswer()       {}
func (ExpressionAnswer) isAnswer()  {}
func (OrderedListAnswer) isAnswer() {}

// NumberAnswer is compared with exact float equality
type NumberAnswer struct {
	Value float64
}

type ExpressionAnswer struct {
	Expression string
}

// RadioAnswer is compared against a reference index
type RadioAnswer struct {
	Index int
}

type OrderedListAnswer struct {
	Indices []int // indices into the question field proposals
}

// QuestionAnswersIn map the field ids to their answer
type QuestionAnswersIn struct {
	Data map[int]Answer
}

func (out *QuestionAnswersIn) UnmarshalJSON(src []byte) error {
	var wr struct {
		Data map[int]AnswerWrapper
	}

	err := json.Unmarshal(src, &wr)
	out.Data = make(map[int]Answer)
	for i, v := range wr.Data {
		out.Data[i] = v.Data
	}

	return err
}

func (out QuestionAnswersIn) MarshalJSON() ([]byte, error) {
	var tmp struct {
		Data map[int]AnswerWrapper
	}
	tmp.Data = make(map[int]AnswerWrapper)
	for k, v := range out.Data {
		tmp.Data[k] = AnswerWrapper{v}
	}
	return json.Marshal(tmp)
}

type QuestionAnswersOut struct {
	Data map[int]bool
}

// QuestionSyntaxCheckIn is emitted by the client
// to perform a preliminary check of the syntax,
// without validating the answer
type QuestionSyntaxCheckIn struct {
	Answer Answer
	ID     int
}

func (out *QuestionSyntaxCheckIn) UnmarshalJSON(src []byte) error {
	var wr struct {
		Answer AnswerWrapper
		ID     int
	}
	err := json.Unmarshal(src, &wr)
	out.Answer = wr.Answer.Data
	out.ID = wr.ID

	return err
}

func (out QuestionSyntaxCheckIn) MarshalJSON() ([]byte, error) {
	wr := struct {
		Answer AnswerWrapper
		ID     int
	}{
		Answer: AnswerWrapper{out.Answer},
		ID:     out.ID,
	}
	return json.Marshal(wr)
}

type QuestionSyntaxCheckOut struct {
	Reason  string
	IsValid bool
}
