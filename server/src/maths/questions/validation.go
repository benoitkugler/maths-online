package questions

import (
	"errors"
	"fmt"

	"github.com/benoitkugler/maths-online/server/src/maths/expression"
)

// maxFunctionBound is the maximum value a function
// may reached. Higher values are either a bug, or won't be properly
// displayed on the student client
const maxFunctionBound = 100

type ErrParameters struct {
	Origin  string
	Details string
}

func (err ErrParameters) Error() string {
	return fmt.Sprintf("invalid random parameters in %s: %s", err.Origin, err.Details)
}

// Validate ensure the given `Parameters` are sound,
// by parsing the expression, checking for duplicate parameters,
// and detecting definition cycles.
// It the error is not nil, it will be of type `ErrParameters`.
// Once called without error, `ToMap` may be safely used.
func (pr Parameters) Validate() error {
	params := make(expression.RandomParameters)
	for _, item := range pr {
		err := item.mergeTo(params)
		if err != nil {
			return ErrParameters{
				Origin:  item.String(),
				Details: err.Error(),
			}
		}
	}

	for v := range params {
		if v.Name == 'e' {
			return ErrParameters{
				Origin:  v.String(),
				Details: "La variable e n'est pas autorisée (car utilisée pour exp).",
			}
		}
	}

	err := params.Validate()
	if err != nil {
		return ErrParameters{
			Origin:  "Paramètres aléatoires",
			Details: err.Error(),
		}
	}

	return nil
}

type errEnonce struct {
	Error string            // detailed error
	Block int               // index of the invalid block
	Vars  map[string]string // the actual values used when the error was encountered, or nil
}

type ErrorKind uint8

const (
	ErrParameters_ ErrorKind = iota
	ErrEnonce
	ErrCorrection
)

// ErrQuestionInvalid is returned by  Question.Validate()
// It is either an error about the random parameters, or the blocks content (enonce or correction).
type ErrQuestionInvalid struct {
	ErrParameters ErrParameters
	ErrEnonce     errEnonce
	ErrCorrection errEnonce
	Kind          ErrorKind // indicates which field is valid
}

func (e ErrQuestionInvalid) Error() string {
	switch e.Kind {
	case ErrParameters_:
		return fmt.Sprintf("invalid question parameters: %v", e.ErrParameters)
	case ErrEnonce:
		return fmt.Sprintf("invalid question blocks: %v", e.ErrEnonce)
	case ErrCorrection:
		return fmt.Sprintf("invalid correction blocks: %v", e.ErrCorrection)
	default:
		panic("exhaustive switch")
	}
}

func (en Enonce) validate(params expression.RandomParameters) (bool, errEnonce) {
	// setup the validators
	var err error
	validators := make([]validator, len(en))
	for i, block := range en {
		validators[i], err = block.setupValidator(params)
		if err != nil {
			return false, errEnonce{Block: i, Error: err.Error()}
		}
	}

	const nbTries = 1_000
	for try := 0; try < nbTries; try++ {
		// instantiate the parameters for this try
		vars, _ := params.Instantiate()

		// run through the blocks
		for i, v := range validators {
			err := v.validate(vars)
			if err != nil {
				// export the current parameters as strings
				varsS := make(map[string]string, len(vars))
				for k, v := range vars {
					varsS[k.String()] = v.String()
				}

				return false, errEnonce{Block: i, Error: err.Error(), Vars: varsS}
			}
		}
	}

	return true, errEnonce{}
}

// Validate ensure the random parameters and enonce blocks are sound.
// If not, an `ErrQuestionInvalid` is returned.
func (qu QuestionPage) Validate() error {
	if err := qu.Parameters.Validate(); err != nil {
		return ErrQuestionInvalid{Kind: ErrParameters_, ErrParameters: err.(ErrParameters)}
	}

	params := qu.Parameters.ToMap()

	if ok, err := qu.Enonce.validate(params); !ok {
		return ErrQuestionInvalid{Kind: ErrEnonce, ErrEnonce: err}
	}

	if ok, err := qu.Correction.validate(params); !ok {
		return ErrQuestionInvalid{Kind: ErrCorrection, ErrCorrection: err}
	}

	return nil
}

