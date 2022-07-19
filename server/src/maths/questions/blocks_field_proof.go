package questions

import (
	"errors"
	"strings"

	"github.com/benoitkugler/maths-online/maths/expression"
	"github.com/benoitkugler/maths-online/maths/questions/client"
)

//go:generate ../../../../../structgen/structgen -source=blocks_field_proof.go -mode=itfs-json:gen_itfs_proof.go

type ProofFieldBlock struct {
	Answer ProofSequence
}

// ID is only used by answer fields
func (pr ProofFieldBlock) instantiate(params expression.Vars, ID int) (instance, error) {
	ins, err := pr.Answer.instantiate(params)
	if err != nil {
		return nil, err
	}
	return ProofFieldInstance{
		Answer: ins.(proofSequenceIns),
		ID:     ID,
	}, nil
}

func (pr ProofFieldBlock) setupValidator(expression.RandomParameters) (validator, error) {
	err := pr.Answer.validate()
	return noOpValidator{}, err
}

// Placeholder used when editing the proof
type ProofInvalid struct{}

type ProofStatement struct {
	Content Interpolated
}

type ProofEquality struct {
	Terms string // always interpreted as LaTeX, splitted by equals
}

type ProofNode struct {
	Left, Right ProofAssertion `structgen-data:"ignore"`
	Op          client.Binary
}

type ProofAssertions []ProofAssertion

type ProofSequence struct {
	Parts ProofAssertions `structgen-data:"ignore"`
}

type ProofAssertion interface {
	instantiate(params expression.Vars) (proofAssertionIns, error)
	validate() error
}

func (ProofInvalid) instantiate(params expression.Vars) (proofAssertionIns, error) {
	return nil, errors.New("La preuve ne doit pas contenir d'assertion en construction.")
}

func (ProofInvalid) validate() error {
	return errors.New("La preuve ne doit pas contenir d'assertion en construction.")
}

func (v ProofStatement) instantiate(params expression.Vars) (proofAssertionIns, error) {
	out, err := v.Content.instantiate(params)
	return proofStatementIns(out), err
}

func (v ProofStatement) validate() error {
	_, err := v.Content.parse()
	return err
}

func ensureLaTeX(s string) string {
	s = strings.TrimSpace(s)
	if s != "" && !strings.HasPrefix(s, "$") { // ensure LaTeX
		s = "$" + s + "$"
	}
	return s
}

func (v ProofEquality) terms() (members []Interpolated, avecDef Interpolated) {
	// start by the optional avec separator
	const avecSep = "avec"
	chunks := strings.SplitN(v.Terms, avecSep, 2)

	if len(chunks) >= 2 {
		avecDef = Interpolated(chunks[1])
	}

	terms := strings.Split(chunks[0], "=")
	members = make([]Interpolated, len(terms))
	for i, t := range terms {
		t = ensureLaTeX(t)
		members[i] = Interpolated(t)
	}
	return members, avecDef
}

func (v ProofEquality) instantiate(params expression.Vars) (proofAssertionIns, error) {
	members, avecDef := v.terms()
	out := proofEqualityIns{
		Terms: make([]client.TextLine, len(members)),
	}
	var err error
	for i, term := range members {
		line, err := term.instantiate(params)
		if err != nil {
			return nil, err
		}
		out.Terms[i] = line
	}
	out.Def, err = avecDef.instantiate(params)
	if err != nil {
		return nil, err
	}
	out.WithDef = len(out.Def) != 0

	return out, nil
}

func (v ProofEquality) validate() error {
	members, avecDef := v.terms()
	for _, term := range members {
		if term == "" {
			return errors.New("Les membres d'une égalité ne peuvent pas être vides.")
		}
		_, err := term.parse()
		if err != nil {
			return err
		}
	}
	if _, err := avecDef.parse(); err != nil {
		return err
	}
	return nil
}

func (v ProofNode) instantiate(params expression.Vars) (proofAssertionIns, error) {
	left, err := v.Left.instantiate(params)
	if err != nil {
		return nil, err
	}
	right, err := v.Right.instantiate(params)
	if err != nil {
		return nil, err
	}
	return proofNodeIns{Op: v.Op, Left: left, Right: right}, nil
}

func (v ProofNode) validate() error {
	if err := v.Left.validate(); err != nil {
		return err
	}
	if err := v.Right.validate(); err != nil {
		return err
	}
	return nil
}

func (v ProofSequence) instantiate(params expression.Vars) (proofAssertionIns, error) {
	out := make(proofSequenceIns, len(v.Parts))
	var err error
	for i, p := range v.Parts {
		out[i], err = p.instantiate(params)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func (v ProofSequence) validate() error {
	for _, p := range v.Parts {
		if err := p.validate(); err != nil {
			return err
		}
	}
	return nil
}
