package expression

import (
	"reflect"
	"strconv"
	"testing"
)

func Test_matricesOperations(t *testing.T) {
	m := RandomParameters{defs: map[Variable]*Expr{
		NewVar('A'): mustParse(t, "[[1;2]; [3;4]]"),
		NewVar('B'): mustParse(t, "[[1;2]; [4;5]]"),
		NewVar('C'): mustParse(t, "[[1;0]; [1;1]; [-1;1]]"),
		NewVar('D'): mustParse(t, "[[2;0]; [0;5]]"),
		NewVar('E'): mustParse(t, "[[2;x]; [0;5]]"),
	}}
	ops := []struct {
		expr string
		want matrix
	}{
		{"A + B", matrix{{newNb(2), newNb(4)}, {newNb(7), newNb(9)}}},
		{"A - B", matrix{{newNb(0), newNb(0)}, {newNb(-1), newNb(-1)}}},
		{"A * B", matrix{{newNb(9), newNb(12)}, {newNb(19), newNb(26)}}},
		{"B * A", matrix{{newNb(7), newNb(10)}, {newNb(19), newNb(28)}}},
		{"C * A", matrix{{newNb(1), newNb(2)}, {newNb(4), newNb(6)}, {newNb(2), newNb(2)}}},
		{"2 * A", matrix{{newNb(2), newNb(4)}, {newNb(6), newNb(8)}}},
		{"A * 2", matrix{{newNb(2), newNb(4)}, {newNb(6), newNb(8)}}},
		{"A ^ 0", matrix{{newNb(1), newNb(0)}, {newNb(0), newNb(1)}}},
		{"A ^ 1", matrix{{newNb(1), newNb(2)}, {newNb(3), newNb(4)}}},
		{"A ^ 2", matrix{{newNb(7), newNb(10)}, {newNb(15), newNb(22)}}},
		{"A ^ 3", matrix{{newNb(37), newNb(54)}, {newNb(81), newNb(118)}}},
		{"A ^ 10", matrix{{newNb(4783807), newNb(6972050)}, {newNb(10458075), newNb(15241882)}}},
		{"-A", matrix{{newNb(-1), newNb(-2)}, {newNb(-3), newNb(-4)}}},
		{"trans(A)", matrix{{newNb(1), newNb(3)}, {newNb(2), newNb(4)}}},
		{"transpose(A)", matrix{{newNb(1), newNb(3)}, {newNb(2), newNb(4)}}},
		{"inv(D)", matrix{{newNb(0.5), newNb(0)}, {newNb(0), newNb(0.2)}}},
	}

	for i, op := range ops {
		m.defs[NewVarI('o', strconv.Itoa(i))] = mustParse(t, op.expr)
	}
	vars, err := m.Instantiate()
	if err != nil {
		t.Fatal(err)
	}
	for i, op := range ops {
		got := vars[NewVarI('o', strconv.Itoa(i))]
		want := &Expr{atom: op.want}
		want.DefaultSimplify()
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("for %s expected %s, got %s", op.expr, op.want, got)
		}
	}
}

func BenchmarkMatrixProduct(b *testing.B) {
	A, _ := newNumberMatrixFrom(mustParse(b, `[[1 ; 2; 3]; 
						[4 ; 5; 6];
						[7 ; 8; 9]]`).atom.(matrix))

	for i := 0; i < b.N; i++ {
		_, _ = A.pow(10)
	}
}

func Test_matricesOperations_invalid(t *testing.T) {
	ops := []string{
		"A + B",
		"A - B",
		"A * B",
		"B ^ 2",
		"B ^ (-1)",
		"1 + A",
		"A + 1",
		"1 - A",
		"A - 1",
		"trace(B)",
		"inv(B)",
	}

	for _, op := range ops {
		m := RandomParameters{defs: map[Variable]*Expr{
			NewVar('A'): mustParse(t, "[[1;2]; [3;4]]"),
			NewVar('B'): mustParse(t, "[[1;2]; [4;5]; [6;7]]"),
			NewVar('C'): mustParse(t, op),
		}}
		_, err := m.Instantiate()
		if err == nil {
			t.Fatal("expected error on invalid matrix operation")
		}
	}
}

func Test_matrix_submatrix(t *testing.T) {
	tests := []struct {
		A    string
		i, j int
		want string
	}{
		{`[[1 ; 2]; 
		   [3 ; 4]]`, 0, 0, "[[4]]"},
		{`[[1 ; 2]; 
		   [3 ; 4]]`, 1, 0, "[[2]]"},
		{`[[1 ; 2]; 
		   [3 ; 4]]`, 0, 1, "[[3]]"},
		{`[[1 ; 2]; 
		   [3 ; 4]]`, 1, 1, "[[1]]"},
		{`[[1 ; 2; 3]; 
		   [4 ; 5; 6];
		   [7 ; 8; 9]]`, 2, 1, "[[1; 3];[4;6]]"},
	}
	for _, tt := range tests {
		A, _ := newNumberMatrixFrom(mustParse(t, tt.A).atom.(matrix))
		want, _ := newNumberMatrixFrom(mustParse(t, tt.want).atom.(matrix))
		if got := A.submatrix(tt.i, tt.j); !reflect.DeepEqual(got, want) {
			t.Errorf("matrix.submatrix() = %v, want %v", got, tt.want)
		}
	}
}

func Test_matrix_determinant(t *testing.T) {
	tests := []struct {
		A       string
		want    float64
		wantErr bool
	}{
		{`[[1 ; 2]]`, 0, true},
		{`[[1 ; 2]; 
		   [3 ; 4]]`, 1*4 - 3*2, false},
		{`[[1 ; 0]; 
		   [0 ; 4]]`, 4, false},
		{`[[1]]`, 1, false},
		{`[[1 ; 2; 3]; 
		   [0 ; 5; 6];
		   [0 ; 0; 9]]`, 9 * 5, false},
		{`[[1 ; 2; 3]; 
		   [0 ; 0; 0];
		   [0 ; 0; 9]]`, 0, false},
		{`[[1 ; 2; 3]; 
		   [0 ; 5; 6];
		   [1 ; 2; 1]]`, (2*6 - 5*3) + 2*(-6) + 5, false},
	}
	for _, tt := range tests {
		A, _ := newNumberMatrixFrom(mustParse(t, tt.A).atom.(matrix))
		got, err := A.determinant()
		if (err != nil) != tt.wantErr {
			t.Errorf("matrix.determinant() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if got != tt.want {
			t.Errorf("matrix.determinant() = %v, want %v", got, tt.want)
		}
	}
}

func Test_matrix_invert(t *testing.T) {
	tests := []struct {
		A       string
		want    string
		wantErr bool
	}{
		{"[[1; 2]]", "[[]]", true},
		{
			`[[1 ; 0]; 
		   	  [0 ; 1]]`,
			`[[1 ; 0]; 
		   	  [0 ; 1]]`, false,
		},
		{
			`[[1 ; 0]; 
		   	  [0 ; 2]]`,
			`[[1 ; 0]; 
		   	  [0 ; 1/2]]`, false,
		},
	}
	for _, tt := range tests {
		A, _ := newNumberMatrixFrom(mustParse(t, tt.A).atom.(matrix))
		got, err := A.invert()
		if (err != nil) != tt.wantErr {
			t.Errorf("matrix.invert() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		want, _ := newNumberMatrixFrom(mustParse(t, tt.want).atom.(matrix))
		if !tt.wantErr && !reflect.DeepEqual(got, want) {
			t.Errorf("matrix.invert() = %v, want %v", got, want)
		}
	}
}