type validator interface {
	// validate the field given the instantiated values
	validate(vars expression.Vars) error
}

type noOpValidator struct{}

func (noOpValidator) validate(vars expression.Vars) error { return nil }

type parsedCoord struct {
	X, Y *expression.Expr
}

func (c parsedCoord) validate(vars expression.Vars, checkPrecision bool) error {
	if err := c.X.IsValidNumber(vars, checkPrecision, true); err != nil {
		return err
	}
	if err := c.Y.IsValidNumber(vars, checkPrecision, true); err != nil {
		return err
	}
	return nil
}

type linearEquationValidator struct {
	expr *expression.Expr
}

func (v linearEquationValidator) validate(vars expression.Vars) error {
	return v.expr.IsValidLinearEquation(vars)
}

type variationTableValidator struct {
	label TextParts
	xs    []*expression.Expr
	fxs   []*expression.Expr
}

func (v variationTableValidator) validate(vars expression.Vars) error {
	for _, c := range v.fxs {
		err := c.IsValidNumber(vars, false, true)
		if err != nil {
			return err
		}
	}

	return expression.AreSortedNumbers(v.xs, vars)
}

type figureValidator struct {
	pointNames []*expression.Expr
	points     []*expression.Expr // X,Y
	references []*expression.Expr

	circlesDims []*expression.Expr // center and radius

	lines [][2]*expression.Expr // A, B
}

func (v figureValidator) pointStrings(vars expression.Vars) map[string]bool {
	out := make(map[string]bool, len(v.pointNames))
	for _, expr := range v.pointNames {
		expr = expr.Copy()
		expr.Substitute(vars)
		out[expr.AsLaTeX()] = true
	}
	return out
}

func (v figureValidator) validate(vars expression.Vars) error {
	for _, point := range v.points {
		if err := point.IsValidNumber(vars, false, true); err != nil {
			return err
		}
	}

	points := v.pointStrings(vars)

	// check for duplicates ...
	if len(points) != len(v.pointNames) {
		return errors.New("Les noms des points ne sont pas distincts.")
	}

	// .. and undefined points
	for _, ref := range v.references {
		ref = ref.Copy()
		ref.Substitute(vars)
		resolvedRef := ref.AsLaTeX()
		if hasPoint := points[resolvedRef]; !hasPoint {
			return fmt.Errorf("L'expression %s ne définit pas un point connu.", resolvedRef)
		}
	}

	// check for valid circle dimensions
	for _, circleNum := range v.circlesDims {
		if err := circleNum.IsValidNumber(vars, false, true); err != nil {
			return err
		}
	}

	// check for affine line coefficients
	for _, line := range v.lines {
		if err := line[0].IsValidNumber(vars, false, false); err != nil {
			return err
		}

		if err := line[1].IsValidNumber(vars, false, true); err != nil {
			return err
		}
	}

	return nil
}

type function struct {
	label  string
	domain expression.Domain
	expression.FunctionExpr
}

func newFunction(fn FunctionDefinition, params expression.RandomParameters) (function, error) {
	fnExpr, from, to, err := fn.parse()
	if err != nil {
		return function{}, err
	}

	// check that the function variable is not used
	if params[fn.Variable] != nil {
		return function{}, fmt.Errorf("La variable <b>%s</b> est déjà utilisée dans les paramètres aléatoires.", fn.Variable)
	}

	return function{label: fn.Decoration.Label, FunctionExpr: fnExpr, domain: expression.Domain{From: from, To: to}}, nil
}

type areaVData struct {
	top, bottom TextParts
	domain      expression.Domain
}

