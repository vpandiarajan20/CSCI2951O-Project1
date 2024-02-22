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
	Vars       map[int]int
	StackAssignments Stack
	Clauses    [](map[int]bool)
	WatchedLiterals   (map[int]struct {
		Literal1 int // watched literal 1
		Literal2 int // watched literal 2
	})
}

func NewSATInstance(numVars, numClauses int) *SATInstance {
	return &SATInstance{
		NumVars:    numVars,
		NumClauses: numClauses,
		Vars:       make(map[int]bool),
		Clauses:    make([]map[int]bool, 0),
		VarCount: make(map[int]struct {
			PosCount int
			NegCount int
		}, 0),
	}
}
func NewSATInstanceVars(numVars int) *SATInstance {
	return &SATInstance{
		NumVars:    numVars,
		NumClauses: 0,
		Vars:       make(map[int]bool),
		Clauses:    make([]map[int]bool, 0),
		VarCount: make(map[int]struct {
			PosCount int
			NegCount int
		}, 0),
	}
}

func (s *SATInstance) addVariable(literal int) {
	s.Vars[abs(literal)] = true
}

func (s *SATInstance) AddClause(clause map[int]bool) {
	// function adds clause to the SATInstance
	s.Clauses = append(s.Clauses, clause)
	s.NumClauses += 1

}

func (s *SATInstance) String() string {
	var buf = new([]string)

	*buf = append(*buf, fmt.Sprintf("Number of variables: %d\n", s.NumVars))
	*buf = append(*buf, fmt.Sprintf("Number of clauses: %d\n", s.NumClauses))
	*buf = append(*buf, fmt.Sprintf("Variables: %v\n", SortedKeys(s.Vars)))

	for c, clause := range s.Clauses {
		*buf = append(*buf, fmt.Sprintf("Clause %d: %v\n", c, SortedKeys(clause)))
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
		*buf = append(*buf, fmt.Sprintf("Clause %d: %v\n", c, SortedKeys(clause)))
	}
	return fmt.Sprint(*buf)
}

func DeepCopySATInstance(instance SATInstance) *SATInstance {
	newVars := make(map[int]bool, len(instance.Vars))
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

	newVarCount := make(map[int]struct {
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

func SortedKeys(m map[int]bool) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}
