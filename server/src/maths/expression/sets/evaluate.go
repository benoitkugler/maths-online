package sets

const MaxNumberOfLeaves = 16

// returns an equivalent expression following the form U ∩ (S_i ou neg(S_i))
// more precisely, the returned tree is "sorted" by operators : Union > Intersection > Complement > Leaf
func normalize(node BinNode) BinNode {
	switch node := node.(type) {
	case Set: // nothing to do
		return node
	case Union:
		// simply recurse
		left := normalize(node.Left)
		right := normalize(node.Right)
		return Union{left, right}
	case Inter:
		// recurse and develop
		left := normalize(node.Left)
		right := normalize(node.Right)
		switch left := left.(type) {
		case Union: // ( ... U ...) n (...)
			switch right := right.(type) {
			case Union: // double develop
				A, B, C, D := left.Left, left.Right, right.Left, right.Right
				return Union{
					Union{
						normalize(Inter{A, C}),
						normalize(Inter{A, D}),
					},
					Union{
						normalize(Inter{B, C}),
						normalize(Inter{B, D}),
					},
				}
			default: // simple develop
				A, B, C := left.Left, left.Right, right
				return Union{normalize(Inter{A, C}), normalize(Inter{B, C})}
			}
		default:
			switch right := right.(type) {
			case Union:
				// simple develop
				A, B, C := left, right.Left, right.Right
				return Union{normalize(Inter{A, B}), normalize(Inter{A, C})}
			default:
				return Inter{left, right}
			}
		}
	case Complement: // apply de Morgan laws
		kid := node.Right
		switch kid := kid.(type) {
		case Set:
			return node // nothing to do
		case Complement: // simplify by removing the double negation
			return normalize(kid.Right)
		case Union:
			return normalize(Inter{Complement{kid.Left}, Complement{kid.Right}})
		case Inter:
			return normalize(Union{Complement{kid.Left}, Complement{kid.Right}})
		default:
			panic("exhaustive switch")
		}
	default:
		panic("exhaustive switch")
	}
}

type intersection uint16 // seen as {0, 1}^{numberbOfSets}

// sorted, unique intersections
type normalized []intersection

func (no normalized) isEqual(other normalized) bool {
	if len(no) != len(other) {
		return false
	}
	for i := range no {
		if no[i] != other[i] {
			return false
		}
	}
	return true
}

type crible []bool

func newCrible(nbSet int) crible {
	maxN := 1 << nbSet // 2^nbSet
	return make(crible, maxN)
}

// add to crible all intersection built from mask and value
func (cr crible) add(mask, value intersection) {
	for i := range cr {
		got := intersection(i)&^mask | value
		cr[got] = true
	}
}

func (cr crible) toList() (out normalized) {
	for i, b := range cr {
		if b {
			out = append(out, intersection(i))
		}
	}
	return out
}

func (set BinarySet) toNormalized() normalized {
	out := newCrible(len(set.Sets))
	expr := normalize(set.Root) // normalize
	// convert to a list tree ...
	l := toList(expr)
	// ... which has at most one level per operator, thanks to the previous step

	if l.Op != SUnion { // ensure the first level is Union
		l = union(l)
	}
	for _, inters := range l.Args {
		if inters.Op != SInter { // ensure the level is Inter
			inters = inter(inters)
		}

		var value, mask intersection
		for _, set := range inters.Args {
			// here we either have Complement or Leaf
			var (
				leaf  Set
				isNeg bool
			)
			if set.Op == SComplement {
				leaf = set.Args[0].Leaf
				isNeg = true
			} else {
				leaf = set.Leaf
			}
			mask |= 1 << leaf
			if !isNeg {
				value |= 1 << leaf
			}
		}
		out.add(mask, value)
	}

	return out.toList()
}

// IsEquivalent assumes [other] is built with the same leaves,
// and returns true if the two expressions define the same set.
// For instance, A n A == A
// and (A n B) u (A n neg(B)) == A
func (s BinarySet) IsEquivalent(other BinNode) bool {
	ref := s.toNormalized()
	otherS := BinarySet{Sets: s.Sets, Root: other}
	return ref.isEqual(otherS.toNormalized())
}