type functionPointVData struct {
	fnLabel TextParts
	x       *expression.Expr
}
type functionsGraphValidator struct {
	functions          []function
	variationValidator []variationTableValidator
	sequences          []function
	areas              []areaVData
	points             []functionPointVData
}

func (v functionsGraphValidator) validate(vars expression.Vars) error {
	for _, f := range v.functions {
		if err := f.FunctionExpr.IsValidAsFunction(f.domain, vars, maxFunctionBound); err != nil {
			return err
		}
	}
	for _, varTable := range v.variationValidator {
		if err := varTable.validate(vars); err != nil {
			return err
		}
	}
	for _, f := range v.sequences {
		if err := f.FunctionExpr.IsValidAsSequence(f.domain, vars, maxFunctionBound); err != nil {
			return err
		}
	}

	// checks that function with same label are defined on non overlapping intervals,
	// so that area references can't be ambiguous
	byNames := make(map[string][]expression.Domain)
	for _, fn := range v.functions {
		byNames[fn.label] = append(byNames[fn.label], fn.domain)
	}
	for _, vt := range v.variationValidator {
		label, err := vt.label.instantiateAndMerge(vars)
		if err != nil {
			return err
		}
		byNames[label] = append(byNames[label], expression.Domain{
			From: vt.xs[0], To: vt.xs[len(vt.xs)-1], // vt.validate checks that these calls are safe
		})
	}

	for name, domains := range byNames {
		if err := expression.AreDisjointsDomains(domains, vars); err != nil {
			return fmt.Errorf("Pour la fonction %s, %s.", name, err)
		}
	}

	// checks that areas are referencing known functions
	// and that the domains are valid
	for _, area := range v.areas {
		top, err := area.top.instantiateAndMerge(vars)
		if err != nil {
			return err
		}
		bottom, err := area.bottom.instantiateAndMerge(vars)
		if err != nil {
			return err
		}
		domainsTop := byNames[top]
		if top == "" {
			domainsTop = []expression.Domain{{}} // use the abscisse axis, which has no constraints
		}
		if len(domainsTop) == 0 {
			return fmt.Errorf("La fonction %s n'est pas définie.", top)
		}
		domainsBottom := byNames[bottom]
		if bottom == "" {
			domainsBottom = []expression.Domain{{}} // use the abscisse axis, which has no constraints
		}
		if len(domainsBottom) == 0 {
			return fmt.Errorf("La fonction %s n'est pas définie.", bottom)
		}

		// check that the domain in included in one of the domain for f1 and f2
		if err := area.domain.IsIncludedIntoOne(domainsTop, vars); err != nil {
			return err
		}
		if err := area.domain.IsIncludedIntoOne(domainsBottom, vars); err != nil {
			return err
		}
	}

	// check that points reference known functions and
	// are in valid domains
	for _, point := range v.points {
		fnLabel, err := point.fnLabel.instantiateAndMerge(vars)
		if err != nil {
			return err
		}
		domains := byNames[fnLabel]
		if fnLabel == "" { // use the abscisse axis, which has no constraints
			domains = []expression.Domain{{}}
		}
		if len(domains) == 0 {
			return fmt.Errorf("La fonction %s n'est pas définie.", fnLabel)
		}
		// check that the abscisse is in one domain
		if err := point.x.IsIncludedIntoOne(domains, vars); err != nil {
			return err
		}
	}

	return nil
}

type numberValidator struct {
	expr *expression.Expr
}

func (v numberValidator) validate(vars expression.Vars) error {
	// note that we dont allow non decimal solutions, since it is confusing for the student.
	// they should rather be handled with an expression field, or rounded using the
	// builtin round() function
	return v.expr.IsValidNumber(vars, true, true)
}

type radioValidator struct {
	proposalsLength int
	expr            *expression.Expr
}

func (v radioValidator) validate(vars expression.Vars) error {
	return v.expr.IsValidIndex(vars, v.proposalsLength)
}

type figurePointValidator struct {
	figure validator
	answer parsedCoord
}

