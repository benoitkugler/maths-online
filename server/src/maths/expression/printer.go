package expression

import (
	"fmt"
	"math"
	"strings"
)

//go:generate go run unicode-latex/gen.go

// defaultLatexResolver maps from unicode values to LaTeX commands
// for usual maths symbols.
func defaultLatexResolver(v Variable) string { return unicodeToLaTeX[v.Name] }

func (mt matrix) defaultSimplify() {
	for i := range mt {
		for j := range mt[i] {
			mt[i][j].DefaultSimplify()
		}
	}
}

// DefaultSimplify performs a serie of basic simplifications,
// removing zero values and -(-...).
// It is used before printing an expression, and should also be used
// after substitutions.
// Note that it applies a subset of the simplifications defined by
// the [SimpleSimplifications] flag.
func (expr *Expr) DefaultSimplify() {
	if expr == nil {
		return
	}

	if mat, ok := expr.atom.(matrix); ok {
		mat.defaultSimplify()
		return
	}

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
	expr.DefaultSimplify()

	return expr.atom.asLaTeX(expr.left, expr.right)
}

// String returns a human readable form of the expression.
// The expression is prettified, meaning the structure of the
// returned expression may slightly differ, but is garanteed
// to be mathematically equivalent.
// See `AsLaTeX` for a better display format.
func (expr *Expr) String() string {
	expr = expr.Copy()
	expr.DefaultSimplify()
	return expr.Serialize()
}

