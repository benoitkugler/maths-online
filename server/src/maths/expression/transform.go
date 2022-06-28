package expression

import (
	"fmt"
	"sort"
	"strings"
)

// return a slice of all the operands of the `op` operator at the level of `expr`,
// or a one length slice
func (expr *Expr) extractOperator(op operator) []*Expr {
	if expr == nil {
		return nil
	}

	var out []*Expr
	if expr.atom == op {
		out = append(out, expr.left.extractOperator(op)...)
		out = append(out, expr.right.extractOperator(op)...)
	} else {
		out = append(out, expr)
	}
	return out
}

// returns -1 if n1 < n2, 0 if n1 == n2, 1 if n1 > n2
// in the sense of a lexical ordering
// examples :
//	2 < 3 < e < pi < x < y < func < + < - < mult < div
func compareNodes(n1, n2 *Expr) int {
	a1, a2 := n1.atom, n2.atom
	l1, l2 := a1.lexicographicOrder(), a2.lexicographicOrder()
	if l1 < l2 {
		return -1
	} else if l1 > l2 {
		return +1
	} else {
		// both atoms have the same type
		switch a1 := a1.(type) {
		case operator:
			a2 := a2.(operator)
			if a1 < a2 {
				return -1
			} else if a1 > a2 {
				return 1
			} else {
				// use the children
				if n1.left == nil && n2.left != nil {
					return -1
				} else if n1.left != nil && n2.left == nil {
					return +1
				} else if n1.left == nil && n2.left == nil {
					return compareNodes(n1.right, n2.right)
				}
				if c := compareNodes(n1.left, n2.left); c != 0 {
					return c
				}
				return compareNodes(n1.right, n2.right)
			}
		case function:
			a2 := a2.(function)
			if a1 < a2 {
				return -1
			} else if a1 > a2 {
				return 1
			} else {
				return compareNodes(n1.right, n2.right)
			}
		case roundFn:
			a2 := a2.(roundFn)
			if a1.nbDigits < a2.nbDigits {
				return -1
			} else if a1.nbDigits > a2.nbDigits {
				return 1
			} else {
				return compareNodes(n1.right, n2.right)
			}
		case specialFunctionA:
			a2 := a2.(specialFunctionA)
			if a1.kind < a2.kind {
				return -1
			} else if a1.kind > a2.kind {
				return 1
			} else {
				return strings.Compare(fmt.Sprint(a1.args), fmt.Sprint(a2.args))
			}
		case Variable:
			a2 := a2.(Variable)
			if a1.Name < a2.Name {
				return -1
			} else if a1.Name > a2.Name {
				return 1
			} else {
				return strings.Compare(a1.Indice, a2.Indice)
			}
		case constant:
			a2 := a2.(constant)
			if a1 < a2 {
				return -1
			} else if a1 > a2 {
				return 1
			} else {
				return 0
			}
		case Number:
			a2 := a2.(Number)
			if a1 < a2 {
				return -1
			} else if a1 > a2 {
				return 1
			} else {
				return 0
			}
		default:
			panic(exhaustiveAtomSwitch)
		}
	}
}

// will panic if nodes is empty
func nodesAsTree(nodes []*Expr, op operator) Expr {
	expr := nodes[0]
	for _, n := range nodes[1:] {
		expr = &Expr{atom: op, left: expr, right: n}
	}
	return *expr
}

// use associativity and commutativity to reorder (in place) + and * operations
// in a "cannonical" order
func (expr *Expr) sortPlusAndMultOperands() {
	if expr == nil {
		return
	}

	if op, isOp := expr.atom.(operator); isOp && (op == plus || op == mult) {
		nodes := expr.extractOperator(op)
		// begin by ordering children
		for _, n := range nodes {
			n.sortPlusAndMultOperands()
		}

		// then sort
		sort.SliceStable(nodes, func(i, j int) bool { return compareNodes(nodes[i], nodes[j]) < 0 })

		// finaly insert nodes back with the new order
		*expr = nodesAsTree(nodes, op) // nodes is not empty since expr.atom == op
	} else {
		expr.left.sortPlusAndMultOperands()
		expr.right.sortPlusAndMultOperands()
	}
}

