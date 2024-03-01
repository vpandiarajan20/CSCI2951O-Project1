package pkg

type VariableAssignment struct {
	PropagatedVariables []uint
}

type StackVA struct {
	elements []VariableAssignment
}

func (s *StackVA) PushEmpty() {
	s.elements = append(s.elements, VariableAssignment{PropagatedVariables: make([]uint, 0)})
}
func (s *StackVA) Push(VA VariableAssignment) {
	s.elements = append(s.elements, VA)
}

func (s *StackVA) Pop() (VariableAssignment, bool) {
	if len(s.elements) == 0 {
		return VariableAssignment{}, false
	}

	index := len(s.elements) - 1
	element := s.elements[index]
	s.elements = s.elements[:index]
	return element, true
}

func (s *StackVA) Peek() (VariableAssignment, bool) {
	if len(s.elements) == 0 {
		return VariableAssignment{}, false
	}

	index := len(s.elements) - 1
	return s.elements[index], true
}
