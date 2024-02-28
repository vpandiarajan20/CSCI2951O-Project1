package pkg

import (
	"errors"
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

func CDCL(f *SATInstance) (bool, error) {
	isSuccessful, err := preprocessFormula(f)
	if !isSuccessful {
		return false, err
	}

	for !allVariablesAssigned(f) {
		conflictClause, err := unitPropagate(f)
		if err != nil {
			return false, err
		}

		if conflictClause != nil {
			f.NumConflicts += 1
			log.Println("Conflict Clause", conflictClause)
			log.Println("Implication Nodes", f.ImplicationGraph)

			level, learnedClause, err := analyzeConflict(f, conflictClause)
			if err != nil {
				return false, err
			}
			if level == -1 {
				return false, nil
				// UNSAT!!
			}
			f.Clauses = append(f.Clauses, learnedClause)
			// add learned clause to formula
			backtrack(f, level)
			// backtrack to level
			f.Level = level
		} else if allVariablesAssigned(f) {
			// SAT!!
			break
		} else {
			// apply splitting rule
			varToAssign, valToAssign := SplittingRule(f)
			if valToAssign {
				f.Vars[varToAssign] = True
			} else {
				f.Vars[varToAssign] = False
			}
			f.Level += 1
			f.NumBranches += 1
			f.BranchingHist[f.Level] = varToAssign
			f.PropagateHist[f.Level] = make([]uint, 0)
			updateImplicationGraph(f, varToAssign, nil)
		}
	}
	return true, nil
}

func allVariablesAssigned(f *SATInstance) bool {
	for _, val := range f.Vars {
		if val == Unassigned {
			return false
		}
	}
	return true
}

func preprocessFormula(f *SATInstance) (bool, error) {
	pureLiteralElim(f)
	// checking for empty clause unsat - should never happen
	for _, clause := range f.Clauses {
		if len(clause) == 0 {
			return false, errors.New("empty clause after pure literal elim - parser messed up")
		}
	}
	return true, nil
}

func pureLiteralElim(f *SATInstance) {
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
			newClauses := []map[int]bool{}
			for _, clause := range f.Clauses {
				_, containsVal := clause[literal]
				// remove clauses that contain the pure literal
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

func determineClause(f *SATInstance, clause map[int]bool) int {
	for variable := range clause { // iterate through all variables in clause
		if f.Vars[uint(abs(variable))] == True {
			// if any variable is true, clause is true
			return True
		}
	}
	for variable := range clause {
		if f.Vars[uint(abs(variable))] == Unassigned {
			// returns Unassigned if there is an unassigned variable and not alr True
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
		if f.Vars[uint(abs(literal))] == False {
			numFalses += 1
		} else if f.Vars[uint(abs(literal))] == Unassigned {
			numUnassigned += 1
			litUnassigned = literal
		}
	}
	if numFalses == len(clause)-1 && numUnassigned == 1 {
		return true, litUnassigned
	}
	return false, 0
}

func updateImplicationGraph(f *SATInstance, varToAssign uint, clause map[int]bool) {
	impNode := f.ImplicationGraph[varToAssign]
	impNode.Value = f.Vars[varToAssign]
	impNode.Level = f.Level

	if clause != nil {
		// if clause is not nil, then it is an implication not a branched decision
		for literal := range clause {
			connection := f.ImplicationGraph[uint(abs(literal))]
			impNode.Parents[&connection] = true  // add connection to parents
			connection.Children[&impNode] = true // bidirection add to children
		}
		impNode.Clause = clause

	}

	f.ImplicationGraph[varToAssign] = impNode
}

func analyzeConflict(f *SATInstance, conflictClause map[int]bool) (int, map[int]bool, error) {
	if f.Level == 0 {
		// if conflict at level 0, then UNSAT
		return -1, nil, nil
	}
	history := make([]uint, 1)
	history[0] = f.BranchingHist[f.Level]
	history = append(history, f.PropagateHist[f.Level]...)
	log.Println("History for level ", f.Level, history)
	poolLiterals := conflictClause
	finishedLiterals := make(map[int]bool)
	currLevelLiterals := make(map[int]bool)
	prevLevelLiterals := make(map[int]bool)

	for {
		for literal := range poolLiterals {
			if f.ImplicationGraph[uint(abs(literal))].Level == f.Level {
				currLevelLiterals[literal] = true
				// if literal set at current branch, add to current level literals
			} else {
				prevLevelLiterals[literal] = true
				// if literal was set at a previous branch, add to previous level literals
			}
		}
		if len(currLevelLiterals) == 1 {
			// WHY IS THIS 1
			// if one literal is at the current level, then we are done
			break
		}

		lastAssigned, others, err := findLastAssigned(history, poolLiterals)
		if err != nil {
			return -1, nil, err
		}
		finishedLiterals[abs(lastAssigned)] = true // done processing this literal
		currLevelLiterals = others                 // rest of the literals

		poolClause := f.ImplicationGraph[uint(abs(lastAssigned))].Clause
		poolLiterals = make(map[int]bool)

		for literal := range poolClause {
			if _, found := finishedLiterals[abs(literal)]; found {
				continue
			}

			poolLiterals[literal] = true
		}

	}
	learnedClause := make(map[int]bool)
	for literal := range currLevelLiterals {
		learnedClause[literal] = true
	}
	for literal := range prevLevelLiterals {
		learnedClause[literal] = true
	}

	level := 0
	if len(prevLevelLiterals) != 0 {
		// if there are literals in the previous level, then the level is the max of the previous level
		for literal := range prevLevelLiterals {
			currlitLevel := f.ImplicationGraph[uint(abs(literal))].Level
			if currlitLevel > level {
				level = currlitLevel
			}
		}
	} else {
		// if there are no literals in the previous level, then the level is one less than the current level
		level = f.Level - 1
	}
	return level, learnedClause, nil
}

func findLastAssigned(history []uint, clause map[int]bool) (int, map[int]bool, error) {
	v := 0

	sort.Slice(history, func(i, j int) bool { return history[i] > history[j] })
	// reverses history

	for _, varCurr := range history {
		// iterate backwards through history to find last assigned var in clause
		others := make(map[int]bool) // others in clause

		for literal := range clause {
			if uint(abs(literal)) == varCurr {
				v = literal
				continue
			}
			others[literal] = true
		}
		if v != 0 {
			return v, others, nil
		}
	}
	return 0, nil, errors.New("no last assigned var found")
}

func backtrack(f *SATInstance, level int) {
	for currVar, node := range f.ImplicationGraph {
		if node.Level > level {
			node.Value = Unassigned
			node.Level = -1
			node.Parents = make(map[*ImplicationNode]bool)
			node.Children = make(map[*ImplicationNode]bool)
			node.Clause = make(map[int]bool)
			f.ImplicationGraph[currVar] = node
			f.Vars[currVar] = Unassigned
		} else {
			nodeNewChildren := make(map[*ImplicationNode]bool)
			for child := range node.Children {
				if child.Level > level {
					continue
				}
				nodeNewChildren[child] = true
			}
			node.Children = nodeNewChildren
			f.ImplicationGraph[currVar] = node
		}
	}

	remainingBranchingVars := make(map[uint]bool)
	for currVar, assignment := range f.Vars {
		if assignment != Unassigned && len(f.ImplicationGraph[currVar].Parents) == 0 {
			remainingBranchingVars[currVar] = true
		}
	}
	f.BranchingVars = remainingBranchingVars
	levelsHist := make([]int, 0)
	for propLevel := range f.PropagateHist {
		levelsHist = append(levelsHist, propLevel)
	}
	for _, levelCurr := range levelsHist {
		if levelCurr > level {
			delete(f.PropagateHist, levelCurr)
			delete(f.BranchingHist, levelCurr)
		}
	}
}

func unitPropagate(f *SATInstance) (map[int]bool, error) {
	for {
		propQueue := make([]PropStruct, 0)
		for _, clause := range f.Clauses {
			clauseVal := determineClause(f, clause)
			if clauseVal == True {
				continue
			} else if clauseVal == False {
				// returns conflict clause
				return clause, nil
			}
			isUnit, unitLit := isUnitClause(f, clause)
			if isUnit {
				continue
			}
			propStruct := PropStruct{
				Literal: unitLit,
				Clause:  clause,
			}
			propQueue = append(propQueue, propStruct)
		}
		if len(propQueue) == 0 {
			// nothing to propogate/force so returns nil
			return nil, nil
		}
		for _, propStruct := range propQueue {
			propVar := uint(abs(propStruct.Literal))
			if f.Vars[propVar] == Unassigned {
				// assigns to forced value based on unit propogation
				if propStruct.Literal > 0 {
					f.Vars[propVar] = True
				} else {
					f.Vars[propVar] = False
				}
				updateImplicationGraph(f, propVar, propStruct.Clause)

				_, found := f.PropagateHist[f.Level]
				if found {
					f.PropagateHist[f.Level] = append(f.PropagateHist[f.Level], propVar)
				} else {
					// otherwise conflict clause is on the 0th level so should be UNSAT
					log.Println("PropagateHist Key Access at ", f.Level, " not found as expected")
				}
			} else {
				return nil, errors.New("propogating a variable that is already assigned")
			}
		}
	}
}

func SplittingRule(f *SATInstance) (uint, bool) {
	// TODO: Change to make return a variable, NOT A LITERAL

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
				return uint(abs(variable)), true
			}
		}
	}
	switch CountFunc {
	case DLCS, DLIS:
		return uint(abs(keys[0])), f.VarCount[keys[0]].PosCount > f.VarCount[keys[0]].NegCount // explore true or false
	case RDLCS, RDLIS: // uniform at random first 5
		validLiterals := 0
		for i := 0; i < 5; i++ { // messed up if varcount less than 5 but like
			iCounts := f.VarCount[keys[i]]
			if iCounts.NegCount+iCounts.PosCount > 0 {
				validLiterals += 1
			}
		}
		keyToReturn := keys[rand.Intn(validLiterals)]
		return uint(abs(keyToReturn)), f.VarCount[keyToReturn].PosCount > f.VarCount[keyToReturn].NegCount
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

// func getUnitClauses(f *SATInstance) []map[int]bool {
// 	unitClauses := make([]map[int]bool, 0)
// 	for _, clause := range f.Clauses {
// 		isUnit, _ := isUnitClause(f, clause)
// 		if isUnit {
// 			unitClauses = append(unitClauses, clause)
// 		}
// 	}
// 	return unitClauses
// }
//
//
// func determineAllClauses(f *SATInstance) int {
// 	for _, clause := range f.Clauses {
// 		if determineClause(f, clause) == False {
// 			return False
// 			// if any clause is false, formula is false
// 		}
// 	}
// 	for _, clause := range f.Clauses {
// 		if determineClause(f, clause) == Unassigned {
// 			return Unassigned
// 			// returns Unassigned if there is an unassigned variable and not alr False
// 		}
// 	}
// 	// if all clauses are true, formula is true
// 	return True
// }
