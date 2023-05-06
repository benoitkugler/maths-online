package expression

import "fmt"

// Implements matrices operations (multiplication, trace, determinant, etc..)

func add(a, b *Expr) *Expr  { return &Expr{atom: plus, left: a, right: b} }
func prod(a, b *Expr) *Expr { return &Expr{atom: mult, left: a, right: b} }

func (A matrix) dims() (rows, cols int) {
	if len(A) == 0 {
		return 0, 0
	}
	return len(A), len(A[0])
}

// factor . A
func (A matrix) scale(factor *Expr) matrix {
	mA, nA := A.dims()
	out := make(matrix, mA)
	for i := range out {
		row := make([]*Expr, nA)
		for j := range row {
			Aij := A[i][j]
			row[j] = prod(factor, Aij)
		}
		out[i] = row
	}
	return out
}

// A + B
func (A matrix) plus(B matrix) (matrix, error) {
	mA, nA := A.dims()
	mB, nB := B.dims()
	if mA != mB || nA != nB {
		return nil, fmt.Errorf("matrices de dimensions incompatibles (%d, %d != %d, %d)", mA, nA, mB, nB)
	}

	out := make(matrix, mA)
	for i := range out {
		row := make([]*Expr, nA)
		for j := range row {
			Aij := A[i][j]
			Bij := B[i][j]
			row[j] = add(Aij, Bij)
		}
		out[i] = row
	}

	return out, nil
}

// A - B
func (A matrix) minus(B matrix) (matrix, error) {
	B = B.scale(newNb(-1))
	return A.plus(B)
}

// A * B
func (A matrix) prod(B matrix) (matrix, error) {
	mA, nA := A.dims()
	mB, nB := B.dims()
	if nA != mB {
		return nil, fmt.Errorf("matrices de dimensions incompatibles (%d != %d)", nA, mB)
	}

	out := make(matrix, mA)
	for i := range out {
		row := make([]*Expr, nB)
		for j := range row {
			s := newNb(0)
			for k := range A[i] {
				Aik := A[i][k]
				Bkj := B[k][j]
				e := prod(Aik, Bkj)
				s = add(s, e)
			}
			row[j] = s
		}
		out[i] = row
	}
	return out, nil
}

// fow now, only handle positive values
func (A matrix) pow(nf float64) (matrix, error) {
	n, ok := IsInt(nf)
	if !ok || n < 0 {
		return nil, fmt.Errorf("puissance de matrice %f non entière positive", nf)
	}

	m, n := A.dims()
	if m != n {
		return nil, fmt.Errorf("puissance d'une matrice non carrée (%d, %d)", m, n)
	}

	// identity
	out := make(matrix, m)
	for i := range out {
		row := make([]*Expr, n)
		for j := range row {
			if i == j {
				row[j] = newNb(1)
			} else {
				row[j] = newNb(0)
			}
		}
		out[i] = row
	}

	for k := 0; k < n; k++ {
		out, _ = out.prod(A) // prod is safe here
	}

	return out, nil
}

func (A matrix) trace() (*Expr, error) {
	m, n := A.dims()
	if m != n {
		return nil, fmt.Errorf("trace d'une matrice non carrée (%d, %d)", m, n)
	}

	s := newNb(0)
	for k := range A {
		Akk := A[k][k]
		s = add(s, Akk)
	}
	return s, nil
}

func (A matrix) transpose() matrix {
	m, n := A.dims()

	out := make(matrix, n)
	for i := range out {
		row := make([]*Expr, m) // out
		for j := range row {
			row[j] = A[j][i]
		}
		out[i] = row
	}

	return out
}

func (A matrix) determinant() (*Expr, error) {
	m, n := A.dims()
	if m != n {
		return nil, fmt.Errorf("déterminant d'une matrice non carrée (%d, %d)", m, n)
	}

	if m == 0 {
		return newNb(0), nil
	} else if m == 1 {
		return A[0][0], nil
	} else if m == 2 {
		return &Expr{atom: minus, left: prod(A[0][0], A[1][1]), right: prod(A[1][0], A[0][1])}, nil
	} else { // last row dev
		det := newNb(0)
		for j := range A {
			minor := A.minor(m-1, j)
			coeff := A[m-1][j]
			det = add(det, prod(coeff, minor))
		}
		return det, nil
	}
}

