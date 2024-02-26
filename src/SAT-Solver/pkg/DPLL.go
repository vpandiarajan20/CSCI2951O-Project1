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

const (
	NoChange = iota
	SuccessfulChange
	FailedChange
)

var CountFunc = 3

func DPLL(f *SATInstance) bool, error {
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

	isSuccessful, err := unitPropagate(f, 0)
	if !isSuccessful {
		return false, nil
	}
	if err != nil {
		return 
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
func moveAllWatchedLiterals(f *SATInstance, wlToChange int) (bool, error) {
	if wlToChange <= 0 {
		return false, errors.New("wlToChange must be greater than or equal to zero")
	}
	for {
		finished := true
		for _, clauseNum := range f.LiteralToClauses[wlToChange] {
			change, err := moveWatchedLiteral(f, wlToChange, clauseNum)
			// check in moveWatchedLiteral if actually need to move

			if err != nil || change == FailedChange {
				return false, err
			}
			if change == SuccessfulChange {
				// if change something then f.LiteralToClauses[wlToChange] is stale so break
				finished = false
				break
			}
		}

		if finished {
			// if iterates through all the watched clauses to change and don't need to change any of them finish
			break
		}
	}
	return true, nil
}

func moveWatchedLiteral(f *SATInstance, wlToChange, clauseNumber int) (int, error) {
	wl := f.WatchedLiterals[clauseNumber]

	if wlToChange <= 0 {
		return FailedChange, errors.New("wlToChange must be greater than or equal to zero")
	}

	_, isPresent := f.Clauses[clauseNumber][wlToChange]
	if isPresent && (f.Vars[abs(wlToChange)] == True) {
		return NoChange, nil
	}

	_, isPresent = f.Clauses[clauseNumber][-wlToChange]
	if isPresent && (f.Vars[abs(wlToChange)] == False) {
		return NoChange, nil
	}

	// if we reach this point, it means the wlToChange has an unsatisfying assignment in this clause
	changeWL1 := abs(wl.Literal1) == wlToChange

	for _, literal := range f.Clauses[clauseNumber] {
		if abs(literal) == abs(wlToChange) {
			continue
		}
		// if we are changing WL1, we don't want it to point to WL2
		if changeWL1 && wl.Literal2 == literal {
			continue
		}
		// if we are changing WL2, we don't wnt it to point to WL1
		if !changeWL1 && wl.Literal1 == literal {
			continue
		}

		if f.Vars[abs(literal)] == Unassigned || (literal > 0 && f.Vars[abs(literal)] == True) || (f.Vars[abs(literal)] == False && literal < 0) {
			if changeWL1 {
				// var name is bad
				newLiteralToClauses, err := removeElement(f.LiteralToClauses[wl.Literal1], clauseNumber)
				f.LiteralToClauses[wl.Literal1] = newLiteralToClauses
				if err != nil {
					return FailedChange, err
				}
				wl.Literal1 = literal
			} else {
				newLiteralToClauses, err := removeElement(f.LiteralToClauses[wl.Literal2], clauseNumber)
				f.LiteralToClauses[wl.Literal2] = newLiteralToClauses
				if err != nil {
					return FailedChange, err
				}
				wl.Literal2 = literal
			}

			f.WatchedLiterals[clauseNumber] = wl
			f.LiteralToClauses[literal] = append(f.LiteralToClauses[literal], clauseNumber)

			return SuccessfulChange, nil
		}
	}

	// called both when both variables are assigned + when one variable is unassigned
	return resolveImplication(f, clauseNumber, changeWL1)
}

func unitPropagate(f *SATInstance, literalProp int) (bool, error) {
	if literalProp < 0 {
		return false, errors.New("Literal has to Be Non-negative Integer")
	}

	unitClauses, isFound := f.LiteralToClauses[literalProp]
	if !isFound {
		return true, nil
	}
	for _, clauseNum := range unitClauses {

		// checking if the propagating literal satifies the clause
		_, isFound := f.Clauses[clauseNum][literalProp]
		if isFound && f.Vars[literalProp] == True {
			continue
		}
		_, isFound = f.Clauses[clauseNum][-literalProp]
		if isFound && f.Vars[literalProp] == False {
			continue
		}

		moveSuccesful, err := moveAllWatchedLiterals(f, literalProp)
		return moveSuccesful, err
		// moveSuccessful is if fails => then backtrack/unsat if init
	}
	return true, nil
}

// should never be called, except after moving a literal
func resolveImplication(f *SATInstance, clauseNum int, changeW1 bool) (int, error) {
	wl := f.WatchedLiterals[clauseNum]
	isLiteral1Satisfying := (f.Vars[abs(wl.Literal1)] == True && wl.Literal1 > 0) || (f.Vars[abs(wl.Literal1)] == False && wl.Literal1 < 0)
	isLiteral2Satisfying := (f.Vars[abs(wl.Literal2)] == True && wl.Literal2 > 0) || (f.Vars[abs(wl.Literal2)] == False && wl.Literal2 < 0)
	isLiteral1Unassigned := f.Vars[abs(wl.Literal1)] == Unassigned
	isLiteral2Unassigned := f.Vars[abs(wl.Literal2)] == Unassigned

	successfulProp, err := false, error(nil)
	if changeW1 {
		// L2 unassigned => then set value of L2 and unit propogate
		// if L2 is assigned poorly => then backtrack
		// if L2 is assigned well => then we are fine

		if !isLiteral2Satisfying && !isLiteral2Unassigned { //
			return FailedChange, nil
		}
		if isLiteral2Satisfying {
			return NoChange, nil
		}
		if wl.Literal2 > 0 {
			f.Vars[abs(wl.Literal2)] = True
		} else {
			f.Vars[abs(wl.Literal2)] = False
		}
		successfulProp, err = unitPropagate(f, wl.Literal2)
	} else {
		// if L1 is assigned well => then we are fine
		if !isLiteral1Satisfying && !isLiteral1Unassigned {
			return FailedChange, nil
		}
		if isLiteral1Satisfying {
			return NoChange, nil
		}
		if wl.Literal1 > 0 {
			f.Vars[abs(wl.Literal1)] = True
		} else {
			f.Vars[abs(wl.Literal1)] = False
		}
		successfulProp, err = unitPropagate(f, wl.Literal1)
	}
	if successfulProp {
		return SuccessfulChange, err
	} else {
		return FailedChange, err
	}
}

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

// removeElement removes the first occurrence of elem from slice and returns the modified slice.
func removeElement(slice []int, elem int) ([]int, error) {
	for i, v := range slice {
		if v == elem {
			// Swap the element to remove with the last element in the slice.
			slice[i] = slice[len(slice)-1]
			// Truncate the slice by one to remove the last element.
			return slice[:len(slice)-1], nil
		}
	}
	// If the element was not found, return the original slice.
	return nil, errors.New("tried to remove an element not in the slice")
}
