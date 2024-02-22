package pkg

import (
	"errors"
	"log"
)

const (
	DLCS = iota
	DLIS
	RDLCS
	RDLIS
)

var CountFunc = 3

func DPLL(f *SATInstance) bool {
	PureLiteralElim(f)

	// checking for empty clause unsat
	for _, clause := range f.Clauses {
		if len(clause) == 0 {
			return false
		}
	}
	for i, clause := range f.Clauses {
		literal1 := 0
		literal2 := 0

		for j, literal := range clause {
			if j == 0 {
				literal1 = literal
			}
			if j == 1 {
				literal2 = literal
			}
		}
		// setting watched literals
		wl := f.WatchedLiterals[i]
		wl.Literal1 = literal1
		wl.Literal2 = literal2

		// setting struct vals
		f.WatchedLiterals[i] = wl
		f.LiteralToClauses[wl.Literal1] = append(f.LiteralToClauses[wl.Literal1], i)
		f.LiteralToClauses[wl.Literal2] = append(f.LiteralToClauses[wl.Literal2], i)
	}

	isSuccessful := initialUnitPropagate(f)
	if !isSuccessful {
		return false
	}

	// if stack is empty, and we try to pop, then unsat

	// literal, literalVal := SplittingRule(fPrime)
	// // fmt.Println("Split on:", literal)

	// fRightPrime := DeepCopySATInstance(*fPrime)

	// if !literalVal {
	// 	literal *= -1
	// }

	// newClause := make(map[int]int, 0)
	// newClause[literal] = False
	// fPrime.AddClause(newClause)
	// retSAT, isSAT := DPLL(fPrime)
	// if isSAT {
	// 	// fmt.Println("0 clauses, returning true")
	// 	// fmt.Println("truth assigns", retSAT.Vars)
	// 	return retSAT, isSAT
	// }
	// // fmt.Println("Split on:", literal, "left failed")

	// newClause = make(map[int]int, 0)
	// newClause[-literal] = False
	// fRightPrime.AddClause(newClause)
	// retSAT, isSAT = DPLL(fRightPrime)
	// if isSAT {
	// 	// fmt.Println("0 clauses, returning true")
	// 	// fmt.Println("truth assigns", retSAT.Vars)
	// 	return retSAT, isSAT
	// }
	// // fmt.Println("Split on:", literal, "right failed")

	// fmt.Println("Doesn't work, returning False", retSAT)
	return false
}

//	func propagateWatchedLiteral(f *SATInstance) bool {
//		for literal, clauses := range f.LiteralToClauses {
//			if f.Vars[abs(literal)] == Unassigned {
//				continue
//			}
//			for _, clauseNum := range clauses {
//				// check if literl needs to be moved
//				// if so
//				// call moveWatchedLitrl
//			}
//		}
//	}
//
// every

// wlToChange MUST be positive
func moveAllWatchedLiterals(f *SATInstance, wlToChange int) (map[int]bool, error) {
	if wlToChange <= 0 {
		return nil, errors.New("wlToChange must be greater than or equal to zero")
	}
	successfulMoves := make(map[int]bool, 0)
	for _, clauseNum := range f.LiteralToClauses[wlToChange] {
		successfulMove, err := moveWatchedLiteral(f, wlToChange, clauseNum)
		if err != nil {
			return nil, err
		}
		successfulMoves[clauseNum] = successfulMove
	}
	return successfulMoves, nil
}

func moveWatchedLiteral(f *SATInstance, wlToChange, clauseNumber int) (bool, error) {
	wl := f.WatchedLiterals[clauseNumber]
	changeWL1 := wl.Literal1 == wlToChange

	if wlToChange <= 0 {
		return false, errors.New("wlToChange must be greater than or equal to zero")
	}

	_, isPresent := f.Clauses[clauseNumber][wlToChange]
	if isPresent && (f.Vars[abs(wlToChange)] == True) {
		return true, nil
	}

	_, isPresent = f.Clauses[clauseNumber][-wlToChange]
	if isPresent && (f.Vars[abs(wlToChange)] == False) {
		return true, nil
	}

	for _, literal := range f.Clauses[clauseNumber] {
		if abs(literal) == abs(wlToChange) {
			continue
		}
		// if we are changing WL1, we don't want it to point to WL2
		if changeWL1 && wl.Literal2 == literal {
			continue
		}
		// if we are changing WL2, we don't want it to point to WL1
		if !changeWL1 && wl.Literal1 == literal {
			continue
		}
		if f.Vars[abs(literal)] == Unassigned || (literal > 0 && f.Vars[abs(literal)] == True) || (f.Vars[abs(literal)] == False && literal < 0) {
			if changeWL1 {
				wl.Literal1 = literal
			} else {
				wl.Literal2 = literal
			}
			return true, nil
		}
	}
	return false, nil
}

func initialUnitPropagate(f *SATInstance) (bool, error) {
	literalsToChange := make([]int, 0)
	unitClauses, isFound := f.LiteralToClauses[0]
	if !isFound {
		return true, nil
	}
	for _, clauseNum := range unitClauses {
		literal := f.WatchedLiterals[clauseNum].Literal1
		if f.Vars[abs(literal)] != Unassigned && ((f.Vars[abs(literal)] != True && literal > 0) || (f.Vars[abs(literal)] != False && literal < 0)) {
			return false, nil
		}
		switch literal > 0 {
		case true:
			f.Vars[abs(literal)] = True
		case false:
			f.Vars[abs(literal)] = False
		}
		literalsToChange = append(literalsToChange, literal)
	}
	for _, literal := range literalsToChange {
		implicationMap, err := moveAllWatchedLiterals(f, abs(literal))
		if err != nil {
			return false, err
		}
		if !resolveImplications(f, implicationMap) {
			return false, nil
		}
		// might want to move resolveImplications to end of moveWatchedLiteral
	}

	return true, nil
}