// assume A is squared
func (A matrix) minor(i, j int) *Expr {
	minor, _ := A.submatrix(i, j).determinant() // determinant is safe here
	if (i+j)%2 != 0 {                           // i and j is shifted by one, compensating
		minor = &Expr{atom: minus, left: nil, right: minor}
	}
	return minor
}

// return 1/detA * transpose(comatrice(A))
func (A matrix) invert() (matrix, error) {
	det, err := A.determinant()
	if err != nil {
		return nil, err
	}
	comTransp := make(matrix, len(A))
	for i := range A {
		row := make([]*Expr, len(A))
		for j := range row {
			row[j] = A.minor(j, i)
		}
		comTransp[i] = row
	}
	return comTransp.scale(&Expr{atom: div, left: newNb(1), right: det}), nil
}

// remove the row i and column j
// assume A is squared, with dims >= 1
// the coefficient expressions are not deep copied
func (A matrix) submatrix(i, j int) matrix {
	out := make(matrix, 0, len(A)-1)
	for i2, row := range A {
		if i2 == i { // remove row i
			continue
		}
		newRow := make([]*Expr, len(A)-1)
		copy(newRow, row[:j])
		copy(newRow[j:], row[j+1:])
		out = append(out, newRow)
	}
	return out
}

func (A matrix) instantiate(ctx *paramsInstantiater) (matrix, error) {
	out := make(matrix, len(A))
	var err error
	for i, row := range A {
		cols := make([]*Expr, len(row))
		for j, col := range row {
			cols[j], err = col.instantiate(ctx)
			if err != nil {
				return nil, err
			}
		}
		out[i] = cols
	}
	return out, nil
}

func matricesOperation(op operator, left, right *Expr, ctx *paramsInstantiater) (*Expr, error) {
	switch op {
	case plus:
		if leftA, ok := left.atom.(matrix); ok {
			if rightA, ok := right.atom.(matrix); ok {
				m, err := leftA.plus(rightA)
				if err != nil {
					return nil, err
				}
				m, err = m.instantiate(ctx)
				if err != nil {
					return nil, err
				}
				return &Expr{atom: m}, nil
			} else {
				return nil, fmt.Errorf("somme entre une matrice et un scalaire")
			}
		} else if _, ok := right.atom.(matrix); ok {
			return nil, fmt.Errorf("somme entre un scalaire et une matrice")
		}
	case minus:
		if rightA, ok := right.atom.(matrix); ok {
			if left == nil { // - A
				rightA = rightA.scale(newNb(-1))
				rightA, err := rightA.instantiate(ctx)
				if err != nil {
					return nil, err
				}
				return &Expr{atom: rightA}, nil
			}

			if leftA, ok := left.atom.(matrix); ok {
				m, err := leftA.minus(rightA)
				if err != nil {
					return nil, err
				}
				m, err = m.instantiate(ctx)
				if err != nil {
					return nil, err
				}
				return &Expr{atom: m}, nil
			} else {
				return nil, fmt.Errorf("différence entre une matrice et un scalaire")
			}
		} else if left != nil {
			if _, ok := left.atom.(matrix); ok {
				return nil, fmt.Errorf("différence entre un scalaire et une matrice")
			}
		}
	case mult:
		if leftA, ok := left.atom.(matrix); ok {
			if rightA, ok := right.atom.(matrix); ok {
				m, err := leftA.prod(rightA)
				if err != nil {
					return nil, err
				}
				m, err = m.instantiate(ctx)
				if err != nil {
					return nil, err
				}
				return &Expr{atom: m}, nil
			} else {
				m := leftA.scale(right)
				m, err := m.instantiate(ctx)
				if err != nil {
					return nil, err
				}
				return &Expr{atom: m}, nil
			}
		} else if rightA, ok := right.atom.(matrix); ok {
			m := rightA.scale(left)
			m, err := m.instantiate(ctx)
			if err != nil {
				return nil, err
			}
			return &Expr{atom: m}, nil
		}
	case pow:
		if leftA, ok := left.atom.(matrix); ok {
			// only expand fixed int numbers
			if val, ok := right.isConstantTerm(); ok {
				m, err := leftA.pow(val)
				if err != nil {
					return nil, err
				}
				m, err = m.instantiate(ctx)
				if err != nil {
					return nil, err
				}
				return &Expr{atom: m}, nil
			}
		}
	case equals, greater, strictlyGreater, lesser, strictlyLesser, div, mod, rem, factorial:
		// pass
	default:
		panic(exhaustiveOperatorSwitch)
	}
	return &Expr{atom: op, left: left, right: right}, nil
}
