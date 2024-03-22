package sets

type SetOp uint8

const (
	SLeaf SetOp = iota
	SUnion
	SInter
	SComplement
)

// ListNode is a tree where operator are grouped when possible,
// used by the student app.
type ListNode struct {
	// with length 0 if Op == SLeaf, one if Op == SComplement, > 1 otherwise
	Args []ListNode
	Op   SetOp
	Leaf Set // valid only if Op == SLeaf
}

func leaf(l Set) ListNode { return ListNode{Op: SLeaf, Leaf: l} }

func inter(args ...ListNode) ListNode { return ListNode{Op: SInter, Args: args} }

func union(args ...ListNode) ListNode { return ListNode{Op: SUnion, Args: args} }

func comp(arg ListNode) ListNode { return ListNode{Op: SComplement, Args: []ListNode{arg}} }

// ListSet is a tree where operator are grouped when possible,
// and leaf contents are extracted
type ListSet struct {
	Leaves []string // indexed by Leaf
	Expr   ListNode
}

func (bin BinarySet) ToList() ListSet {
	return ListSet{Leaves: bin.Sets, Expr: toList(bin.Root)}
}

func toList(node BinNode) ListNode {
	switch node := node.(type) {
	case Set:
		return ListNode{Op: SLeaf, Leaf: node}
	case Union:
		// group all same ops into one node
		left := toList(node.Left)
		right := toList(node.Right)
		var args []ListNode
		if left.Op == SUnion { // move up
			args = append(args, left.Args...)
		} else {
			args = append(args, left)
		}
		if right.Op == SUnion {
			args = append(args, right.Args...)
		} else {
			args = append(args, right)
		}
		return ListNode{Op: SUnion, Args: args}
	case Inter:
		// group all same ops into one node
		left := toList(node.Left)
		right := toList(node.Right)
		var args []ListNode
		if left.Op == SInter { // move up
			args = append(args, left.Args...)
		} else {
			args = append(args, left)
		}
		if right.Op == SInter {
			args = append(args, right.Args...)
		} else {
			args = append(args, right)
		}
		return ListNode{Op: SInter, Args: args}
	case Complement:
		right := toList(node.Right)
		return ListNode{Op: SComplement, Args: []ListNode{right}}
	default:
		panic("exhaustive switch")
	}
}

// ToBin converts the given tree to a binary version.
func (ls ListNode) ToBin() BinNode {
	switch ls.Op {
	case SLeaf:
		return ls.Leaf
	case SComplement:
		return Complement{ls.Args[0].ToBin()}
	case SUnion:
		node := ls.Args[0].ToBin()
		for _, arg := range ls.Args[1:] {
			node = Union{Left: node, Right: arg.ToBin()}
		}
		return node
	case SInter:
		node := ls.Args[0].ToBin()
		for _, arg := range ls.Args[1:] {
			node = Inter{Left: node, Right: arg.ToBin()}
		}
		return node
	default:
		panic("exhaustive switch")
	}
}