func resolveImplications(f *SATInstance, implicationMap map[int]bool) bool {
	for clauseNum, hasMoved := range implicationMap {
		if hasMoved {
			continue
		}
		wl := f.WatchedLiterals[clauseNum]
		if (f.Vars[abs(wl.Literal2)] == True && wl.Literal2 > 0) ||
			(f.Vars[abs(wl.Literal2)] != False && wl.Literal2 < 0) ||
			(f.Vars[abs(wl.Literal1)] == True && wl.Literal1 > 0) ||
			(f.Vars[abs(wl.Literal1)] != False && wl.Literal1 < 0) {
			continue
		}
		changeWL1 := wl.Literal1 == Unassigned

	}
}

// 0 1
// -1 2
// -2 -1

// func UnitPropagate(f *SATInstance) {
// 	for {
// 		toRemove := 0
// 		for _, clause := range f.Clauses {
// 			if len(clause) == 1 {
// 				//first key of clause map
// 				for k := range clause {
// 					toRemove = k
// 				}
// 				break
// 			}
// 		}
// 		// fmt.Println("unit propping", f.PrintClauses())
// 		// fmt.Println("removing", toRemove)
// 		if toRemove == 0 {
// 			break
// 			// break if no unit clause
// 		} else if toRemove < 0 {
// 			f.Vars[toRemove*-1] = false
// 			// all variables stored positive in map
// 		} else {
// 			f.Vars[toRemove] = true
// 		}
// 		newClauses := []map[int]bool{}
// 		for _, clause := range f.Clauses {
// 			// remove clause if it contains the value
// 			_, containsVal := clause[toRemove]
// 			if containsVal {
// 				f.RemoveClauseFromCount(clause)
// 				continue // can have both value and negation, but then still remove
// 			}
// 			// remove value from clause if it contains the negation
// 			_, containsNegVal := clause[-toRemove]
// 			if containsNegVal {
// 				f.RemoveLiteralFromCount(-toRemove)
// 				delete(clause, -toRemove)
// 			}
// 			// if len(clause) == 0 {
// 			// 	fmt.Println("Unsat, clause empty")
// 			// }
// 			newClauses = append(newClauses, clause)
// 		}
// 		f.Clauses = newClauses
// 	}
// }

func PureLiteralElim(f *SATInstance) {
	for {
		pureLiterals := make(map[int]bool, 0)
		for _, clause := range f.Clauses {
			for variable := range clause {
				// checking to see if the negation is present in another clause
				_, containsVal := pureLiterals[-variable]
				if containsVal {
					// if present, set both pos and neg to false
					pureLiterals[-variable] = false
					pureLiterals[variable] = false
				} else {
					// if not present, set to true
					pureLiterals[variable] = true
				}
			}
		}
		noChanges := true
		for literal, isPure := range pureLiterals {
			if !isPure {
				// skip entries that are not pure
				continue
			}
			noChanges = false
			// assign truth vals to literals
			if literal > 0 {
				f.Vars[literal] = True
			} else {
				f.Vars[-literal] = False
			}
			newClauses := []map[int]int{}
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

	// keys := make([]int, len(f.VarCount))
	// i := 0
	// for k := range f.VarCount {
	// 	keys[i] = k
	// 	i++
	// }
	// switch CountFunc {
	// case DLCS, RDLCS:
	// 	sort.SliceStable(keys, func(i, j int) bool {
	// 		iCounts := f.VarCount[keys[i]]
	// 		jCounts := f.VarCount[keys[j]]
	// 		return (iCounts.NegCount + iCounts.PosCount) > (jCounts.NegCount + jCounts.PosCount)
	// 		// counts stored in struct with NegCount and PosCount
	// 	})
	// case DLIS, RDLIS:
	// 	sort.SliceStable(keys, func(i, j int) bool {
	// 		iCounts := f.VarCount[keys[i]]
	// 		jCounts := f.VarCount[keys[j]]
	// 		return max(iCounts.NegCount, iCounts.PosCount) > max(jCounts.NegCount, jCounts.PosCount)
	// 	})
	// default:
	// 	for _, clause := range f.Clauses {
	// 		for variable := range clause {
	// 			return variable, true
	// 		}
	// 	}
	// }
	// switch CountFunc {
	// case DLCS, DLIS:
	// 	return keys[0], f.VarCount[keys[0]].PosCount > f.VarCount[keys[0]].NegCount // explore true or false
	// case RDLCS, RDLIS: // uniform at random first 5
	// 	validLiterals := 0
	// 	for i := 0; i < 5; i++ { // messed up if varcount less than 5 but like
	// 		iCounts := f.VarCount[keys[i]]
	// 		if iCounts.NegCount+iCounts.PosCount > 0 {
	// 			validLiterals += 1
	// 		}
	// 	}
	// 	keyToReturn := keys[rand.Intn(validLiterals)]
	// 	return keyToReturn, f.VarCount[keyToReturn].PosCount > f.VarCount[keyToReturn].NegCount
	// }
	for _, clause := range f.Clauses {
		for variable := range clause {
			return variable, true
		}
	}
	log.Fatal("splitting went wrong", f.PrintClauses())
	return 0, false
}

// func max(a, b int) int {
// 	if a > b {
// 		return a
// 	}
// 	return b
// }
