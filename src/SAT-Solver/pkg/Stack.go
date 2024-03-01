package pkg

type VariableAssignment struct {
	Literal  int
	IsBranch bool
}

type StackVA struct {
	elements []VariableAssignment
}

func (s *StackVA) Push(value int, flag bool) {
	s.elements = append(s.elements, VariableAssignment{Literal: value, IsBranch: flag})
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

type StackUInt struct {
	elements []uint
}

func NewStackUint(elements []uint) StackUInt {
	return StackUInt{elements: elements}
}

func (s *StackUInt) Push(value uint) {
	s.elements = append(s.elements, value)
}

func (s *StackUInt) Pop() (uint, bool) {
	if len(s.elements) == 0 {
		return 0, false
	}

	index := len(s.elements) - 1
	element := s.elements[index]
	s.elements = s.elements[:index]
	return element, true
}
