package pkg

import (
	"log"
	"math/rand"
	"sort"
)

const (
	DLCS = iota
	DLIS
	RDLCS
	RDLIS
)

var CountFunc = 3

func DPLL(f *SATInstance) (*SATInstance, bool) {
	if len(f.Clauses) == 0 {
		// fmt.Println("0 clauses, returning true")
		// fmt.Println("truth assigns", f.Vars)
		return f, true
	}

	fPrime := DeepCopySATInstance(*f) // couldn't be asked to backtrack

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
		// fmt.Println("truth assigns", f.Vars)
		return fPrime, true
	}

	literal, literalVal := SplittingRule(fPrime)
	// fmt.Println("Split on:", literal)

	fRightPrime := DeepCopySATInstance(*fPrime)

	if !literalVal {
		literal *= -1
	}

	newClause := make(map[int]bool, 0)
	newClause[literal] = false
	fPrime.AddClause(newClause)
	retSAT, isSAT := DPLL(fPrime)
	if isSAT {
		// fmt.Println("0 clauses, returning true")
		// fmt.Println("truth assigns", retSAT.Vars)
		return retSAT, isSAT
	}
	// fmt.Println("Split on:", literal, "left failed")

	newClause = make(map[int]bool, 0)
	newClause[-literal] = false
	fRightPrime.AddClause(newClause)
	retSAT, isSAT = DPLL(fRightPrime)
	if isSAT {
		// fmt.Println("0 clauses, returning true")
		// fmt.Println("truth assigns", retSAT.Vars)
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
				//first key of clause map 
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
			// break if no unit clause
		} else if toRemove < 0 {
			f.Vars[toRemove*-1] = false
			// all variables stored positive in map
		} else {
			f.Vars[toRemove] = true
		}
		newClauses := []map[int]bool{}
		for _, clause := range f.Clauses {
			// remove clause if it contains the value
			_, containsVal := clause[toRemove]
			if containsVal {
				f.RemoveClauseFromCount(clause) 
				continue // can have both value and negation, but then still remove
			}
			// remove value from clause if it contains the negation
			_, containsNegVal := clause[-toRemove]
			if containsNegVal {
				f.RemoveLiteralFromCount(-toRemove)
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

func SplittingRule(f *SATInstance) (int, bool) {

	keys := make([]int, len(f.VarCount))
	i := 0
	for k := range f.VarCount {
		keys[i] = k
		i++
	}
	switch CountFunc {
	case DLCS, RDLCS:
		sort.SliceStable(keys, func(i, j int) bool {
			iCounts := f.VarCount[keys[i]]
			jCounts := f.VarCount[keys[j]]
			return (iCounts.NegCount + iCounts.PosCount) > (jCounts.NegCount + jCounts.PosCount)
			// counts stored in struct with NegCount and PosCount
		})
	case DLIS, RDLIS:
		sort.SliceStable(keys, func(i, j int) bool {
			iCounts := f.VarCount[keys[i]]
			jCounts := f.VarCount[keys[j]]
			return max(iCounts.NegCount, iCounts.PosCount) > max(jCounts.NegCount, jCounts.PosCount)
		})
	default:
		for _, clause := range f.Clauses {
			for variable := range clause {
				return variable, true
			}
		}
	}
	switch CountFunc {
	case DLCS, DLIS:
		return keys[0], f.VarCount[keys[0]].PosCount > f.VarCount[keys[0]].NegCount // explore true or false
	case RDLCS, RDLIS: // uniform at random first 5
		validLiterals := 0
		for i := 0; i < 5; i++ { // messed up if varcount less than 5 but like
			iCounts := f.VarCount[keys[i]]
			if iCounts.NegCount+iCounts.PosCount > 0 {
				validLiterals += 1
			}
		}
		keyToReturn := keys[rand.Intn(validLiterals)]
		return keyToReturn, f.VarCount[keyToReturn].PosCount > f.VarCount[keyToReturn].NegCount
	}

	log.Fatal("splitting went wrong", f.PrintClauses())
	return 0, false
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
