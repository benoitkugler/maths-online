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

// ParseIntrinsic interprets `s` as a special function definition,
// adding it to the parameters.
// It returns `ErrIntrinsic` is the definition is invalid.
func (rd *RandomParameters) ParseIntrinsic(s string) error {
	chunks := strings.Split(s, "=")
	if len(chunks) != 2 {
		return ErrIntrinsic{
			Reason: fmt.Sprintf("Une fonction spéciale doit contenir un seul symbol = (%d parties reçues)", len(chunks)),
		}
	}

	varNames := strings.Split(chunks[0], ",")

	startArg := strings.IndexByte(chunks[1], '(')
	endArg := strings.IndexByte(chunks[1], ')')
	if startArg == -1 || endArg == -1 || endArg < startArg {
		return ErrIntrinsic{
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
		p, err := parsePythagorians(varNames, args)
		if err != nil {
			return err
		}
		rd.specials = append(rd.specials, p)
		return nil
	case "projection":
		p, err := parseProjection(varNames, args)
		if err != nil {
			return err
		}
		return p.mergeTo(rd)
	default:
		_ = exhaustiveIntrinsicSwitch
		return ErrIntrinsic{
			Reason: fmt.Sprintf("Fonction spéciale %s inconnue", funcName),
		}
	}
}

func parseVariable(s string) Variable {
	tk := newTokenizer([]byte(strings.TrimSpace(s)))
	return tk.readVariable()
}

func parsePythagorians(variables []string, arguments []string) (out pythagorianTriplet, err error) {
	if len(variables) != 3 {
		return out, ErrIntrinsic{
			Reason: fmt.Sprintf("La fonction 'pythagorians' définit 3 variables (%d reçues)", len(variables)),
		}
	}

	switch len(arguments) {
	case 0: // bound is optionnal
		out.bound = 10
	case 1:
		out.bound, err = strconv.Atoi(strings.TrimSpace(arguments[0]))
		if err != nil || out.bound < 2 {
			return out, ErrIntrinsic{
				Reason: "L'argument optionnel de la fonction 'pythagorians' doit être un nombre entier >= 2",
			}
		}
	default:
		return out, ErrIntrinsic{
			Reason: fmt.Sprintf("La fonction 'pythagorians' accepte un seul paramètre (optionnel) : %d reçus", len(arguments)),
		}
	}

	out.a = parseVariable(variables[0])
	out.b = parseVariable(variables[1])
	out.c = parseVariable(variables[2])

	return out, nil
}

func parseProjection(variables []string, arguments []string) (out orthogonalProjection, err error) {
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
