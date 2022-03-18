package expression

import (
	"sort"
)

// return a slice of all the operands of the `op` operator at the level of `expr`,
// or a one length slice
func (expr *Expression) extractOperator(op operator) []*Expression {
	if expr == nil {
		return nil
	}

	var out []*Expression
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
func compareNodes(n1, n2 *Expression) int {
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
		case random:
			a2 := a2.(random)
			if a1.isPrime && !a2.isPrime {
				return 1
			} else if !a1.isPrime && a2.isPrime {
				return -1
			} else {
				if a1.start == a2.start {
					return a1.end - a2.end
				} else {
					return a1.start - a2.start
				}
			}
		case Variable:
			a2 := a2.(Variable)
			if a1 < a2 {
				return -1
			} else if a1 > a2 {
				return 1
			} else {
				return 0
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
func nodesAsTree(nodes []*Expression, op operator) Expression {
	expr := nodes[0]
	for _, n := range nodes[1:] {
		expr = &Expression{atom: op, left: expr, right: n}
	}
	return *expr
}

// use associativity and commutativity to reorder (in place) + and * operations
// in a "cannonical" order
func (expr *Expression) sortPlusAndMultOperands() {
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
func (expr *Expression) expandMult() {
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
		expr.left = &Expression{atom: mult, left: a, right: c}
		expr.right = &Expression{atom: mult, left: b, right: c}
	} else if expr.right.atom == plus { // c * (a+b) => c*a + c*b
		a, b := expr.right.left, expr.right.right
		c := expr.left
		expr.atom = plus
		expr.left = &Expression{atom: mult, left: c, right: a}
		expr.right = &Expression{atom: mult, left: c, right: b}
	} // double distributivity is handled by recursion
}

// replace integral power by product, so that
// expandMult can trigger
func (expr *Expression) expandPow() {
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
			exprNew := expr.left.copy()
			for k := 1; k < asInt; k++ {
				exprNew = &Expression{atom: mult, left: exprNew, right: expr.left.copy()}
			}
			*expr = *exprNew
		}
	}
}

// replace a + a + a by 3*a
// should be applied after sorting operands
func (expr *Expression) groupAdditions() {
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
		newNodes []*Expression
		count    = 1
		ref      = nodes[0]
	)
	for _, n := range nodes[1:] {
		if ref.equals(n) {
			count++
		} else { // combine the nodes
			if count > 1 {
				newNodes = append(newNodes, &Expression{
					atom:  mult,
					left:  newNumber(float64(count)),
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
		newNodes = append(newNodes, &Expression{
			atom:  mult,
			left:  newNumber(float64(count)),
			right: ref,
		})
	} else {
		newNodes = append(newNodes, ref)
	}

	combined := nodesAsTree(newNodes, plus)
	*expr = combined
}

// replace a-c  by a + (-c)
// so that plus operation may trigger
func (expr *Expression) expandMinus() {
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
		expr.right = &Expression{atom: minus, right: expr.right}
	}
}

const maxIterations = 10_000 // very very unlikely in pratice

func (expr *Expression) basicSimplification() (nbPasses int) {
	ref := expr.copy()

	// apply each transformation until no one triggers a change
	for nbPasses = 1; nbPasses < maxIterations; nbPasses++ {
		expr.simplifyNumbers()
		expr.sortPlusAndMultOperands()
		expr.expandMinus()

		if expr.equals(ref) {
			break
		}
		ref = expr.copy() // update the reference and start a new pass
	}

	return nbPasses
}

func (expr *Expression) fullSimplification() (nbPasses int) {
	ref := expr.copy()

	// apply each transformation until no one triggers a change
	for nbPasses = 1; nbPasses < maxIterations; nbPasses++ {
		expr.simplifyNumbers()
		expr.expandPow()
		expr.expandMult()
		expr.sortPlusAndMultOperands()
		expr.groupAdditions()
		expr.expandMinus()

		if expr.equals(ref) {
			break
		}
		ref = expr.copy() // update the reference and start a new pass
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
	// For instance, x + 2 != 2 + x according to this level
	Strict ComparisonLevel = iota
	// Expressions are compared after performing some basic transformations,
	// such as reording operands.
	SimpleSubstitutions
	// Apply many subsitutions to get a very robust comparison
	// For instance, multiplications are expanded and equal terms grouped.
	ExpandedSubstitutions
)

// AreExpressionsEquivalent compares the two expressions using
// mathematical knowledge, as hinted by `level`.
// For instance, (a+b)^2 and (a^2 + 2ab + b^2) are equivalent
// if level == ExpandedSubstitutions, but not with other levels.
func AreExpressionsEquivalent(e1, e2 *Expression, level ComparisonLevel) bool {
	switch level {
	case Strict:
		// pass
	case SimpleSubstitutions:
		e1, e2 = e1.copy(), e2.copy()
		e1.basicSimplification()
		e2.basicSimplification()
	default:
		e1, e2 = e1.copy(), e2.copy()
		e1.fullSimplification()
		e2.fullSimplification()
	}

	return e1.equals(e2)
}