// expandMult distribute * over + (in place)
func (expr *Expr) expandMult() {
	if expr == nil {
		return
	}

	expr.left.expandMult()
	expr.right.expandMult()

	if expr.atom != mult {
		return
	}

	// a multiplication always has left and right children
	if expr.left.atom == plus { // (a+b) * c => a*c + b*c
		a, b := expr.left.left, expr.left.right
		c := expr.right
		expr.atom = plus
		expr.left = &Expr{atom: mult, left: a, right: c}
		expr.right = &Expr{atom: mult, left: b, right: c}
	} else if expr.right.atom == plus { // c * (a+b) => c*a + c*b
		a, b := expr.right.left, expr.right.right
		c := expr.left
		expr.atom = plus
		expr.left = &Expr{atom: mult, left: c, right: a}
		expr.right = &Expr{atom: mult, left: c, right: b}
	} // double distributivity is handled by recursion
}

// replace integral power by product, so that
// expandMult can trigger
func (expr *Expr) expandPow() {
	if expr == nil {
		return
	}

	expr.left.expandPow()
	expr.right.expandPow()

	if expr.atom != pow {
		return
	}

	if power, ok := expr.right.atom.(Number); ok {
		asInt := int(power)
		if float64(asInt) == float64(power) && asInt > 0 { // c^d = c * c * c * c ... * c
			exprNew := expr.left.Copy()
			for k := 1; k < asInt; k++ {
				exprNew = &Expr{atom: mult, left: exprNew, right: expr.left.Copy()}
			}
			*expr = *exprNew
		}
	}
}

// replace a + a + a by 3*a
// should be applied after sorting operands
func (expr *Expr) groupAdditions() {
	if expr == nil {
		return
	}

	if expr.atom != plus { // recurse and return early
		expr.left.groupAdditions()
		expr.right.groupAdditions()
		return
	}

	nodes := expr.extractOperator(plus)

	// recurse
	for _, n := range nodes {
		n.groupAdditions()
	}

	var (
		newNodes []*Expr
		count    = 1
		ref      = nodes[0]
	)
	for _, n := range nodes[1:] {
		if ref.equals(n) {
			count++
		} else { // combine the nodes
			if count > 1 {
				newNodes = append(newNodes, &Expr{
					atom:  mult,
					left:  NewNb(float64(count)),
					right: ref,
				})
			} else {
				newNodes = append(newNodes, ref)
			}
			// reset the trackers
			ref = n
			count = 1
		}
	}
	// add the last chunk
	if count > 1 {
		newNodes = append(newNodes, &Expr{
			atom:  mult,
			left:  NewNb(float64(count)),
			right: ref,
		})
	} else {
		newNodes = append(newNodes, ref)
	}

	combined := nodesAsTree(newNodes, plus)
	*expr = combined
}

// replace a-c by a + (-c)
// so that plus operation may trigger
func (expr *Expr) expandMinus() {
	if expr == nil {
		return
	}

	expr.left.expandMinus()
	expr.right.expandMinus()

	if expr.atom != minus {
		return
	}

	// do not transform (-a)
	if expr.left != nil { // a - c => a + (-c)
		expr.atom = plus
		if number, isNumber := expr.right.atom.(Number); isNumber {
			expr.right = &Expr{atom: Number(-number)}
		} else { // general case
			expr.right = &Expr{atom: minus, right: expr.right}
		}
		return
	}
}

// replace + (- 8) by -8 to have a better formatted output
func (expr *Expr) contractPlusMinus() {
	if expr == nil {
		return
	}

	expr.left.contractPlusMinus()
	expr.right.contractPlusMinus()

	if expr.atom != plus {
		return
	}

	// ... + (-...) => ... - ...
	if isNegative, opposite := expr.right.isNegativeExpr(); isNegative {
		expr.atom = minus
		expr.right = opposite
		return
	}
}

