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

func (TextBlock) isBlock()         {}
func (FormulaBlock) isBlock()      {}
func (NumberFieldBlock) isBlock()  {}
func (ListFieldBlock) isBlock()    {}
func (FormulaFieldBlock) isBlock() {}

type TextBlock struct {
	Text string
}

type FormulaBlock struct {
	Content  string // as latex
	IsInline bool
}

// NumberFieldBlock is an answer field where only
// numbers are allowed
// answers are compared as float values
type NumberFieldBlock struct {
	ID int
}

// TODO:
type ListFieldBlock struct {
	Choices []string
	ID      int
}

// TODO:
type FormulaFieldBlock struct {
	Expression string // a valid expression, in the format used by expression.Expression
	ID         int
}

// Answer is a sum type for the possible answers
// of question fields
type Answer interface {
	isAnswer()
}

func (NumberAnswer) isAnswer() {}
func (RadioAnswer) isAnswer()  {}

// NumberAnswer is compared with exact float equality
type NumberAnswer struct {
	Value float64
}

// RadioAnswer is compared against a reference index
type RadioAnswer struct {
	Index int
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