func (v figurePointValidator) validate(vars expression.Vars) error {
	if err := v.figure.validate(vars); err != nil {
		return err
	}
	if err := v.answer.validate(vars, false); err != nil {
		return err
	}
	return nil
}

type figureVectorValidator struct {
	figure       validator
	answer       parsedCoord
	answerOrigin *parsedCoord // optional
}

func (v figureVectorValidator) validate(vars expression.Vars) error {
	if err := v.figure.validate(vars); err != nil {
		return err
	}
	if err := v.answer.validate(vars, false); err != nil {
		return err
	}
	if v.answerOrigin != nil {
		return v.answerOrigin.validate(vars, false)
	}
	return nil
}

type functionPointsValidator struct {
	xGrid    []*expression.Expr
	function function
}

// checks the x grid only contains integers values with no duplicates,
// and that the y values are integers.
func (v functionPointsValidator) validate(vars expression.Vars) error {
	seen := make(map[int]bool)

	fnExpr := expression.FunctionExpr{
		Function: v.function.Function.Copy(),
		Variable: v.function.Variable,
	}
	fnExpr.Function.Substitute(vars)
	f := fnExpr.Closure()

	// checks that all grid values are integers
	for _, xExpr := range v.xGrid {
		xValue, err := xExpr.Evaluate(vars)
		if err != nil {
			return err
		}

		val, ok := expression.IsInt(xValue)
		if !ok {
			return fmt.Errorf("L'expression %s ne définit par un antécédent <b>entier</b> (%g).", xExpr, expression.RoundFloat(xValue))
		}

		if seen[val] {
			return fmt.Errorf("Les antécédents ne sont pas uniques.")
		}
		seen[val] = true

		y := f(xValue)
		if _, ok = expression.IsInt(y); !ok {
			return fmt.Errorf("L'expression %s ne définit pas des images <b>entières</b> (%g)", fnExpr.Function, expression.RoundFloat(y))
		}
	}

	return nil
}

type figureAffineLineValidator struct {
	figure validator
	a, b   *expression.Expr
}

func (v figureAffineLineValidator) validate(vars expression.Vars) error {
	if err := v.figure.validate(vars); err != nil {
		return err
	}

	if err := v.a.IsValidNumber(vars, false, false); err != nil {
		return err
	}
	if err := v.b.IsValidNumber(vars, false, true); err != nil {
		return err
	}

	b, err := v.b.Evaluate(vars)
	if err != nil {
		return err
	}

	if _, ok := expression.IsInt(b); !ok {
		return fmt.Errorf("L'expression %s de B n'est pas un nombre entier (%f).", v.b, b)
	}

	return nil
}

// NOTE: as an optimisation, we could parse
// earlier the expression
type treeValidator struct {
	data TreeBlock
}

func (v treeValidator) validate(vars expression.Vars) error {
	var checkTree func(node TreeNodeAnswer) error
	checkTree = func(node TreeNodeAnswer) error {
		if node.Value < 0 || node.Value >= len(v.data.EventsProposals) {
			return fmt.Errorf("L'index %d n'est pas compatible avec le nombre de propositions.", node.Value)
		}

		for _, c := range node.Probabilities {
			_, err := expression.Parse(c) // we accept any valid expression to allow for instance "x"
			if err != nil {
				return err
			}
		}
		for _, c := range node.Children {
			if err := checkTree(c); err != nil {
				return err
			}
		}
		return nil
	}
	return checkTree(v.data.AnswerRoot)
}

type tableValidator struct {
	answer [][]*expression.Expr
}

func (v tableValidator) validate(vars expression.Vars) error {
	for _, row := range v.answer {
		for _, cell := range row {
			if err := cell.IsValidNumber(vars, true, true); err != nil {
				return err
			}
		}
	}
	return nil
}

type vectorValidator struct {
	answer parsedCoord
}

func (v vectorValidator) validate(vars expression.Vars) error {
	return v.answer.validate(vars, true)
}