// replace negative numbers -3 by -(3) to simplify further processing
func (expr *Expr) normalizeNegativeNumbers() {
	if expr == nil {
		return
	}

	expr.left.normalizeNegativeNumbers()
	expr.right.normalizeNegativeNumbers()

	if number, isNumber := expr.atom.(Number); isNumber && number < 0 {
		expr.atom = minus
		expr.right = &Expr{atom: -number}
	}
}

// returns true and the right term for expression of the form - (...)
func (expr *Expr) isNegativeExpr() (bool, *Expr) {
	if expr.atom == minus && expr.left == nil {
		return true, expr.right
	}
	return false, nil
}

// replace - (- 8) by +8 to have a better formatted output
func (expr *Expr) contractMinusMinus() {
	if expr == nil {
		return
	}

	expr.left.contractMinusMinus()
	expr.right.contractMinusMinus()

	if expr.atom != minus {
		return
	}

	// ... - (-...) => ... + ...
	if isNegative, opposite := expr.right.isNegativeExpr(); isNegative {
		expr.atom = plus
		expr.right = opposite
		return
	}
}

// remove unnecessary 1 and 0 such as in
// 	1 * x -> x
//	 0x -> 0
// -1 * x -> -x
func (expr *Expr) simplify0And1() {
	if expr == nil {
		return
	}

	expr.left.simplify0And1()
	expr.right.simplify0And1()

	op, ok := expr.atom.(operator)
	if !ok {
		return
	}

	left := expr.left
	if expr.left == nil { // 0 is a valid default value
		left = NewNb(0)
	}
	right := expr.right

	// multiplying or dividing by 1;
	// adding or substracting 0 are no-ops
	switch op {
	case plus:
		if left.atom == Number(0) { // 0 + x = x
			*expr = *expr.right
			return
		} else if right.atom == Number(0) { // x + 0 = x
			*expr = *left
			return
		}
	case minus:
		if right.atom == Number(0) { // x - 0 = x
			*expr = *left
			return
		}
	case mult:
		if left.atom == Number(1) { // 1 * x = x
			*expr = *expr.right
			return
		} else if right.atom == Number(1) { // x * 1 = x
			*expr = *left
			return
		} else if left.atom == Number(0) { // 0 * x = 0
			*expr = Expr{atom: Number(0)}
			return
		} else if right.atom == Number(0) {
			*expr = Expr{atom: Number(0)}
			return
		}
	case div:
		if right.atom == Number(1) { // x / 1 = x
			*expr = *left
			return
		} else if left.atom == Number(0) && right.atom != Number(0) { // 0 / x = 0
			*expr = Expr{atom: Number(0)}
			return
		}
	case pow:
		if right.atom == Number(1) { // x ^ 1 = x
			*expr = *left
			return
		} else if left.atom == Number(1) { // 1 ^ x = 1
			*expr = Expr{atom: Number(1)}
			return
		}
	case mod:
	case rem:
	default:
		panic(exhaustiveOperatorSwitch)
	}
}

// replace (-a) * b by -(a * b) (same with division)
// requiring normalizeNegativeNumbers() to have been called
func (expr *Expr) extractNegativeInMults() {
	if expr == nil {
		return
	}

	expr.left.extractNegativeInMults()
	expr.right.extractNegativeInMults()

	if expr.atom != mult && expr.atom != div {
		return
	}

	// replace (-a) * b by -(a * b)
	changeSign := false
	newLeft := expr.left
	if isNegative, opposite := expr.left.isNegativeExpr(); isNegative {
		newLeft = opposite
		changeSign = !changeSign
	}
	newRight := expr.right
	if isNegative, opposite := expr.right.isNegativeExpr(); isNegative {
		newRight = opposite
		changeSign = !changeSign
	}

	newExpr := &Expr{atom: expr.atom, left: newLeft, right: newRight}
	if changeSign { // wrap with minus
		expr.atom = minus
		expr.left = nil
		expr.right = newExpr
	} else {
		*expr = *newExpr
	}
}

