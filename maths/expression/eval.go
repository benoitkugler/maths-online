package expression

type Expression node

type VariablesBinding interface {
	Resolve(v variable) float64
}

func (e *Expression) Evaluate(bindings VariablesBinding) float64 {
	// TODO:
	return 0
}
