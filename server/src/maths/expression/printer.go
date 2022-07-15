package expression

import (
	"fmt"
	"strings"
)

//go:generate go run unicode-latex/gen.go

// defaultLatexResolver maps from unicode values to LaTeX commands
// for usual maths symbols.
func defaultLatexResolver(v Variable) string { return unicodeToLaTeX[v.Name] }

func (expr *Expr) simplifyForPrint() {
	for nbPasses := 0; nbPasses < 2; nbPasses++ {
		expr.normalizeNegativeNumbers()
		expr.extractNegativeInMults()
		expr.contractPlusMinus()
		expr.contractMinusMinus()
		expr.simplify0And1()
	}
}

// AsLaTeX returns a valid LaTeX code displaying the expression.
func (expr *Expr) AsLaTeX() string {
	if expr == nil {
		return ""
	}

	expr = expr.Copy()
	expr.simplifyForPrint()

	return expr.atom.asLaTeX(expr.left, expr.right)
}

// String returns a human readable form of the expression.
// The expression is prettified, meaning the structure of the
// returned expression may slightly differ, but is garanteed
// to be mathematically equivalent.
// See `AsLaTeX` for a better display format.
func (expr *Expr) String() string {
	expr = expr.Copy()
	expr.simplifyForPrint()
	return expr.Serialize()
}

func (v Variable) String() string {
	// we have to output valid expression syntax
	out := string(v.Name)
	if v.Indice != "" {
		return out + "_" + v.Indice + " " // notice the white space to avoid x_Ay_A
	}
	return out
}

func addParenthesisLatex(s string) string {
	return `\left(` + s + `\right)`
}

func addParenthesisPlain(s string) string {
	return `(` + s + `)`
}

func (op operator) asLaTeX(left, right *Expr) string {
	leftCode, rightCode := left.AsLaTeX(), right.AsLaTeX()

	leftHasParenthesis := op.needParenthesis(left, true)
	rightHasParenthesis := op.needParenthesis(right, false)

	// the latex syntax allow to spare some redundant parenthesis
	ignoreParenthesis := op == div || op == rem || op == mod

	if leftHasParenthesis && !ignoreParenthesis {
		leftCode = addParenthesisLatex(leftCode)
	}
	if rightHasParenthesis && !ignoreParenthesis {
		rightCode = addParenthesisLatex(rightCode)
	}

	switch op {
	case plus:
		if leftCode == "" {
			return rightCode // plus is optional
		}
		// special case when rightCode starts with a - (without parenthesis),
		// such as in x + (-a + 2)
		if strings.HasPrefix(rightCode, "-") { // add the space back
			return fmt.Sprintf("%s - %s", leftCode, rightCode[1:])
		}
		return fmt.Sprintf("%s + %s", leftCode, rightCode)
	case minus:
		if leftCode == "" { // remove the space
			return fmt.Sprintf("-%s", rightCode)
		}
		return fmt.Sprintf("%s - %s", leftCode, rightCode)
	case mult:
		// check for implicit multiplication
		if shouldOmitTimes(rightHasParenthesis, right) {
			return fmt.Sprintf(`%s%s`, leftCode, rightCode)
		}
		return fmt.Sprintf(`%s \times %s`, leftCode, rightCode)
	case div:
		return fmt.Sprintf(`\frac{%s}{%s}`, leftCode, rightCode)
	case mod:
		return fmt.Sprintf(`\text{mod}\left(%s; %s\right)`, leftCode, rightCode)
	case rem:
		return fmt.Sprintf(`\text{rem}\left(%s; %s\right)`, leftCode, rightCode)
	case pow:
		return fmt.Sprintf(`{%s}^{%s}`, leftCode, rightCode)
	default:
		panic(exhaustiveOperatorSwitch)
	}
}

func (fn function) asLaTeX(left, right *Expr) string {
	arg := right.AsLaTeX()
	switch fn {
	case logFn:
		return fmt.Sprintf(`\log\left(%s\right)`, arg)
	case expFn:
		return fmt.Sprintf(`\exp\left(%s\right)`, arg)
	case sinFn:
		return fmt.Sprintf(`\sin\left(%s\right)`, arg)
	case cosFn:
		return fmt.Sprintf(`\cos\left(%s\right)`, arg)
	case tanFn:
		return fmt.Sprintf(`\tan\left(%s\right)`, arg)
	case asinFn:
		return fmt.Sprintf(`\arcsin\left(%s\right)`, arg)
	case acosFn:
		return fmt.Sprintf(`\arccos\left(%s\right)`, arg)
	case atanFn:
		return fmt.Sprintf(`\arctan\left(%s\right)`, arg)
	case absFn:
		return fmt.Sprintf(`\left|%s\right|`, arg)
	case floorFn:
		return fmt.Sprintf(`\left\lfloor %s \right\rfloor`, arg)
	case sqrtFn:
		return fmt.Sprintf(`\sqrt{%s}`, arg)
	case sgnFn:
		return fmt.Sprintf(`\text{sgn}\left(%s\right)`, arg)
	case isZeroFn:
		return fmt.Sprintf(`\text{isZero}\left(%s\right)`, arg)
	case isPrimeFn:
		return fmt.Sprintf(`\text{isPrime}\left(%s\right)`, arg)
	default:
		panic(exhaustiveFunctionSwitch)
	}
}