const maxIterations = 10_000 // very very unlikely in pratice

func (expr *Expr) basicSimplification() (nbPasses int) {
	ref := expr.Copy()

	// apply each transformation until no one triggers a change
	for nbPasses = 1; nbPasses < maxIterations; nbPasses++ {
		expr.normalizeNegativeNumbers()

		expr.expandMinus()
		expr.sortPlusAndMultOperands()
		expr.simplifyForPrint()
		expr.extractNegativeInMults()

		if expr.equals(ref) {
			break
		}
		ref = expr.Copy() // update the reference and start a new pass
	}

	return nbPasses
}

func (expr *Expr) fullSimplification() (nbPasses int) {
	ref := expr.Copy()

	// apply each transformation until no one triggers a change
	for nbPasses = 1; nbPasses < maxIterations; nbPasses++ {
		expr.simplifyNumbers()
		expr.normalizeNegativeNumbers()

		expr.expandPow()
		expr.expandMinus()
		expr.expandMult()
		expr.sortPlusAndMultOperands()
		expr.groupAdditions()
		expr.simplify0And1()
		expr.simplifyNumbers()
		if expr.equals(ref) {
			break
		}
		ref = expr.Copy() // update the reference and start a new pass
	}

	// extractNegativeInMults interfers with other transforms, do it later
	for nbPasses = 1; nbPasses < maxIterations; nbPasses++ {
		expr.normalizeNegativeNumbers()
		expr.extractNegativeInMults()
		expr.simplifyNumbers()

		if expr.equals(ref) {
			break
		}
		ref = expr.Copy() // update the reference and start a new pass
	}

	return nbPasses
}

// ComparisonLevel speficies how mathematical expressions should be
// compared.
// Depending on the context, it is preferable to ask for exact matching
// (like when learning distributivity), or to accept broader answers,
// such as for derivatives.
type ComparisonLevel uint8

const (
	// Expressions are compared structurally
	// This is usually too restrictive to be useful
	// For instance, x + 2 != 2 + x according to this level,
	// or 1/4 != 0.25
	Strict ComparisonLevel = iota
	// Expressions are compared after performing some basic transformations,
	// such as reording operands.
	// Also, if both expressions may be evaluated, equal values are considered equal.
	SimpleSubstitutions
	// Apply many subsitutions to get a very robust comparison
	// For instance, multiplications are expanded and equal terms grouped.
	// Operations on numbers are also performed.
	ExpandedSubstitutions
)

// AreExpressionsEquivalent compares the two expressions using
// mathematical knowledge, as hinted by `level`.
// For instance, (a+b)^2 and (a^2 + 2ab + b^2) are equivalent
// if level == ExpandedSubstitutions, but not with other levels.
func AreExpressionsEquivalent(e1, e2 *Expr, level ComparisonLevel) bool {
	if level == Strict {
		return e1.equals(e2)
	}

	// start by evaluating
	v1, err1 := e1.Evaluate(nil)
	v2, err2 := e2.Evaluate(nil)
	if err1 == nil && err2 == nil && AreFloatEqual(v1, v2) {
		return true
	}

	e1, e2 = e1.Copy(), e2.Copy() // make sur e1 and e2 are not mutated
	if level == SimpleSubstitutions {
		e1.basicSimplification()
		e2.basicSimplification()
	} else {
		e1.fullSimplification()
		e2.fullSimplification()
	}

	return e1.equals(e2)
}

// partial evaluation a.k.a substitution

// Substitute replaces variables contained in `vars`, returning
// a copy.
func (expr *Expr) Substitute(vars Vars) {
	if expr == nil {
		return
	}
	expr.left.Substitute(vars)
	expr.right.Substitute(vars)

	if v, isVariable := expr.atom.(Variable); isVariable {
		value, has := vars[v]
		if has {
			*expr = *value.Copy()
		}
	}
}
