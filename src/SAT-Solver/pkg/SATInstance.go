package pkg

import (
	"fmt"
	"sort"
)

type SATInstance struct {
	NumVars    int
	NumClauses int
	Vars       map[int]bool
	Clauses    [](map[int]bool)
}

func NewSATInstance(numVars, numClauses int) *SATInstance {
	return &SATInstance{
		numVars, numClauses, make(map[int]bool), make([]map[int]bool, 0),
	}
}

func (s *SATInstance) addVariable(literal int) {
	s.Vars[abs(literal)] = true
}

func (s *SATInstance) AddClause(clause map[int]bool) {
	s.Clauses = append(s.Clauses, clause)
}
func (s *SATInstance) String() string {
	var buf = new([]string)

	*buf = append(*buf, fmt.Sprintf("Number of variables: %d\n", s.NumVars))
	*buf = append(*buf, fmt.Sprintf("Number of clauses: %d\n", s.NumClauses))
	*buf = append(*buf, fmt.Sprintf("Variables: %v\n", sortedKeys(s.Vars)))

	for c, clause := range s.Clauses {
		*buf = append(*buf, fmt.Sprintf("Clause %d: %v\n", c, sortedKeys(clause)))
	}

	return fmt.Sprint(*buf)
}

func (s *SATInstance) PrintClauses() string {
	var buf = new([]string)
	*buf = append(*buf, "\n")
	for c, clause := range s.Clauses {
		*buf = append(*buf, fmt.Sprintf("Clause %d: %v\n", c, sortedKeys(clause)))
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

	copiedInstance := &SATInstance{
		NumVars:    instance.NumVars,
		NumClauses: instance.NumClauses,
		Vars:       newVars,
		Clauses:    newClauses,
	}

	return copiedInstance
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func sortedKeys(m map[int]bool) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}
