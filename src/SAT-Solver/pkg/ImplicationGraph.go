package pkg

import "fmt"

type ImplicationNode struct {
	Var     uint
	Value   int
	Level   int
	Parents map[uint]bool
	Clause  map[int]bool
	// TODO: parents and children should probaly be a list?
}

func NewImplicationNode(Var uint, Value int) *ImplicationNode {
	return &ImplicationNode{
		Var:      Var,
		Value:    Value,
		Level:    -1,
		Parents:  make(map[uint]bool),
		Children: make(map[uint]bool),
		Clause:   make(map[int]bool),
	}
}

func NewImplicationNodeAll(Var uint, Value int, Level int, Parents, Children map[uint]bool, Clause map[int]bool) *ImplicationNode {
	return &ImplicationNode{
		Var:      Var,
		Value:    Value,
		Level:    Level,
		Parents:  Parents,
		Children: Children,
		Clause:   Clause,
	}
}

// func AllParents(Node *ImplicationNode, s *SATInstance) map[*ImplicationNode]bool {

// 	for k := range Node.Parents {
// 		parent := s.ImplicationGraph[k]
// 		for k_par, v_par := range AllParents(&parent, s) {
// 			allParents[k_par] = v_par
// 		}
// 	}
// 	return allParents
// }

func (n *ImplicationNode) String() string {
	retVal := fmt.Sprintf("Var: %d, Value: %d, Level: %d, Parents: %v, Children: %v, Clause: %v", n.Var, n.Value, n.Level, n.Parents, n.Children, n.Clause)
	// retVal := fmt.Sprintf("Var: %d, Value: %d, Level: %d, Clause: %v", n.Var, n.Value, n.Level, n.Clause)
	// if len(n.Parents) > 0 {
	// 	retVal += fmt.Sprintln()
	// 	retVal += "Parents:"
	// 	for parent := range n.Parents {
	// 		retVal = retVal + fmt.Sprintf(" %d ", parent)
	// 	}
	// }
	// if len(n.Children) > 0 {
	// 	retVal += fmt.Sprintln()
	// 	retVal += "Children:"
	// 	for child := range n.Children {
	// 		retVal = retVal + fmt.Sprintf(" %d ", child)
	// 	}
	// }
	return retVal
}