func (r roundFn) asLaTeX(_, right *Expr) string {
	return fmt.Sprintf(`\text{round(%s; %d)}`, right.AsLaTeX(), r.nbDigits)
}

func (r specialFunctionA) asLaTeX(_, _ *Expr) string {
	return fmt.Sprintf(`\text{%s}`, r.String())
}

func (r randVariable) asLaTeX(_, _ *Expr) string {
	return fmt.Sprintf(`\text{%s}`, r.String())
}

func (v Variable) asLaTeX(_, _ *Expr) string {
	// special case for @_variable to allow custom symbols
	if v.Name == '@' {
		return v.Indice
	}

	name := defaultLatexResolver(v)
	if v.Indice != "" {
		name += "_{" + v.Indice + "}"
	}
	return name
}

func (c constant) asLaTeX(_, _ *Expr) string {
	switch c {
	case piConstant:
		return "\\pi"
	case eConstant:
		return "e"
	default:
		panic(exhaustiveConstantSwitch)
	}
}

func (v Number) asLaTeX(_, _ *Expr) string { return v.String() }

// returns `true` is the expression is compound and requires parenthesis
// when used with `op`
// if `isLeft` is true, this is :  expr op ...
// else this is :                  ...  op expr
func (op operator) needParenthesis(expr *Expr, isLeftArg bool) bool {
	if expr == nil {
		return false
	}

	switch atom := expr.atom.(type) {
	case Number, constant, function, Variable, randVariable, roundFn, specialFunctionA:
		return false
	case operator:
		switch op {
		case minus:
			if isLeftArg { // actually rever required
				return false
			}
			return atom <= op // - or +
		case pow:
			if atom == pow {
				return isLeftArg // (2^3)^5, but 2^(3^5) is usually written 2^3^5
			}
			return atom < op
		default:
			return atom < op
		}
	default:
		panic(exhaustiveAtomSwitch)
	}
}

// plain text, pretty output

func (op operator) serialize(left, right *Expr) string {
	leftCode, rightCode := left.Serialize(), right.Serialize()

	if leftHasParenthesis := op.needParenthesis(left, true); leftHasParenthesis {
		leftCode = addParenthesisPlain(leftCode)
	}

	rightHasParenthesis := op.needParenthesis(right, false)
	if rightHasParenthesis {
		rightCode = addParenthesisPlain(rightCode)
	}

	switch op {
	case plus:
		if leftCode == "" {
			return rightCode // plus is optional
		}
		// special case when rightCode starts with a - (without parenthesis),
		// such as in x + (-a + 2)
		if strings.HasPrefix(rightCode, "-") { // add the space back
			return fmt.Sprintf("%s - %s", leftCode, rightCode[1:])
		}
		return fmt.Sprintf("%s + %s", leftCode, rightCode)
	case minus:
		if leftCode == "" { // remove the space
			return fmt.Sprintf("-%s", rightCode)
		}
		return fmt.Sprintf("%s - %s", leftCode, rightCode)
	case mult:
		// check for implicit multiplication
		if shouldOmitTimes(rightHasParenthesis, right) {
			return fmt.Sprintf(`%s%s`, leftCode, rightCode)
		}
		return fmt.Sprintf(`%s * %s`, leftCode, rightCode)
	case div:
		return fmt.Sprintf(`%s / %s`, leftCode, rightCode)
	case mod:
		return fmt.Sprintf(`%s %% %s`, leftCode, rightCode)
	case rem:
		return fmt.Sprintf(`%s // %s`, leftCode, rightCode)
	case pow:
		return fmt.Sprintf(`%s ^ %s`, leftCode, rightCode)
	default:
		panic(exhaustiveOperatorSwitch)
	}
}

// left and right are assumed not to be nil
// we are conservative here : only omit * when we are
// certain it results in valid ouput
func shouldOmitTimes(rightHasParenthesis bool, right *Expr) bool {
	// with parenthesis, it is always valid to omit times
	if rightHasParenthesis {
		return true
	}

	switch right.atom.(type) {
	case Variable, constant, function, specialFunctionA, roundFn:
		return true
	case operator:
		// if the first term is not a number, it is safe to remove the *
		// but otherwise not : 2 * 4 ^ 2
		if right.left != nil {
			if _, isNumber := right.left.atom.(Number); !isNumber {
				return true
			}
		}
	}
	return false
}