func (v Variable) String() string {
	// we have to output valid expression syntax
	if v.Name == 0 { // custom symbols
		return fmt.Sprintf(`"%s"`, v.Indice)
	}

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

	leftHasParenthesis := op.needParenthesis(left, true, true)
	rightHasParenthesis := op.needParenthesis(right, false, true)

	if leftHasParenthesis {
		leftCode = addParenthesisLatex(leftCode)
	}
	if rightHasParenthesis {
		rightCode = addParenthesisLatex(rightCode)
	}

	switch op {
	case equals:
		return fmt.Sprintf("%s = %s", leftCode, rightCode)
	case greater:
		return fmt.Sprintf("%s \\ge %s", leftCode, rightCode)
	case strictlyGreater:
		return fmt.Sprintf("%s > %s", leftCode, rightCode)
	case lesser:
		return fmt.Sprintf("%s \\le %s", leftCode, rightCode)
	case strictlyLesser:
		return fmt.Sprintf("%s < %s", leftCode, rightCode)
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
		// special case to handle -inf ; avoid -+inf
		if strings.HasPrefix(rightCode, latexPlusInf) { // remove the plus
			rightCode = rightCode[1:]
		}
		if leftCode == "" { // remove the space
			return fmt.Sprintf("-%s", rightCode)
		}
		return fmt.Sprintf("%s - %s", leftCode, rightCode)
	case mult:
		// check for implicit multiplication
		if shouldOmitTimes(left, rightHasParenthesis, right) {
			return fmt.Sprintf(`%s %s`, leftCode, rightCode)
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
	case factorial:
		return fmt.Sprintf(`%s!`, leftCode)
	case union:
		return fmt.Sprintf(`%s \cup %s`, leftCode, rightCode)
	case intersection:
		return fmt.Sprintf(`%s \cap %s`, leftCode, rightCode)
	case complement:
		return fmt.Sprintf(`\overline{%s}`, rightCode)
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
	case isPrimeFn:
		return fmt.Sprintf(`\text{isPrime}\left(%s\right)`, arg)
	case forceDecimalFn:
		return fmt.Sprintf(`\text{forceDecimal}\left(%s\right)`, arg)
	case detFn:
		// for explicit matrice, use do not add ()
		if _, isMat := right.atom.(matrix); isMat {
			return fmt.Sprintf(`\det %s`, arg)
		}
		return fmt.Sprintf(`\det\left(%s\right)`, arg)
	case traceFn:
		// for explicit matrice, use do not add ()
		if _, isMat := right.atom.(matrix); isMat {
			return fmt.Sprintf(`\text{trace} %s`, arg)
		}
		return fmt.Sprintf(`\text{trace}\left(%s\right)`, arg)
	case invertFn:
		// add parenthesis on op
		switch right.atom.(type) {
		case operator, function:
			return fmt.Sprintf(`\left(%s\right)^{-1}`, arg)
		}
		return fmt.Sprintf(`%s^{-1}`, arg)
	case transposeFn:
		// add parenthesis on op
		switch right.atom.(type) {
		case operator, function:
			return fmt.Sprintf(`\left(%s\right)^{T}`, arg)
		}
		return fmt.Sprintf(`%s^{T}`, arg)
	default:
		panic(exhaustiveFunctionSwitch)
	}
}

func (i indice) asLaTeX(left, right *Expr) string {
	return i.serialize(left, right) // the syntaxe actually matches latex
}

func (r roundFunc) asLaTeX(_, right *Expr) string {
	return fmt.Sprintf(`\text{round(%s; %d)}`, right.AsLaTeX(), r.nbDigits)
}

// try to write the sequence explicitely
func (r specialFunction) expandSequence(eval bool) (string, error) {
	start, err := evalInt(r.args[1], nil)
	if err != nil {
		return "", err
	}
	end, err := evalInt(r.args[2], nil)
	if err != nil {
		return "", err
	}
	k, _ := r.args[0].atom.(Variable)
	expr := r.args[3]
	var chunks []string
	for i := start; i <= end; i++ {
		v := newRealInt(i)
		term := expr.Copy()
		term.Substitute(Vars{k: v.toExpr()})
		if eval {
			term.reduce()
		}
		chunks = append(chunks, term.AsLaTeX())
	}

	var sep string
	switch r.kind {
	case sumFn:
		sep = " + "
	case prodFn:
		sep = " \\times "
	case unionFn:
		sep = " \\cup "
	case interFn:
		sep = " \\cap "
	}
	return strings.Join(chunks, sep), nil
}

func (r specialFunction) asLaTeX(_, _ *Expr) string {
	_ = exhaustiveSpecialFunctionSwitch
	switch r.kind {
	case binomial:
		k, n := r.args[0], r.args[1]
		return fmt.Sprintf(`\binom{%s}{%s}`, n.AsLaTeX(), k.AsLaTeX())
	case sumFn, prodFn, unionFn, interFn:
		k, start, end, expr := r.args[0], r.args[1], r.args[2], r.args[3]
		if len(r.args) == 5 {
			switch r.args[4].atom {
			case Variable{Indice: "expand-eval"}:
				// try to write the sequence explicitely
				code, err := r.expandSequence(true)
				if err == nil {
					return code
				}
				fallthrough // try only expanding
			case Variable{Indice: "expand"}:
				// try to write the sequence explicitely
				code, err := r.expandSequence(false)
				if err == nil {
					return code
				}
			}
			// else, default to general notation
		}
		var latexOp string
		switch r.kind {
		case sumFn:
			latexOp = "sum"
		case prodFn:
			latexOp = "prod"
		case unionFn:
			latexOp = "bigcup"
		case interFn:
			latexOp = "bigcap"
		}
		return fmt.Sprintf(`\%s_{%s=%s}^{%s} %s`, latexOp, k.AsLaTeX(), start.AsLaTeX(), end.AsLaTeX(), expr.AsLaTeX())
	default:
		return fmt.Sprintf(`\text{%s}`, r.String())
	}
}

func (v Variable) asLaTeX(_, _ *Expr) string {
	// special case for "" variable to allow custom symbols
	if v.Name == 0 {
		return v.Indice
	}

	name := defaultLatexResolver(v)
	if v.Indice != "" {
		name += "_{" + v.Indice + "}"
	}
	return name
}

func (mt matrix) asLaTeX(_, _ *Expr) string {
	rows := make([]string, len(mt))
	for i, row := range mt {
		cols := make([]string, len(row))
		for j, col := range row {
			cols[j] = col.AsLaTeX()
		}
		rows[i] = strings.Join(cols, " & ")
	}
	return fmt.Sprintf(`\begin{pmatrix} 
		%s
	\end{pmatrix}`, strings.Join(rows, "\\\\\n"))
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

const latexPlusInf = "+\\infty"

func (v Number) asLaTeX(_, _ *Expr) string {
	if math.IsInf(float64(v), 1) {
		return latexPlusInf
	} else if math.IsInf(float64(v), -1) {
		return "-\\infty"
	}
	return v.String()
}

// returns `true` is the expression is compound and requires parenthesis
// when used with `op`
// if `isLeft` is true, this is :  expr op ...
// else this is :                  ...  op expr
func (op operator) needParenthesis(expr *Expr, isLeftArg, isLaTex bool) bool {
	if expr == nil {
		return false
	}

	// the latex syntax allow to spare some redundant parenthesis
	if isLaTex && (op == div || op == rem || op == mod || op == complement) {
		return false
	}
	if isLaTex && !isLeftArg && op == pow {
		return false
	}

	switch atom := expr.atom.(type) {
	case Number, constant, function, Variable, roundFunc, specialFunction, indice, matrix:
		return false
	case operator:
		if isLaTex && atom == complement {
			// the latex overline is enough to delimit
			return false
		}

		_ = exhaustiveOperatorSwitch
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
		case mult:
			// the parser recognize 1/2*3 as (1/2)*3, as
			// calculator and programing languages do, but
			// we usually add parenthesis for clarity

			// (1/2) * y, 3 * (1/2), 1/2 + 3, 2 - 1/2
			if !isLaTex && isLeftArg && atom == div {
				return true
			}
			return op > atom
		case factorial:
			// parenthesis are always needed if the argument is an expression
			return true
		case union, intersection, complement:
			return true
		default:
			return op > atom //  (1+2)*4, (1-2)*4, 4*(1+2)
		}
	default:
		panic(exhaustiveAtomSwitch)
	}
}

// plain text, pretty output

func (op operator) serialize(left, right *Expr) string {
	leftCode, rightCode := left.Serialize(), right.Serialize()

	if leftHasParenthesis := op.needParenthesis(left, true, false); leftHasParenthesis {
		leftCode = addParenthesisPlain(leftCode)
	}

	rightHasParenthesis := op.needParenthesis(right, false, false)
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
		if shouldOmitTimes(left, rightHasParenthesis, right) {
			return fmt.Sprintf(`%s%s`, leftCode, rightCode)
		}
		return fmt.Sprintf(`%s * %s`, leftCode, rightCode)
	case factorial:
		return fmt.Sprintf(`%s!`, leftCode)
	case div:
		// compact fractions
		return leftCode + op.String() + rightCode
	case equals, greater, strictlyGreater, lesser, strictlyLesser, union, intersection, complement,
		mod, rem, pow:
		return fmt.Sprintf(`%s %s %s`, leftCode, op.String(), rightCode)
	default:
		panic(exhaustiveOperatorSwitch)
	}
}

// left and right are assumed not to be nil
// we are conservative here : only omit * when we are
// certain it results in valid ouput
func shouldOmitTimes(left *Expr, rightHasParenthesis bool, right *Expr) bool {
	// with parenthesis, it is always valid to omit times
	if rightHasParenthesis {
		return true
	}

	switch right.atom.(type) {
	case Variable, constant, function, specialFunction, roundFunc, indice:
		return true
	case Number:
		return false
	case matrix:
		_, leftIsMatrix := left.atom.(matrix)
		return !leftIsMatrix
	case operator:
		// if the first term is not a number, it is safe to remove the *
		// but otherwise not : 2 * 4 ^ 2
		if right.left != nil {
			if _, isNumber := right.left.atom.(Number); !isNumber {
				return true
			}
		}
		return false
	default:
		panic(exhaustiveAtomSwitch)
	}
}
