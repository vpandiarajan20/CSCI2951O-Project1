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
	NumVars      int
	NumClauses   int
	NumConflicts int
	NumBranches  int
	Vars         map[uint]int     // vars are always positive, literals not
	Clauses      [](map[int]bool) // array of set of literals
	VarCount     map[int]struct {
		PosCount int
		NegCount int
	} // counts of positive and negative literals so that we can use them for branching
	BranchingVars    map[uint]bool
	Level            int // current level of the search
	LearnedClauses   [](map[int]bool)
	ImplicationGraph map[uint]ImplicationNode // map of variable to implication node
	BranchingHist    map[int]uint
	PropagateHist    map[int][]uint
}

type PropStruct struct {
	Literal int
	Clause  map[int]bool
}

func NewSATInstance(numVars, numClauses int) *SATInstance {
	return &SATInstance{
		NumVars:      numVars,
		NumClauses:   numClauses,
		NumConflicts: 0,
		NumBranches:  0,
		Vars:         make(map[uint]int),
		Clauses:      make([]map[int]bool, 0),
		VarCount: make(map[int]struct {
			PosCount int
			NegCount int
		}, 0),
		BranchingVars:    make(map[uint]bool),
		Level:            0,
		ImplicationGraph: make(map[uint]ImplicationNode),
		BranchingHist:    make(map[int]uint),
		PropagateHist:    make(map[int][]uint),
	}
}

func NewSATInstanceVars(numVars int) *SATInstance {
	return &SATInstance{
		NumVars:      numVars,
		NumClauses:   0,
		NumConflicts: 0,
		NumBranches:  0,
		Vars:         make(map[uint]int),
		Clauses:      make([]map[int]bool, 0),
		VarCount: make(map[int]struct {
			PosCount int
			NegCount int
		}, 0),
		BranchingVars:    make(map[uint]bool),
		Level:            0,
		ImplicationGraph: make(map[uint]ImplicationNode),
		BranchingHist:    make(map[int]uint),
		PropagateHist:    make(map[int][]uint),
	}
}

func (s *SATInstance) addVariable(literal int) {
	s.Vars[uint(abs(literal))] = Unassigned
	s.ImplicationGraph[uint(abs(literal))] = *NewImplicationNode(uint(abs(literal)), Unassigned)
}

func (s *SATInstance) AddClause(clause map[int]bool) {
	for k := range clause {
		if k < 0 {
			varStruct := s.VarCount[-k]
			varStruct.NegCount += 1
			s.VarCount[-k] = varStruct
		} else {
			varStruct := s.VarCount[k]
			varStruct.PosCount += 1
			s.VarCount[k] = varStruct
		}
	}
	s.Clauses = append(s.Clauses, clause)
	s.NumClauses += 1

}

func (s *SATInstance) RemoveClauseFromCount(clause map[int]bool) {
	for literal := range clause {
		if literal < 0 {
			varStruct := s.VarCount[-literal]
			varStruct.NegCount -= 1
			// fmt.Println("Trying to subtract from map where k is", -literal)
			// fmt.Println("Map", s.VarCount)
			// fmt.Println("Map entry", s.VarCount[-literal])
			s.VarCount[-literal] = varStruct
		} else {
			varStruct := s.VarCount[literal]
			varStruct.PosCount -= 1
			s.VarCount[literal] = varStruct
		}
	}
}

func (s *SATInstance) String() string {
	var buf = new([]string)

	*buf = append(*buf, fmt.Sprintf("Number of variables: %d\n", s.NumVars))
	*buf = append(*buf, fmt.Sprintf("Number of clauses: %d\n", s.NumClauses))
	*buf = append(*buf, fmt.Sprintf("Variables: %v\n", SortedKeysUint(s.Vars)))

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

func SortedKeysUint(m map[uint]int) []uint {
	keys := make([]uint, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}
