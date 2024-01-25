package sets

const MaxNumberOfLeaves = 16

type intersection uint16 // seen as {0, 1}^{numberbOfSets}

type normalized []intersection // sorted

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
