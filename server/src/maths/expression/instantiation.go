package expression

import (
	"errors"
	"fmt"
	"math/rand"
)

// This file implements the instantiation/resolution of a set of variable/expression pairs.
// The general idea is to replace known variable by their expression (recursively)
// and to perform a partial evaluation.
// The selectors involved in random generators are evaluated, and the resulting
// expression are evaluated if valid.
// The result may be a plain number, but also a symbol, or even a mix of expression.

// Validate calls `Instantiate` many times to make sure the parameters are always
// valid regardless of the random value chosen.
// If not, it returns the first error encountered.
func (rv RandomParameters) Validate() error {
	const nbTries = 200
	for i := 0; i < nbTries; i++ {
		_, err := rv.Instantiate()
		if err != nil {
			return err
		}
	}
	return nil
}

// ErrInvalidRandomParameters is returned when instantiating
// invalid parameter definitions
type ErrInvalidRandomParameters struct {
	Detail string
	Cause  Variable
}

func (irv ErrInvalidRandomParameters) Error() string {
	return fmt.Sprintf("%s -> %s", irv.Cause, irv.Detail)
}

// Instantiate generates a random version of the variables, resolving possible dependencies.
//
// The general idea is to replace known variables by their expression (recursively)
// and to perform a partial evaluation.
// The selectors involved in random generators are evaluated, and the resulting
// expression are evaluated if valid. However, free variables are accepted, so that
// the result may be a plain number, but also a symbol, or even a mix of expression.
//
// It returns an `ErrInvalidRandomParameters` error for invalid cycles, like a = a + 1
// or a = b + 1; b = a.
//
// See `Validate` to statistically check for invalid parameters.
func (rv RandomParameters) Instantiate() (Vars, error) {
	resolver := &paramsInstantiater{
		defs:    rv,
		seen:    make(map[Variable]bool),
		results: make(map[Variable]*Expr),
	}

	out := make(Vars, len(rv))
	for v := range rv {
		resolver.currentVariable = v
		value, err := resolver.instantiate(v) // this triggers the evaluation of the expression
		if err != nil {
			return nil, err
		}
		out[v] = value
	}

	return out, nil
}

type paramsInstantiater struct {
	defs RandomParameters

	seen    map[Variable]bool  // variable that we are currently resolving
	results map[Variable]*Expr // resulting values being built

	// the top level variable being resolved,
	// or zero if are recursing in the tree
	currentVariable Variable
}

// instantiate instantiate the given variable [v], and its dependencies
// if [v] is not defined, it returns an error
func (rvv *paramsInstantiater) instantiate(v Variable) (*Expr, error) {
	// first, check if it has already been resolved by side effect
	if expr, has := rvv.results[v]; has {
		return expr, nil
	}

	expr, ok := rvv.defs[v]
	if !ok {
		return nil, ErrInvalidRandomParameters{
			Cause:  v,
			Detail: fmt.Sprintf("%s n'est pas définie", v),
		}
	}

	// avoid invalid cycles
	if rvv.seen[v] {
		return nil, ErrInvalidRandomParameters{
			Cause:  v,
			Detail: fmt.Sprintf("%s est présente dans un cycle et ne peut donc pas être résolue.", v),
		}
	}

	// start the resolution : to detect invalid cycles, register the variable
	rvv.seen[v] = true

	// instantiate the definition, recursing when needed
	value, err := expr.instantiate(rvv)
	if err != nil {
		return nil, ErrInvalidRandomParameters{
			Cause:  v,
			Detail: err.Error(),
		}
	}

	value.DefaultSimplify()

	// register the result
	rvv.results[v] = value

	return value, nil
}

func (rvv *paramsInstantiater) resolve(v Variable) (real, error) {
	expr, err := rvv.instantiate(v)
	if err != nil {
		return real{}, err
	}
	return expr.evalReal(rvv)
}

func (expr *Expr) tryEval(ctx *paramsInstantiater) *Expr {
	if mat, ok := expr.atom.(matrix); ok {
		out := make(matrix, len(mat))
		for i, row := range mat {
			out[i] = make([]*Expr, len(row))
			for j := range row {
				v, err := mat[i][j].evalReal(ctx)
				if err == nil {
					out[i][j] = v.toExpr()
				} else {
					out[i][j] = mat[i][j]
				}
			}
		}
		return &Expr{atom: out}
	}
	v, err := expr.evalReal(ctx)
	if err == nil {
		return v.toExpr()
	}
	return expr
}

// we allow zero length cycles so that f = randSymbol(f ; g ; h) is valid
// note that such cycle will be invalid if used in evaluation
// (see Evaluate)
func (expr *Expr) isZeroCycle(currentVariable Variable) bool {
	if v, isVar := expr.atom.(Variable); isVar && v == currentVariable {
		return true
	}
	return false
}

