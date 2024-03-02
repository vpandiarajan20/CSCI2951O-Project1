package pkg

import (
	"fmt"
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

var CountFunc = 1
var Testing = true

func DPLL(f *SATInstance) (*SATInstance, bool) {
	if len(f.Clauses) == 0 {
		// fmt.Println("0 clauses, returning true")
		// fmt.Println("truth assigns", f.Vars)
		return f, true
	}

	// fPrime := DeepCopySATInstance(*f) // couldn't be asked to backtrack

	// fmt.Println("Pre-UnitProp", f)
	// fmt.Println("Pre-UnitProp")
	f.AssignmentStack.PushEmpty()
	f.ClauseStack.PushEmpty()
	UnitPropagate(f)
	// fmt.Println("Post-UnitProp, Pre-Literal")
	// fmt.Println("Post-UnitProp, Pre-Literal", f)
	PureLiteralElim(f)
	// fmt.Println("Post-Literal, Pre-Split")
	// fmt.Println("Post-Literal, Pre-Split", f)
	// fmt.Println("Post-Literal, Pre-Split, Satisfied Clauses", f.ClauseStack.elements)
	// fmt.Println("Post-Literal, Pre-Split, Decisions", f.AssignmentStack.elements)

	// checking for forced false clause unsat
	allClausesSatified := determineAllClauses(f)
	if allClausesSatified == True {
		return f, true
	} else if allClausesSatified == False {
		fmt.Println("Backtracking because conflict cause")
		backtrack(f)
		return nil, false
	}

	// check if all clauses satisfied, if so, return SAT

	Var, varVal := SplittingRule(f)
	fmt.Println("Split on:", Var)

	literalToAdd := int(Var)

	if !varVal {
		literalToAdd *= -1
	}

	newClause := make(map[int]bool, 0)
	newClause[literalToAdd] = false
	f.AddClause(newClause)
	retSAT, isSAT := DPLL(f)
	if isSAT {
		// fmt.Println("0 clauses, returning true")
		// fmt.Println("truth assigns", retSAT.Vars)
		return retSAT, isSAT
	}
	lastClause := f.RemoveLastClause()
	lastLevelAssigned, _ := f.AssignmentStack.Peek()
	fmt.Println("removed clause", lastClause, "split on ", literalToAdd, "assignment stack", lastLevelAssigned.PropagatedVariables)
	// for k := range lastClause {
	// 	f.Vars[uint(abs(k))] = Unassigned
	// }

	_, correctClauseRemoved := lastClause[literalToAdd]
	if !correctClauseRemoved && len(lastClause) == 1 {
		fmt.Println("Added clause:", newClause, "does not match removed clause when backtrcking", correctClauseRemoved)
	}

	fmt.Println("Split on:", Var, "left failed")

	newClause = make(map[int]bool, 0)
	newClause[-literalToAdd] = false
	f.AddClause(newClause)
	retSAT, isSAT = DPLL(f)
	if isSAT {
		// fmt.Println("0 clauses, returning true")
		// fmt.Println("truth assigns", retSAT.Vars)
		return retSAT, isSAT
	}
	// for k := range lastClause {
	// 	f.Vars[uint(abs(k))] = Unassigned
	// }

	_, correctClauseRemoved = lastClause[-literalToAdd]
	if !correctClauseRemoved && len(lastClause) == 1 {
		fmt.Println("Added clause:", newClause, "does not match removed clause when backtrcking", correctClauseRemoved)
	}
	fmt.Println("Split on:", Var, "right failed")

	fmt.Println("Clauses", f.Clauses)
	fmt.Println("Unsatisfied Clauses", f.UnsatisfiedClauses)

	// then backtrack
	fmt.Println("Backtracking b/c both branches failed")
	backtrack(f)

	return nil, false
}

func backtrack(f *SATInstance) {
	currLevelVars, isSuccessful := f.AssignmentStack.Pop()
	if isSuccessful {
		for _, Var := range currLevelVars.PropagatedVariables {
			f.Vars[Var] = Unassigned
		}
	}
	currLevelClauses, isSuccessful := f.ClauseStack.Pop()
	if isSuccessful {
		for _, clauseNum := range currLevelClauses.PropagatedVariables {
			if Testing {
				_, isFound := f.UnsatisfiedClauses[clauseNum]
				if isFound {
					log.Panicln("trying to add a clause to unsatisfiedClauses that is already there")
				}
			}
			f.AddClauseBacktrack(f.Clauses[clauseNum])
		}
	}
}

func UnitPropagate(f *SATInstance) {
	for {
		toRemove := 0
		for clauseNum := range f.UnsatisfiedClauses {
			clause := f.Clauses[clauseNum]
			if Testing {
				clauseState := determineClause(f, clause)
				if clauseState == True || clauseState == False {
					lastAssignments, _ := f.AssignmentStack.Peek()
					fmt.Println("Decision Stack", lastAssignments.PropagatedVariables)
					for literal := range clause {
						fmt.Println("literal:", literal, "val", f.Vars[uint(abs(literal))])
					}
					// TODO: somehow all variables in here are assigned even though one is not in the decsion stack
					log.Panic("satisfied/unsatisfied clause in unsatisfied clauses, state:", clauseState, ", clause:", clause)
				}
			}
			isUnit, unit := isUnitClause(f, clause)
			if isUnit {
				toRemove = unit // this is mega braindead idk how to replac
				break
			}
		}
		// fmt.Println("unit propping", f.PrintClauses())
		fmt.Println("Unit Propping:", toRemove)
		if toRemove == 0 {
			break
			// break if no unit clause
		} else if toRemove < 0 {
			f.Vars[uint(abs(toRemove))] = False
			// all variables stored positive in map
		} else {
			f.Vars[uint(abs(toRemove))] = True
		}
		for i := range f.UnsatisfiedClauses {
			// remove clause from unsatisfied clauses if it contains the value
			clause := f.Clauses[i]
			_, containsVal := clause[toRemove]
			if containsVal {
				f.RemoveClauseFromCount(clause)
				delete(f.UnsatisfiedClauses, i)
				currLevelClauses, doesExist := f.ClauseStack.Pop()
				if !doesExist {
					log.Panicln("nothing on clause stack in unit prop")
				}
				currLevelClauses.PropagatedVariables = append(currLevelClauses.PropagatedVariables, uint(i))
				f.ClauseStack.Push(currLevelClauses)
				continue // can have both value and negation, but then still remove
			}
		}
		currLevelAssignments, doesExist := f.AssignmentStack.Pop()
		if !doesExist {
			log.Panicln("nothing on assignment stack in unit prop")
		}
		currLevelAssignments.PropagatedVariables = append(currLevelAssignments.PropagatedVariables, uint(abs(toRemove)))
		f.AssignmentStack.Push(currLevelAssignments)
		// f.Clauses = newClauses
	}
}

func PureLiteralElim(f *SATInstance) {
	// always shows up in same parity - remove all clauses with that literal
	for {
		pureLiterals := make([]int, 0)
		for k, v := range f.VarCount {
			if f.Vars[k] != Unassigned {
				continue
			}
			if v.NegCount == 0 && v.PosCount > 0 {
				pureLiterals = append(pureLiterals, int(k))
			} else if v.PosCount == 0 && v.NegCount > 0 {
				pureLiterals = append(pureLiterals, -int(k))
			}
		}
		if len(pureLiterals) == 0 {
			return
		}
		for _, literal := range pureLiterals {
			fmt.Println(literal, "is a pure literal")
			// actually filling out var
			if literal > 0 {
				f.Vars[uint(abs(literal))] = True
			} else {
				f.Vars[uint(abs(literal))] = False
			}
			for i := range f.UnsatisfiedClauses {
				clause := f.Clauses[i]
				_, containsVal := clause[literal]
				if containsVal {
					literalCount := f.VarCount[uint(abs(literal))].NegCount + f.VarCount[uint(abs(literal))].PosCount
					f.RemoveClauseFromCount(clause)
					if f.VarCount[uint(abs(literal))].NegCount+f.VarCount[uint(abs(literal))].PosCount != literalCount-1 {
						fmt.Println("literal", literal, "clause:", clause)
						fmt.Println("VarCountNow", f.VarCount[uint(abs(literal))], "clause:", clause)
						log.Panicln("varcounts are messed up?")
					}
					delete(f.UnsatisfiedClauses, i)
					currLevelClauses, doesExist := f.ClauseStack.Pop()
					if !doesExist {
						log.Panicln("nothing on clause stack in pure literal elim")
					}
					currLevelClauses.PropagatedVariables = append(currLevelClauses.PropagatedVariables, uint(i))
					f.ClauseStack.Push(currLevelClauses)
					continue
				}
			}
			currLevelAssignments, doesExist := f.AssignmentStack.Pop()
			if !doesExist {
				log.Panicln("nothing on assignment stack in pure literal elim")
			}
			currLevelAssignments.PropagatedVariables = append(currLevelAssignments.PropagatedVariables, uint(abs(literal)))
			f.AssignmentStack.Push(currLevelAssignments)
		}
	}
}

func SplittingRule(f *SATInstance) (uint, bool) {

	keys := make([]uint, len(f.VarCount))
	i := 0
	for k := range f.VarCount {
		if f.Vars[k] == Unassigned {
			keys[i] = k
			i++
		}
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
				return uint(abs(variable)), true
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
		if validLiterals == 0 {
			return 0, false
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

func determineAllClauses(f *SATInstance) int {
	for _, clause := range f.Clauses {
		if determineClause(f, clause) == False {
			return False
			// if any clause is false, formula is false
		}
	}
	for _, clause := range f.Clauses {
		if determineClause(f, clause) == Unassigned {
			return Unassigned
			// returns Unassigned if there is an unassigned variable and not alr False
		}
	}
	// if all clauses are true, formula is true
	return True
}

func determineClause(f *SATInstance, clause map[int]bool) int {

	for variable := range clause { // iterate through all variables in clause
		if (variable > 0 && f.Vars[uint(abs(variable))] == True) || (variable < 0 && f.Vars[uint(abs(variable))] == False) {
			// if any variable is true and shows up as pos || any var is false and shows up as neg, clause is true
			return True
		}
	}
	for variable := range clause {
		if f.Vars[uint(abs(variable))] == Unassigned {
			// returns Unassigned if there is an unassigned variable
			return Unassigned
		}
	}
	// if all variables are false, clause is false
	return False
}

func isUnitClause(f *SATInstance, clause map[int]bool) (bool, int) {
	numFalses := 0
	numUnassigned := 0
	litUnassigned := 0
	for literal := range clause {
		if (literal < 0 && f.Vars[uint(abs(literal))] == True) || (literal > 0 && f.Vars[uint(abs(literal))] == False) {
			numFalses += 1
		} else if f.Vars[uint(abs(literal))] == Unassigned {
			numUnassigned += 1
			litUnassigned = literal
			if numUnassigned > 1 {
				return false, 0
			}
		}
	}
	if numFalses == len(clause)-1 && numUnassigned == 1 {
		return true, litUnassigned
	}
	return false, 0
}
