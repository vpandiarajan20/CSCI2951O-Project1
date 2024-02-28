package pkg

import (
    "fmt"
)

type VariableAssignment struct {
    Literal int
    IsBranch  bool
}

type Stack struct {
    elements []VariableAssignment
}

func (s *Stack) Push(value int, flag bool) {
    s.elements = append(s.elements, VariableAssignment{Literal: value, IsBranch: flag})
}

func (s *Stack) Pop() (VariableAssignment, bool) {
    if len(s.elements) == 0 {
        return VariableAssignment{}, false
    }

    index := len(s.elements) - 1
    element := s.elements[index]
    s.elements = s.elements[:index]
    return element, true
}

func (s *Stack) Peek() (VariableAssignment, bool) {
    if len(s.elements) == 0 {
        return VariableAssignment{}, false
    }

    index := len(s.elements) - 1
    return s.elements[index], true
}

func main() {
    stack := Stack{}

    stack.Push(1, true)
    stack.Push(2, false)
    stack.Push(3, true)

    if element, ok := stack.Pop(); ok {
        fmt.Println("Popped:", element.Literal, element.IsBranch)
    }
    if element, ok := stack.Pop(); ok {
        fmt.Println("Popped:", element.Literal, element.IsBranch)
    }
    if element, ok := stack.Pop(); ok {
        fmt.Println("Popped:", element.Literal, element.IsBranch)
    }
    if _, ok := stack.Pop(); !ok {
        fmt.Println("Stack is empty")
    }
}
