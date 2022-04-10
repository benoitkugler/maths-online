package expression

import (
	"fmt"
	"strconv"
	"strings"
)

// implements the parsing logic for special definitions, aka intrinsic,
// of the form
// a,b,c = intrinsic(arg1, arg2)

type ErrIntrinsic struct {
	Reason string // in french
}

func (err ErrIntrinsic) Error() string {
	return err.Reason
}

// ParseIntrinsic interprets `s` as a special function definition.
// It returns `ErrIntrinsic` is the definition is invalid.
func ParseIntrinsic(s string) (Intrinsic, error) {
	chunks := strings.Split(s, "=")
	if len(chunks) != 2 {
		return nil, ErrIntrinsic{
			Reason: fmt.Sprintf("Une fonction spéciale doit contenir un seul symbol = (%d parties reçues)", len(chunks)),
		}
	}

	varNames := strings.Split(chunks[0], ",")

	startArg := strings.IndexByte(chunks[1], '(')
	endArg := strings.IndexByte(chunks[1], ')')
	if startArg == -1 || endArg == -1 || endArg < startArg {
		return nil, ErrIntrinsic{
			Reason: "Parenthèses invalides",
		}
	}

	funcName := strings.TrimSpace(chunks[1][:startArg])

	var args []string
	if argS := strings.TrimSpace(chunks[1][startArg+1 : endArg]); argS != "" {
		args = strings.Split(argS, ",")
	}

	switch funcName {
	case "pythagorians":
		return parsePythagorians(varNames, args)
	case "projection":
		return parseProjection(varNames, args)
	default:
		_ = exhaustiveIntrinsicSwitch
		return nil, ErrIntrinsic{
			Reason: fmt.Sprintf("Fonction spéciale %s inconnue", funcName),
		}
	}
}

func parseVariable(s string) Variable {
	tk := newTokenizer([]byte(strings.TrimSpace(s)))
	return tk.readVariable()
}

func parsePythagorians(variables []string, arguments []string) (out PythagorianTriplet, err error) {
	if len(variables) != 3 {
		return out, ErrIntrinsic{
			Reason: fmt.Sprintf("La fonction 'pythagorians' définit 3 variables (%d reçues)", len(variables)),
		}
	}

	switch len(arguments) {
	case 0: // bound is optionnal
		out.Bound = 10
	case 1:
		out.Bound, err = strconv.Atoi(strings.TrimSpace(arguments[0]))
		if err != nil || out.Bound < 2 {
			return out, ErrIntrinsic{
				Reason: fmt.Sprintf("L'argument optionnel de la fonction 'pythagorians' doit être un nombre entier >= 2"),
			}
		}
	default:
		return out, ErrIntrinsic{
			Reason: fmt.Sprintf("La fonction 'pythagorians' accepte un seul paramètre (optionnel) : %d reçus", len(arguments)),
		}
	}

	out.A = parseVariable(variables[0])
	out.B = parseVariable(variables[1])
	out.C = parseVariable(variables[2])

	return out, nil
}

func parseProjection(variables []string, arguments []string) (out OrthogonalProjection, err error) {
	if len(arguments) != 3 {
		return out, ErrIntrinsic{
			Reason: fmt.Sprintf("La fonction 'projection' accepte 3 points en arguments (%d reçus)", len(arguments)),
		}
	}

	switch len(variables) {
	case 1: // syntaxe for point
		indice := strings.TrimSpace(variables[0])
		out.Hx = Variable{Name: 'x', Indice: indice}
		out.Hy = Variable{Name: 'y', Indice: indice}
	case 2:
		out.Hx = parseVariable(variables[0])
		out.Hy = parseVariable(variables[1])
	default:
		return out, ErrIntrinsic{
			Reason: fmt.Sprintf("La fonction 'projection' définit 1 point (%d variables reçues)", len(variables)),
		}
	}

	A := strings.TrimSpace(arguments[0])
	B := strings.TrimSpace(arguments[1])
	C := strings.TrimSpace(arguments[2])
	out.Ax, out.Ay = Variable{Name: 'x', Indice: A}, Variable{Name: 'y', Indice: A}
	out.Bx, out.By = Variable{Name: 'x', Indice: B}, Variable{Name: 'y', Indice: B}
	out.Cx, out.Cy = Variable{Name: 'x', Indice: C}, Variable{Name: 'y', Indice: C}

	return out, nil
}
