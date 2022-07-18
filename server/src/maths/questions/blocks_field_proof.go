package questions

import (
	"strings"

	"github.com/benoitkugler/maths-online/maths/expression"
	"github.com/benoitkugler/maths-online/maths/questions/client"
)

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

type ProofStatement struct {
	Content Interpolated
}

type ProofEquality struct {
	Terms string // always LaTeX, equations separated by equals
}

type ProofNode struct {
	Left, Right ProofAssertion `structgen-data:"ignore"`
	Op          client.Binary
}

type ProofSequence struct {
	Parts []ProofAssertion `structgen-data:"ignore"`
}

type ProofAssertion interface {
	instantiate(params expression.Vars) (proofAssertionIns, error)
	validate() error
}

func (v ProofStatement) instantiate(params expression.Vars) (proofAssertionIns, error) {
	out, err := v.Content.instantiate(params)
	return proofStatementIns(out), err
}

func (v ProofStatement) validate() error {
	_, err := v.Content.parse()
	return err
}

func (v ProofEquality) instantiate(params expression.Vars) (proofAssertionIns, error) {
	exprs := strings.Split(v.Terms, "=")
	out := make(proofEqualityIns, len(exprs))
	for i, expr := range exprs {
		s, err := instantiateLaTeXExpr(expr, params)
		if err != nil {
			return nil, err
		}
		out[i] = client.TextLine{{Text: s, IsMath: true}}
	}
	return out, nil
}

func (v ProofEquality) validate() error {
	exprs := strings.Split(v.Terms, "=")
	for _, e := range exprs {
		_, err := expression.Parse(e)
		if err != nil {
			return err
		}
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
