package pkg

import (
	"fmt"
	"log"
)

func DPLL(f *SATInstance) (*SATInstance, bool) {
	if len(f.Clauses) == 0 {
		fmt.Println("0 clauses, returning true")
		fmt.Println("truth assigns", f.Vars)
		return f, true
	}

	fPrime := DeepCopySATInstance(*f)

	// fmt.Println("Pre-UnitProp", fPrime)
	UnitPropagate(fPrime)
	// fmt.Println("Post-UnitProp, Pre-Literal", fPrime)
	PureLiteralElim(fPrime)
	// fmt.Println("Post-Literal, Pre-Split", fPrime)

	// checking for empty clause unsat
	for _, clause := range fPrime.Clauses {
		if len(clause) == 0 {
			return nil, false
		}
	}

	if len(fPrime.Clauses) == 0 {
		// fmt.Println("0 clauses, returning true", fPrime)
		fmt.Println("truth assigns", f.Vars)
		return fPrime, true
	}

	literal := SplittingRule(fPrime)
	// fmt.Println("Split on:", literal)

	fRightPrime := DeepCopySATInstance(*fPrime)

	newClause := make(map[int]bool, 0)
	newClause[literal] = false
	fPrime.AddClause(newClause)
	retSAT, isSAT := DPLL(fPrime)
	if isSAT {
		// fmt.Println("0 clauses, returning true")
		fmt.Println("truth assigns", retSAT.Vars)
		return retSAT, isSAT
	}
	// fmt.Println("Split on:", literal, "left failed")

	newClause = make(map[int]bool, 0)
	newClause[-literal] = false
	fRightPrime.AddClause(newClause)
	retSAT, isSAT = DPLL(fRightPrime)
	if isSAT {
		fmt.Println("0 clauses, returning true")
		fmt.Println("truth assigns", retSAT.Vars)
		return retSAT, isSAT
	}
	// fmt.Println("Split on:", literal, "right failed")

	// fmt.Println("Doesn't work, returning False", retSAT)
	return nil, false
}

func UnitPropagate(f *SATInstance) {
	for {
		toRemove := 0
		for _, clause := range f.Clauses {
			if len(clause) == 1 {
				for k := range clause {
					toRemove = k
				}
				break
			}
		}
		// fmt.Println("unit propping", f.PrintClauses())
		// fmt.Println("removing", toRemove)
		if toRemove == 0 {
			break
		} else if toRemove < 0 {
			f.Vars[toRemove*-1] = false
		} else {
			f.Vars[toRemove] = true
		}
		newClauses := []map[int]bool{}
		for _, clause := range f.Clauses {
			_, containsVal := clause[toRemove]
			if containsVal {
				continue
			}
			_, containsNegVal := clause[-toRemove]
			if containsNegVal {
				delete(clause, -toRemove)
			}
			// if len(clause) == 0 {
			// 	fmt.Println("Unsat, clause empty")
			// }
			newClauses = append(newClauses, clause)
		}
		f.Clauses = newClauses
	}
}

func PureLiteralElim(f *SATInstance) {
	// map int to bool
	// if seen, put in as true
	// check if opp seen, if so make it false and opp false
	// at end, go through map, pure literal elim all true valued
	for {
		pureLiterals := make(map[int]bool, 0)
		for _, clause := range f.Clauses {
			for variable := range clause {
				_, containsVal := pureLiterals[-variable]
				if containsVal {
					pureLiterals[-variable] = false
					pureLiterals[variable] = false
				} else {
					pureLiterals[variable] = true
				}
			}
		}
		// log.Println("clauses:\n", f.PrintClauses())
		// log.Println("pure literals", pureLiterals)
		noChanges := true
		for literal, isPure := range pureLiterals {
			if !isPure {
				continue
			}
			noChanges = false
			if literal > 0 {
				f.Vars[literal] = true
			} else {
				f.Vars[-literal] = false
			}
			newClauses := []map[int]bool{}
			for _, clause := range f.Clauses {
				_, containsVal := clause[literal]
				if containsVal {
					continue
				}
				newClauses = append(newClauses, clause)
			}
			f.Clauses = newClauses
		}
		if noChanges {
			break
		}
	}
}

func SplittingRule(f *SATInstance) int {
	for _, clause := range f.Clauses {
		for variable := range clause {
			return variable
		}
	}
	log.Fatal("splitting went wrong", f.PrintClauses())
	return 0
}
