package pkg

import "fmt"

type ImplicationNode struct {
	Var      uint
	Value    int
	Level    int
	Parents  map[*ImplicationNode]bool
	Children map[*ImplicationNode]bool
	Clause   map[int]bool
}

func NewImplicationNode(Var uint, Value int) *ImplicationNode {
	return &ImplicationNode{
		Var:      Var,
		Value:    Value,
		Level:    -1,
		Parents:  make(map[*ImplicationNode]bool),
		Children: make(map[*ImplicationNode]bool),
		Clause:   make(map[int]bool),
	}
}

func NewImplicationNodeAll(Var uint, Value int, Level int, Parents map[*ImplicationNode]bool, Children map[*ImplicationNode]bool, Clause map[int]bool) *ImplicationNode {
	return &ImplicationNode{
		Var:      Var,
		Value:    Value,
		Level:    Level,
		Parents:  Parents,
		Children: Children,
		Clause:   Clause,
	}
}

func AllParents(Node *ImplicationNode) map[*ImplicationNode]bool {

	allParents := make(map[*ImplicationNode]bool)
	for k, v := range Node.Parents {
		allParents[k] = v
		for k_par, v_par := range AllParents(k) {
			allParents[k_par] = v_par
		}
	}
	return allParents
}

func (n *ImplicationNode) String() string {
	retVal := fmt.Sprintf("Var: %d, Value: %d, Level: %d, Clause: %v", n.Var, n.Value, n.Level, n.Clause)
	if len(n.Parents) > 0 {
		retVal += fmt.Sprintln()
		retVal += "Parents:"
		for parent := range n.Parents {
			retVal = retVal + fmt.Sprintf(" %d ", parent.Var)
		}
	}
	if len(n.Children) > 0 {
		retVal += fmt.Sprintln()
		retVal += "Children:"
		for child := range n.Children {
			retVal = retVal + fmt.Sprintf(" %d ", child.Var)
		}
	}
	return retVal
}
