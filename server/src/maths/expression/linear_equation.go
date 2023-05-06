package expression

import "fmt"

// linearCoefficients stores the coefficient associated with
// variables in a linear equation (such as 2x - y + 4 = 0)
// the constant term is indexed by the empty (zero value) [Variable]
type linearCoefficients map[Variable]float64

// isEquivalent returns true if the two list of coeffecients
// describe the same equation, that is have the same variables
// and have proportional coefficients
func (lc linearCoefficients) isEquivalent(other linearCoefficients) bool {
	if len(lc) != len(other) {
		return false
	}
	var alpha float64 // such that other = alpha * lc
	for v := range lc {
		coeff1, coeff2 := lc[v], other[v] // coeff1 is not zero
		alpha = coeff2 / coeff1
		break
	}

	for v := range lc {
		coeff1, coeff2 := lc[v], other[v]
		if !AreFloatEqual(coeff1*alpha, coeff2) {
			return false
		}
	}
	return true
}

// duplicate variables are rejected
// term with 0 as coefficients are removed
func (expr *Expr) isLinearEquation() (linearCoefficients, error) {
	// normalize the expression, making sure the receiver is not mutated
	expr = expr.Copy()
	expr.basicSimplification()

	// extract each term
	expr.expandMinus() // required by the following step
	terms := expr.extractOperator(plus)
	out := make(linearCoefficients)
	for _, term := range terms {
		if c, v, ok := term.isLinearTerm(); ok {
			if _, alreadyPresent := out[v]; alreadyPresent {
				return nil, fmt.Errorf("La variable %s apparaît plus d'une fois.", v)
			}
			out[v] = c
		} else if c, ok := term.isConstantTerm(); ok {
			if _, alreadyPresent := out[Variable{}]; alreadyPresent {
				return nil, fmt.Errorf("Le terme constant apparaît plus d'une fois.")
			}
			out[Variable{}] = c
		} else {
			return nil, fmt.Errorf("Le terme %s n'est pas linéaire (ou constant).", term)
		}
	}
	return out, nil
}

// returns true if [expr] is equivalent to [coeff] * [variable]
// expr has already been simplified with basicSimplification
func (expr *Expr) isLinearTerm() (coeff float64, variable Variable, ok bool) {
	if expr == nil {
		return 0, Variable{}, false
	}
	// expr.simplifyNumbers()
	switch atom := expr.atom.(type) {
	case Variable: // x
		return 1, atom, true
	case operator: // -x, -2x, 2x, x / 2
		// recurse
		leftCoeff, leftVar, leftOk := expr.left.isLinearTerm()
		rightCoeff, rightVar, rightOk := expr.right.isLinearTerm()
		switch atom {
		case minus:
			if expr.left == nil && rightOk {
				return -rightCoeff, rightVar, true
			}
		case mult:
			// only one of the term may be linear
			if leftOk && !rightOk {
				if rightValue, ok := expr.right.isConstantTerm(); ok {
					return leftCoeff * rightValue, leftVar, true
				}
			} else if !leftOk && rightOk {
				if leftValue, ok := expr.left.isConstantTerm(); ok {
					return rightCoeff * leftValue, rightVar, true
				}
			}
		case div:
			// only the left term may be linear
			if leftOk {
				if rightValue, ok := expr.right.isConstantTerm(); ok {
					return leftCoeff / rightValue, leftVar, true
				}
			}
		}
		return 0, Variable{}, false
	case Number, constant, function, specialFunction, roundFunc, indice, matrix:
		return 0, Variable{}, false
	default:
		panic(exhaustiveAtomSwitch)
	}
}

// isConstantTerm returns [true] if [expr] is a constant number,
// that is evaluable without context
func (expr *Expr) isConstantTerm() (float64, bool) {
	res, err := expr.evalReal(nil)
	if err != nil {
		return 0, false
	}
	return res.eval(), true
}

// AreLinearEquationsEquivalent returns [true] if both expressions
// defines the same linear equation (up to a non zero factor).
//
// It always returns false if [e1] or [e2] are not [*Expr].
func AreLinearEquationsEquivalent(e1, e2 Compound) bool {
	ee1, ok1 := e1.(*Expr)
	ee2, ok2 := e2.(*Expr)
	if !(ok1 && ok2) {
		return false
	}
	lc1, err1 := ee1.isLinearEquation()
	lc2, err2 := ee2.isLinearEquation()
	return err1 == nil && err2 == nil && lc1.isEquivalent(lc2)
}
