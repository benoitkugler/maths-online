package expression

import (
	"errors"
	"fmt"
)

// Implements matrices operations (multiplication, trace, determinant, etc..)

func (A matrix) dims() (rows, cols int) {
	if len(A) == 0 {
		return 0, 0
	}
	return len(A), len(A[0])
}

func newMatrixEmpty(n, m int) matrix {
	out := make(matrix, n)
	vals := make([]*Expr, n*m)
	for i := range out {
		out[i] = vals[i*m : (i+1)*m]
	}
	return out
}

func add(a, b *Expr) *Expr { return &Expr{atom: plus, left: a, right: b} }

func prod(a, b *Expr) *Expr { return &Expr{atom: mult, left: a, right: b} }

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

// numberMatrix is an efficient version of a matrix
// containing only scalars
type numberMatrix [][]float64

func (A numberMatrix) dims() (rows, cols int) {
	if len(A) == 0 {
		return 0, 0
	}
	return len(A), len(A[0])
}

func newZeros(n, m int) numberMatrix {
	out := make(numberMatrix, n)
	vals := make([]float64, n*m)
	for i := range out {
		out[i] = vals[i*m : (i+1)*m]
	}
	return out
}

func newNumberMatrixFrom(A matrix) (numberMatrix, bool) {
	n, m := A.dims()
	out := newZeros(n, m)
	for i, row := range A {
		for j, v := range row {
			var ok bool
			out[i][j], ok = v.isConstantTerm()
			if !ok {
				return nil, false
			}
		}
	}
	return out, true
}

func (A numberMatrix) toExprMatrix() matrix {
	n, m := A.dims()
	out := newMatrixEmpty(n, m)
	for i, row := range A {
		for j, v := range row {
			out[i][j] = NewNb(v)
		}
	}
	return out
}

// assume A and B are compatible and [out] has the correct size
func (A numberMatrix) prodTo(B, out numberMatrix) {
	for i, row := range out {
		for j := range row {
			s := 0.
			for k := range A[i] {
				Aik := A[i][k]
				Bkj := B[k][j]
				s += Aik * Bkj
			}
			out[i][j] = s
		}
	}
}

// return A * B
func (A numberMatrix) prod(B numberMatrix) (numberMatrix, error) {
	mA, nA := A.dims()
	mB, nB := B.dims()
	if nA != mB {
		return nil, fmt.Errorf("matrices de dimensions incompatibles (%d != %d)", nA, mB)
	}

	out := newZeros(mA, nB)
	A.prodTo(B, out)

	return out, nil
}

// fow now, only handle positive values
func (A numberMatrix) pow(nf float64) (numberMatrix, error) {
	kmax, ok := IsInt(nf)
	if !ok || kmax < 0 {
		return nil, fmt.Errorf("puissance de matrice %f non entière positive", nf)
	}

	m, n := A.dims()
	if m != n {
		return nil, fmt.Errorf("puissance d'une matrice non carrée (%d, %d)", m, n)
	}

	// identity
	out := newZeros(n, n)
	for i := range out {
		for j := range out {
			if i == j {
				out[i][j] = 1
			}
		}
	}

	if kmax == 0 {
		return out, nil
	} else if kmax == 1 {
		return A, nil
	} else if kmax == 2 {
		out, _ = A.prod(A) // prod is safe here
		return out, nil
	} else {
		// only allocate two buffers
		tmp := newZeros(n, n)
		for k := 0; k < kmax; k++ {
			out.prodTo(A, tmp)
			out, tmp = tmp, out // switch buffers
		}
	}

	return out, nil
}

func (A numberMatrix) determinant() (float64, error) {
	m, n := A.dims()
	if m != n {
		return 0, fmt.Errorf("déterminant d'une matrice non carrée (%d, %d)", m, n)
	}

	if m == 0 {
		return 0, nil
	} else if m == 1 {
		return A[0][0], nil
	} else if m == 2 {
		return A[0][0]*A[1][1] - A[1][0]*A[0][1], nil
	} else { // last row dev
		det := 0.
		for j := range A {
			minor := A.minor(m-1, j)
			coeff := A[m-1][j]
			det += coeff * minor
		}
		return det, nil
	}
}

// assume A is squared
func (A numberMatrix) minor(i, j int) float64 {
	minor, _ := A.submatrix(i, j).determinant() // determinant is safe here
	if (i+j)%2 != 0 {                           // i and j is shifted by one, compensating
		minor = -minor
	}
	return minor
}

// remove the row i and column j
// assume A is squared, with dims >= 1
func (A numberMatrix) submatrix(i, j int) numberMatrix {
	out := make(numberMatrix, 0, len(A)-1)
	for i2, row := range A {
		if i2 == i { // remove row i
			continue
		}
		newRow := make([]float64, len(A)-1)
		copy(newRow, row[:j])
		copy(newRow[j:], row[j+1:])
		out = append(out, newRow)
	}
	return out
}

// return 1/detA * transpose(comatrice(A))
// or an error if detA == 0
func (A numberMatrix) invert() (numberMatrix, error) {
	det, err := A.determinant()
	if err != nil {
		return nil, err
	}
	if det == 0 {
		return nil, errors.New("matrice non inversible (déterminant nul)")
	}
	inv := newZeros(len(A), len(A))
	for i, row := range inv {
		for j := range row {
			row[j] = A.minor(j, i) / det // tranpose
		}
	}
	return inv, nil
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

func (A matrix) instantiate(ctx *resolver) (matrix, error) {
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

func matricesOperation(op operator, left, right *Expr, ctx *resolver) (*Expr, error) {
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
				leftN, ok1 := newNumberMatrixFrom(leftA)
				rightN, ok2 := newNumberMatrixFrom(rightA)
				if !(ok1 && ok2) {
					return nil, errors.New("Le produit de deux matrices ne supporte pas le calcul symbolique")
				}
				m, err := leftN.prod(rightN)
				if err != nil {
					return nil, err
				}
				return &Expr{atom: m.toExprMatrix()}, nil
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
				leftN, ok := newNumberMatrixFrom(leftA)
				if !ok {
					return nil, errors.New("La puissance d'une matrice ne supporte pas le calcul symbolique.")
				}
				m, err := leftN.pow(val)
				if err != nil {
					return nil, err
				}
				return &Expr{atom: m.toExprMatrix()}, nil
			}
		}
	case equals, greater, strictlyGreater, lesser, strictlyLesser, div, mod, rem, factorial:
		// pass
	default:
		panic(exhaustiveOperatorSwitch)
	}
	return &Expr{atom: op, left: left, right: right}, nil
}