// instantiate recurse through [expr], resolving defined variables
// and trying to evaluate all expressions. When used in random selectors,
// an error is returned on failure, but not in general
func (expr *Expr) instantiate(ctx *paramsInstantiater) (*Expr, error) {
	if expr == nil {
		return nil, nil
	}

	currentVariable := ctx.currentVariable
	ctx.currentVariable = Variable{} // when recursing, set tracker to zero

	switch atom := expr.atom.(type) {
	case Number, constant, roundFunc, indice: // no-op, simply recurse
		left, err := expr.left.instantiate(ctx)
		if err != nil {
			return nil, err
		}
		right, err := expr.right.instantiate(ctx)
		if err != nil {
			return nil, err
		}
		out := &Expr{atom: atom, left: left, right: right}
		return out.tryEval(ctx), nil
	case function: // we expand matrices functions
		right, err := expr.right.instantiate(ctx)
		if err != nil {
			return nil, err
		}
		var out *Expr
		switch atom {
		case traceFn:
			mat, ok := right.atom.(matrix)
			if !ok {
				return nil, errors.New("La fonction trace attend une matrice en argument.")
			}
			out, err = mat.trace()
			if err != nil {
				return nil, err
			}
		case detFn:
			mat, ok := right.atom.(matrix)
			if !ok {
				return nil, errors.New("La fonction trace attend une matrice en argument.")
			}
			matN, ok := newNumberMatrixFrom(mat)
			if !ok {
				return nil, errors.New("Le déterminant ne supporte pas le calcul symbolique.")
			}
			d, err := matN.determinant()
			if err != nil {
				return nil, err
			}
			return newNb(d), nil
		case transposeFn:
			mat, ok := right.atom.(matrix)
			if !ok {
				return nil, errors.New("La fonction trace attend une matrice en argument.")
			}
			mat = mat.transpose()
			out = &Expr{atom: mat}
		case invertFn:
			mat, ok := right.atom.(matrix)
			if !ok {
				return nil, errors.New("La fonction trace attend une matrice en argument.")
			}
			matN, ok := newNumberMatrixFrom(mat)
			if !ok {
				return nil, errors.New("L'inverse d'une matrice ne supporte pas le calcul symbolique.")
			}
			matN, err = matN.invert()
			if err != nil {
				return nil, err
			}
			out = &Expr{atom: matN.toExprMatrix()}
		default:
			out = &Expr{atom: atom, left: nil, right: right}
		}
		return out.tryEval(ctx), nil
	case operator: // we expand matrices calculus
		// recurse
		left, err := expr.left.instantiate(ctx)
		if err != nil {
			return nil, err
		}
		right, err := expr.right.instantiate(ctx)
		if err != nil {
			return nil, err
		}
		out, err := matricesOperation(atom, left, right, ctx)
		if err != nil {
			return nil, err
		}
		return out.tryEval(ctx), nil
	case Variable:
		// if the variable is not defined, just returns the free variable as an expression
		if _, isDefined := ctx.defs[atom]; !isDefined {
			return NewVarExpr(atom), nil
		}
		return ctx.instantiate(atom)
	case matrix:
		mt, err := atom.instantiate(ctx)
		if err != nil {
			return nil, err
		}
		return &Expr{atom: mt}, nil
	case specialFunction:
		// generate random numbers
		switch atom.kind {
		case randInt, randPrime, randDenominator:
			v, err := atom.evalRat(ctx)
			return v.toExpr(), err
		case randChoice:
			index := rand.Intn(len(atom.args))
			choice := atom.args[index]
			if choice.isZeroCycle(currentVariable) {
				// do not instantiate, resulting in cycle
				return choice, nil
			}
			return choice.instantiate(ctx)
		case choiceFrom:
			// we must evaluate the selector
			choice, err := choiceFromSelect(atom.args, ctx)
			if err != nil {
				return nil, err
			}
			if choice.isZeroCycle(currentVariable) {
				// do not instantiate, resulting in cycle
				return choice, nil
			}
			return choice.instantiate(ctx)
		case minFn, maxFn, matCoeff: // no-op, simply recurse
			inst := specialFunction{
				kind: atom.kind,
				args: make([]*Expr, len(atom.args)),
			}
			for i, arg := range atom.args {
				var err error
				inst.args[i], err = arg.instantiate(ctx)
				if err != nil {
					return nil, err
				}
			}
			out := &Expr{atom: inst}
			return out.tryEval(ctx), nil
		default:
			panic(exhaustiveSpecialFunctionSwitch)
		}
	default:
		panic(exhaustiveAtomSwitch)
	}
}
