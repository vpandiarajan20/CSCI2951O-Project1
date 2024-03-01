package pkg

import (
	"fmt"
	"sort"
)

const (
	Unassigned = iota
	True
	False
)

type SATInstance struct {
	NumVars    int
	NumClauses int
	Vars       map[uint]int
	Clauses    [](map[int]bool)
	VarCount   (map[uint]struct {
		PosCount int
		NegCount int
	})
	UnsatisfiedClauses map[int]bool
	AssignmentStack    StackVA
	ClauseStack        StackVA
}

func NewSATInstance(numVars, numClauses int) *SATInstance {
	return &SATInstance{
		NumVars:    numVars,
		NumClauses: numClauses,
		Vars:       make(map[uint]int),
		Clauses:    make([]map[int]bool, 0),
		VarCount: make(map[uint]struct {
			PosCount int
			NegCount int
		}, 0),
	}
}
func NewSATInstanceVars(numVars int) *SATInstance {
	return &SATInstance{
		NumVars:    numVars,
		NumClauses: 0,
		Vars:       make(map[uint]int),
		Clauses:    make([]map[int]bool, 0),
		VarCount: make(map[uint]struct {
			PosCount int
			NegCount int
		}, 0),
	}
}

func (s *SATInstance) addVariable(literal int) {
	s.Vars[uint(abs(literal))] = Unassigned
}

func (s *SATInstance) AddClause(clause map[int]bool) {
	for k := range clause {
		varStruct := s.VarCount[uint(abs(k))]
		if k < 0 {
			varStruct.NegCount += 1
		} else {
			varStruct := s.VarCount[uint(k)]
			varStruct.PosCount += 1
		}
		s.VarCount[uint(abs(k))] = varStruct
	}
	s.Clauses = append(s.Clauses, clause)
	s.UnsatisfiedClauses = append(s.UnsatisfiedClauses, len(s.Clauses)-1)
	s.NumClauses += 1

}

func (s *SATInstance) RemoveClauseFromCount(clause map[int]bool) {
	for literal := range clause {
		varStruct := s.VarCount[uint(abs(literal))]
		if literal < 0 {
			varStruct.NegCount -= 1
			// fmt.Println("Trying to subtract from map where k is", -literal)
			// fmt.Println("Map", s.VarCount)
			// fmt.Println("Map entry", s.VarCount[-literal])
		} else {
			varStruct.PosCount -= 1
		}
		s.VarCount[uint(abs(literal))] = varStruct
	}
}

func (s *SATInstance) RemoveLiteralFromCount(literal int) {
	varStruct := s.VarCount[uint(abs(literal))]
	if literal < 0 {
		varStruct.NegCount -= 1
	} else {
		varStruct.PosCount -= 1
	}
	s.VarCount[uint(abs(literal))] = varStruct
}

func (s *SATInstance) String() string {
	var buf = new([]string)

	*buf = append(*buf, fmt.Sprintf("Number of variables: %d\n", s.NumVars))
	*buf = append(*buf, fmt.Sprintf("Number of clauses: %d\n", s.NumClauses))
	*buf = append(*buf, fmt.Sprintf("Variables: %v\n", SortedKeysUint(s.Vars)))

	for c, clause := range s.Clauses {
		*buf = append(*buf, fmt.Sprintf("Clause %d: %v\n", c, SortedKeysInt(clause)))
	}

	*buf = append(*buf, "VarCount:\n")
	for varID, counts := range s.VarCount {
		*buf = append(*buf, fmt.Sprintf("  Variable %d: PosCount=%d, NegCount=%d\n", varID, counts.PosCount, counts.NegCount))
	}

	return fmt.Sprint(*buf)
}

func (s *SATInstance) PrintClauses() string {
	var buf = new([]string)
	*buf = append(*buf, "\n")
	for c, clause := range s.Clauses {
		*buf = append(*buf, fmt.Sprintf("Clause %d: %v\n", c, SortedKeysInt(clause)))
	}
	return fmt.Sprint(*buf)
}

func DeepCopySATInstance(instance SATInstance) *SATInstance {
	newVars := make(map[uint]int, len(instance.Vars))
	for k, v := range instance.Vars {
		newVars[k] = v
	}

	newClauses := make([]map[int]bool, len(instance.Clauses))
	for i, clause := range instance.Clauses {
		newClauses[i] = make(map[int]bool, len(clause))
		for k, v := range clause {
			newClauses[i][k] = v
		}
	}

	newVarCount := make(map[uint]struct {
		PosCount int
		NegCount int
	})
	for k, v := range instance.VarCount {
		newVarCount[k] = v
	}

	// Create a copied instance with deep-copied fields
	copiedInstance := &SATInstance{
		NumVars:    instance.NumVars,
		NumClauses: instance.NumClauses,
		Vars:       newVars,
		Clauses:    newClauses,
		VarCount:   newVarCount,
	}

	return copiedInstance
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func SortedKeysUint(m map[uint]int) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	return keys
}

func SortedKeysInt(m map[int]bool) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}
