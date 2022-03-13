package exercice

import "encoding/json"

//go:generate ../../../../../structgen/structgen -source=client_types.go -mode=dart:../../../../eleve/lib/exercices/types.gen.dart  -mode=itfs-json:gen_itfs_client.go

type ClientQuestion struct {
	Title   string
	Content ClientContent
}

type ClientContent []clientBlock

func (evs ClientContent) MarshalJSON() ([]byte, error) {
	tmp := make([]clientBlockWrapper, len(evs))
	for i, v := range evs {
		tmp[i] = clientBlockWrapper{Data: v}
	}
	return json.Marshal(tmp)
}

func (evs *ClientContent) UnmarshalJSON(data []byte) error {
	var tmp []clientBlockWrapper
	err := json.Unmarshal(data, &tmp)
	*evs = make(ClientContent, len(tmp))
	for i, v := range tmp {
		(*evs)[i] = v.Data
	}
	return err
}

type clientBlock interface {
	isClientBlock()
}

func (clientTextBlock) isClientBlock()         {}
func (clientFormulaBlock) isClientBlock()      {}
func (clientNumberFieldBlock) isClientBlock()  {}
func (clientListFieldBlock) isClientBlock()    {}
func (clientFormulaFieldBlock) isClientBlock() {}

type clientTextBlock struct {
	Text string
}

type clientFormulaBlock struct {
	Content  string // as latex
	IsInline bool
}

// clientNumberFieldBlock is an answer field where only
// numbers are allowed
// answers are compared as float values
type clientNumberFieldBlock struct{}

// TODO:
type clientListFieldBlock struct {
	// Id      string
	Choices []string
}

// TODO:
type clientFormulaFieldBlock struct {
	// Id string
	Expression string // a valid expression, in the format used by expression.Expression
}
