package pkg

import (
	"errors"
	"fmt"
	"math/rand"
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

func DPLL(f *SATInstance) (bool, error) {

	triviallySuccessful, err := preprocessFormula(f)
	if !triviallySuccessful || err != nil {
		return false, err
	}

	for {
		literal, literalVal, doneSplitting := SplittingRule(f)

		if doneSplitting || allClausesSatisfied(f) { // can prolly get rid of one of these conditions
			// out of variables to split on or only split on variables that are in alr satisfied clauses
			return true, nil
		}
		// fmt.Println("Split on:", literal, "set to", literalVal)

		f.Vars[uint(abs(literal))] = True
		if !literalVal {
			f.Vars[uint(abs(literal))] = False
		}
		f.StackAssignments.Push(literal, true, false)
		isSuccessful, err := unitPropagate(f, literal)
		if err != nil {
			return false, err
		}
		for !isSuccessful {
			// backtracks until branch where has not tried both true and false
			prevVarAssignment := VariableAssignment{}
			successfulPop := false
			for {
				prevVarAssignment, successfulPop = f.StackAssignments.Pop()
				if !successfulPop {
					// if nothing else on stack, then unsat
					return false, nil
				}
				if prevVarAssignment.IsBranch {
					// fmt.Println("Unsetting branch: ", prevVarAssignment.Literal, " val:", f.Vars[uint(prevVarAssignment.Literal)])
				}
				if prevVarAssignment.TriedBothWays {
					// fmt.Println(prevVarAssignment.Literal, "has been tried both ways")
				}
				if prevVarAssignment.IsBranch && !prevVarAssignment.TriedBothWays {
					break
				}
				f.Vars[uint(abs(prevVarAssignment.Literal))] = Unassigned
			}
			currAssignment := f.Vars[uint(abs(prevVarAssignment.Literal))]
			oppositeAssignment := False
			if currAssignment == False {
				oppositeAssignment = True
			} else if currAssignment == Unassigned {
				return false, errors.New("literal should be assigned to True or False if Branched on it")
			}
			f.Vars[uint(abs(prevVarAssignment.Literal))] = oppositeAssignment
			// fmt.Println("setting ", prevVarAssignment.Literal, " val:", f.Vars[uint(prevVarAssignment.Literal)])
			f.StackAssignments.Push(prevVarAssignment.Literal, true, true)
			isSuccessful, err = unitPropagate(f, prevVarAssignment.Literal)
			if err != nil {
				return false, err
			}
		}
	}
}

func preprocessFormula(f *SATInstance) (bool, error) {
	PureLiteralElim(f)
	SetWatchedLiterals(f)

	// checking for empty clause unsat - should never happen
	for _, clause := range f.Clauses {
		if len(clause) == 0 {
			return false, errors.New("empty clause after pure literal elim - parser messed up")
		}
	}
	// 0 is set to a positive literal, so this should never be the decisive variable in a clause

	isSuccessful, err := unitPropagate(f, 0)
	if !isSuccessful || err != nil {
		fmt.Println("UNSAT After Initial Unit Propagation")
		return false, nil
	}
	return true, nil
}

func SetWatchedLiterals(f *SATInstance) {
	for i, clause := range f.Clauses {
		literal1 := 0
		literal2 := 0

		j := 0
		for literal := range clause {
			if j == 0 {
				literal1 = literal
			}
			if j == 1 {
				literal2 = literal
			}
			j += 1
		}
		// setting watched literals
		f.WatchedLiterals[uint(abs(i))] = struct{ Literal1, Literal2 int }{Literal1: literal1, Literal2: literal2}

		// setting struct vals
		f.LiteralToClauses[uint(abs(literal1))] = append(f.LiteralToClauses[uint(abs(literal1))], i)
		f.LiteralToClauses[uint(abs(literal2))] = append(f.LiteralToClauses[uint(abs(literal2))], i)
	}

	f.Vars[0] = False
}
func allClausesSatisfied(f *SATInstance) bool {
	for i := range f.Clauses {
		literal1 := f.WatchedLiterals[uint(abs(i))].Literal1
		literal2 := f.WatchedLiterals[uint(abs(i))].Literal2

		if (literal1 >= 0 && f.Vars[uint(abs(literal1))] == True) || (f.Vars[uint(abs(literal1))] == False && literal1 < 0) {
			continue
		} else if (literal2 >= 0 && f.Vars[uint(abs(literal2))] == True) || (f.Vars[uint(abs(literal2))] == False && literal2 < 0) {
			continue
		} else {
			return false
		}
	}
	return true
}

// wlToChange MUST be positive
func moveAllWatchedLiterals(f *SATInstance, wlToChange int) (bool, error) {
	if wlToChange < 0 {
		return false, errors.New("wlToChange must be greater than or equal to zero")
	}
	for {
		finished := true
		for _, clauseNum := range f.LiteralToClauses[uint(abs(wlToChange))] {
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

func moveWatchedLiteral(f *SATInstance, wlToChange int, clauseNumber int) (int, error) {
	wl := f.WatchedLiterals[uint(clauseNumber)]

	if wlToChange < 0 {
		return FailedChange, errors.New("wlToChange must be greater than or equal to zero")
	}

	_, isPresent := f.Clauses[clauseNumber][wlToChange]
	if isPresent && (f.Vars[uint(abs(wlToChange))] == True) {
		return NoChange, nil
	}

	_, isPresent = f.Clauses[clauseNumber][-wlToChange]
	if isPresent && (f.Vars[uint(abs(wlToChange))] == False) {
		return NoChange, nil
	}

	// if we reach this point, it means the wlToChange has an unsatisfying assignment in this clause
	changeWL1 := abs(wl.Literal1) == wlToChange

	for literal := range f.Clauses[clauseNumber] {
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

		if f.Vars[uint(abs(literal))] == Unassigned || (literal >= 0 && f.Vars[uint(abs(literal))] == True) || (f.Vars[uint(abs(literal))] == False && literal < 0) {
			if changeWL1 {
				// var name is bad
				newLiteralToClauses, err := removeElement(f.LiteralToClauses[uint(abs(wl.Literal1))], clauseNumber)
				f.LiteralToClauses[uint(abs(wl.Literal1))] = newLiteralToClauses
				if err != nil {
					return FailedChange, err
				}
				wl.Literal1 = literal
			} else {
				newLiteralToClauses, err := removeElement(f.LiteralToClauses[uint(abs(wl.Literal2))], clauseNumber)
				f.LiteralToClauses[uint(abs(wl.Literal2))] = newLiteralToClauses
				if err != nil {
					return FailedChange, err
				}
				wl.Literal2 = literal
			}

			f.WatchedLiterals[uint(abs(clauseNumber))] = wl
			f.LiteralToClauses[uint(abs(literal))] = append(f.LiteralToClauses[uint(abs(literal))], clauseNumber)

			return SuccessfulChange, nil
		}
	}

	// called both when both variables are assigned + when one variable is unassigned
	return resolveImplication(f, clauseNumber, changeWL1)
}

func unitPropagate(f *SATInstance, literalProp int) (bool, error) {
	if literalProp < 0 {
		return false, errors.New("literal has to Be Non-negative Integer")
	}

	unitClauses, isFound := f.LiteralToClauses[uint(abs(literalProp))]
	if !isFound {
		return true, nil
	}
	for _, clauseNum := range unitClauses {

		// checking if the propagating literal satifies the clause
		_, isFound := f.Clauses[clauseNum][literalProp]
		if isFound && f.Vars[uint(abs(literalProp))] == True {
			continue
		}
		_, isFound = f.Clauses[clauseNum][-literalProp]
		if isFound && f.Vars[uint(abs(literalProp))] == False {
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
	wl := f.WatchedLiterals[uint(abs(clauseNum))]
	isLiteral1Satisfying := (f.Vars[uint(abs(wl.Literal1))] == True && wl.Literal1 >= 0) || (f.Vars[uint(abs(wl.Literal1))] == False && wl.Literal1 < 0)
	isLiteral2Satisfying := (f.Vars[uint(abs(wl.Literal2))] == True && wl.Literal2 >= 0) || (f.Vars[uint(abs(wl.Literal2))] == False && wl.Literal2 < 0)
	isLiteral1Unassigned := f.Vars[uint(abs(wl.Literal1))] == Unassigned
	isLiteral2Unassigned := f.Vars[uint(abs(wl.Literal2))] == Unassigned

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
		if wl.Literal2 >= 0 {
			f.Vars[uint(abs(wl.Literal2))] = True
			f.StackAssignments.Push(abs(wl.Literal2), false, false) // false means not branching
		} else {
			f.Vars[uint(abs(wl.Literal2))] = False
			f.StackAssignments.Push(abs(wl.Literal2), false, false)
		}
		successfulProp, err = unitPropagate(f, abs(wl.Literal2))
	} else {
		// if L1 is assigned well => then we are fine
		if !isLiteral1Satisfying && !isLiteral1Unassigned {
			return FailedChange, nil
		}
		if isLiteral1Satisfying {
			return NoChange, nil
		}
		if wl.Literal1 >= 0 {
			f.Vars[uint(abs(wl.Literal1))] = True
			f.StackAssignments.Push(abs(wl.Literal1), false, false)
		} else {
			f.Vars[uint(abs(wl.Literal1))] = False
			f.StackAssignments.Push(abs(wl.Literal1), false, false)
		}
		successfulProp, err = unitPropagate(f, abs(wl.Literal1))
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
			if literal >= 0 {
				f.Vars[uint(abs(literal))] = True
			} else {
				f.Vars[uint(abs(literal))] = False
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

func SplittingRule(f *SATInstance) (int, bool, bool) {

	// for _, clause := range f.Clauses {
	// 	for variable := range clause {
	// 		if f.Vars[abs(variable)] == Unassigned {
	// 			return abs(variable), true
	// 		}
	// 	}
	// }
	bestLiteral := uint(0)
	mostClauses := 0

	for literal, clauses := range f.LiteralToClauses {

		if f.Vars[literal] == Unassigned {
			if len(clauses) > mostClauses {
				mostClauses = len(clauses)
				bestLiteral = literal
			}
		}
	}
	if bestLiteral != 0 {
		return int(bestLiteral), true, false
	}
	// for i, clause := range f.Clauses {
	// 	wl := f.WatchedLiterals[uint(abs(i))]
	// 	isLiteral1Satisfying := (f.Vars[uint(abs(wl.Literal1))] == True && wl.Literal1 >= 0) || (f.Vars[uint(abs(wl.Literal1))] == False && wl.Literal1 < 0)
	// 	isLiteral2Satisfying := (f.Vars[uint(abs(wl.Literal2))] == True && wl.Literal2 >= 0) || (f.Vars[uint(abs(wl.Literal2))] == False && wl.Literal2 < 0)
	// 	if isLiteral1Satisfying || isLiteral2Satisfying {
	// 		continue
	// 	}
	// 	// should i be returning a variable that is a watched literal
	// 	for variable := range clause {
	// 		if f.Vars[uint(abs(variable))] == Unassigned {
	// 			return abs(variable), variable > 0, false
	// 		}
	// 	}
	// }
	return 0, false, true
	// log.Fatal("splitting went wrong", f.PrintClauses())
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
	return nil, errors.New("tried to remove an element not in the slice")
}

func pickRandomKey(m map[int]int) int {
	// Get the slice of keys.
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	// Generate a random index.
	i := rand.Intn(len(keys))

	// Return the key at the random index.
	return keys[i]
}
