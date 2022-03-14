package client

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

// Answer is an sum type for the possible answers
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

// Answers map the field ids to their answer
type Answers map[int]Answer
