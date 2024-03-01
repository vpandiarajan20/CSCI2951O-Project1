package pkg

import (
	"errors"
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
	VSIDS
)

var CountFunc = 4

func CDCL(f *SATInstance) (bool, error) {
	// isSuccessful, err := preprocessFormula(f)
	// if !isSuccessful {
	// 	return false, err
	// }
	// fmt.Println("post pure literl", f)

	for !allVariablesAssigned(f) {
		conflictClause, err := unitPropagate(f)

		if err != nil {
			return false, err
		}

		if conflictClause != nil {
			f.NumConflicts += 1
			log.Println("Learned Clauses", f.LearnedClauses)
			log.Println("Conflict Clause", conflictClause)
			log.Println("Implication Graph")
			for _, i := range f.ImplicationGraph {
				fmt.Println(i.String())
			}
			level, learnedClause, err := analyzeConflict(f, conflictClause)
			f.LearnedClauses = append(f.LearnedClauses, learnedClause)
			if err != nil {
				return false, err
			}
			if level == -1 {
				return false, nil
				// UNSAT!!
			}
			f.AddClause(learnedClause)
			condition := len(f.LearnedClauses) % 20
			if condition == 0 {
				f.DivideVarCounts()
			}

			// add learned clause to formula
			fmt.Println("backtracking to ", level)
			backtrack(f, level)
			isLearnedUnitClause, _ := isUnitClause(f, learnedClause)
			if !isLearnedUnitClause {
				log.Panicln("backtracked too far / created a learned clause that is not going to be acted upon by unitProp")
			}
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
			fmt.Printf("--------------------------------------------Decision Level %d--------------------- \n", f.Level)
			fmt.Println("Branch, Var:", varToAssign, "Assigned", f.Vars[varToAssign] == True)

			f.BranchingHist[f.Level] = varToAssign
			f.PropagateHist[f.Level] = make([]uint, 0)
			updateImplicationGraph(f, varToAssign, nil)
		}
	}
	// fmt.Println("branching history", f.BranchingHist)
	// fmt.Println("propagate history", f.PropagateHist)
	// fmt.Println("implication graph")
	// for _, i := range f.ImplicationGraph {
	// 	fmt.Println(i.String())
	// }
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
	for variable := range clause {
		if f.Vars[uint(abs(variable))] == Unassigned {
			// returns Unassigned if there is an unassigned variable
			return Unassigned
		}
	}

	// TODO: why is this ordered the way it is?

	for variable := range clause { // iterate through all variables in clause
		if (variable > 0 && f.Vars[uint(abs(variable))] == True) || (variable < 0 && f.Vars[uint(abs(variable))] == False) {
			// if any variable is true and shows up as pos || any var is false and shows up as neg, clause is true
			return True
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
			if abs(literal) != int(varToAssign) {
				impNode.Parents[uint(abs(literal))] = true // add connection to parents
				// connection.Children[impNode.Var] = true    // bidirection add to children
			}
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

	learntClause := make(map[int]bool)
	maxLevel := 0
	currLiteral := 0

	history := make([]uint, 1)
	history[0] = f.BranchingHist[f.Level]
	history = append(history, f.PropagateHist[f.Level]...)

	historyOld := make([]uint, len(history))
	copy(historyOld, history)

	// ReverseUIntSlice(history)
	// ASSERTION: history is reversed
	// for i, varCurr := range history {
	// 	if varCurr != historyOld[len(historyOld)-i-1] {
	// 		log.Panic("history not reversed correctly!")
	// 	}
	// }
	literalsToProcess := make(map[int]bool)

	for litCurr := range conflictClause {
		if f.ImplicationGraph[uint(abs(litCurr))].Level == f.Level {
			literalsToProcess[litCurr] = true
		} else {
			toAdd := litCurr
			if f.ImplicationGraph[uint(abs(litCurr))].Value == True {
				toAdd = -abs(litCurr)
			} else if f.ImplicationGraph[uint(abs(litCurr))].Value == False {
				toAdd = abs(litCurr)
			} else {
				log.Panic("Parent is unassigned")
			}
			learntClause[toAdd] = true
			fmt.Println("Adding to Learned Clause", toAdd)
		}
	}

	trailStack := NewStackUint(history)

	fmt.Println("Analyzing this clause", conflictClause, "at level:", f.Level)

	// varsToProcess := conflictClause
	// literalsToProcess := Queue{}
	// for literal := range conflictClause {
	// 	literalsToProcess.Enqueue(literal)
	// }
	varsProcessed := make(map[int]bool)
	for len(literalsToProcess) > 1 {
		fmt.Println("literalSet:", literalsToProcess)
		// 
		fmt.Println("trailStack:", trailStack)
		pVar, isSuccessful := trailStack.Pop()
		if !isSuccessful {
			return 0, nil, errors.New("trail stack is empty")
		}
		if f.Vars[pVar] == True {
			currLiteral = int(pVar)
		} else if f.Vars[pVar] == False {
			currLiteral = -int(pVar)
		} else {
			log.Panic("Parent is unassigned")
		}

		_, isFoundP := conflictClause[currLiteral]
		if isFoundP {
			// should never happen
			delete(literalsToProcess, currLiteral)
		}
		_, isFoundN := conflictClause[-currLiteral]
		if isFoundN {
			delete(literalsToProcess, -currLiteral)
		}
		fmt.Println("Processing parents of", currLiteral)
		if isFoundP || isFoundN {
			for parentVar := range f.ImplicationGraph[uint(abs(currLiteral))].Parents {
				fmt.Println("deciding on parent:", parentVar, "Level:", f.ImplicationGraph[parentVar].Level)
				if f.ImplicationGraph[parentVar].Level == f.Level {
					literalToEnqueue := int(parentVar)
					if f.ImplicationGraph[parentVar].Value == True {
						// toEnqueue = int(parent)
					} else if f.ImplicationGraph[parentVar].Value == False {
						literalToEnqueue = -int(parentVar)
					} else {
						log.Panic("Parent is unassigned")
					}
					_, alreadyProcessed := varsProcessed[abs(literalToEnqueue)]
					if alreadyProcessed {
						continue
					}
					fmt.Println("Adding to Set", literalToEnqueue)
					// literalsToProcess.Enqueue(literalToEnqueue)
					literalsToProcess[literalToEnqueue] = true
					varsProcessed[abs(literalToEnqueue)] = true
				} else {
					// if parent is not at current level, then it is at previous level
					if f.ImplicationGraph[parentVar].Value == True {
						fmt.Println("Adding to Learned Clause", -int(parentVar))
						learntClause[-int(parentVar)] = true
						// flip sign
					} else if f.ImplicationGraph[parentVar].Value == False {
						fmt.Println("Adding to Learned Clause", int(parentVar))
						learntClause[int(parentVar)] = true
						// flip sign
					} else {
						log.Panic("Parent is unassigned")
					}
					maxLevel = max(maxLevel, f.ImplicationGraph[parentVar].Level)
				}
				// be careful bc parents are always positive + varsToProcess isn't
			}
		}
	}

	// need to pop off the decision stack
	// once you pop off the decision stack

	// is it okay for a literal and its negation to be in literalsToProcess?
	// we need to make sure last element has not been visited yet, hence I'm marking things visited as they're put in Q
	// fmt.Println("literalQ:", literalsToProcess)
	// currLiteral, err := literalsToProcess.Dequeue()

	// if err != nil {
	// 	log.Panicln(err.Error())
	// }
	for literal := range literalsToProcess {
		currLiteral = literal
	}
	currVar := abs(currLiteral)
	if f.ImplicationGraph[uint(currVar)].Value == True {
		fmt.Println("Adding to Learned Clause", -int(currVar))
		learntClause[-int(currVar)] = true
		// flip sign
	} else if f.ImplicationGraph[uint(currVar)].Value == False {
		fmt.Println("Adding to Learned Clause", int(currVar))
		learntClause[int(currVar)] = true
		// flip sign
	} else {
		log.Panic("Parent is unassigned")
	}

	// TODO: optimization, but could be commented out for now
	// if len(learntClause) == 1 {
	// 	return 0, learntClause, nil
	// }
	fmt.Println("Learned Clause", learntClause, "at level", maxLevel)
	return maxLevel, learntClause, nil
}

func backtrack(f *SATInstance, level int) {

	for i := f.Level; i > level; i-- {
		for _, currVar := range f.PropagateHist[i] {
			node := f.ImplicationGraph[currVar]
			if node.Value != Unassigned && node.Value != f.Vars[currVar] {
				log.Panic("in stack (propagate hist) backtracking and unassigning something that was never assigned")
			}
			node.Value = Unassigned
			node.Level = -1
			node.Parents = make(map[uint]bool)
			node.Clause = make(map[int]bool)
			f.ImplicationGraph[currVar] = node
			f.Vars[currVar] = Unassigned
		}
		currVar := f.BranchingHist[i]
		node := f.ImplicationGraph[currVar]
		if node.Value != Unassigned && node.Value != f.Vars[currVar] {
			log.Panic("in stack (branching hist) backtracking and unassigning something that was never assigned")
		}
		node.Value = Unassigned
		node.Level = -1
		node.Parents = make(map[uint]bool)
		node.Clause = make(map[int]bool)
		f.ImplicationGraph[currVar] = node
		f.Vars[currVar] = Unassigned
	}

	for currVar, node := range f.ImplicationGraph {
		if node.Level > level {
			node.Value = Unassigned
			node.Level = -1
			node.Parents = make(map[uint]bool)
			// node.Children = make(map[uint]bool)
			node.Clause = make(map[int]bool)
			f.ImplicationGraph[currVar] = node
			f.Vars[currVar] = Unassigned
		} else {
			f.ImplicationGraph[currVar] = node
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
			if !isUnit {
				continue
			}
			propStruct := PropStruct{
				Literal: unitLit,
				Clause:  clause,
			}
			propQueue = append(propQueue, propStruct)
		}

		sort.Slice(propQueue, func(i, j int) bool { return propQueue[i].Literal < propQueue[j].Literal })
		if len(propQueue) == 0 {
			// nothing to propogate/force so returns nil
			return nil, nil
		}
		for _, propStruct := range propQueue {
			propVar := uint(abs(propStruct.Literal))

			// assigns to forced value based on unit propogation
			if propStruct.Literal > 0 {
				f.Vars[propVar] = True
				// fmt.Println("setting", propVar, "to true")
			} else {
				f.Vars[propVar] = False
				// fmt.Println("setting", propVar, "to false")
			}
			updateImplicationGraph(f, propVar, propStruct.Clause)
			_, found := f.PropagateHist[f.Level]
			if found {
				f.PropagateHist[f.Level] = append(f.PropagateHist[f.Level], propVar)
			} else {
				// otherwise conflict clause is on the 0th level so should be UNSAT
				// log.Println("PropagateHist Key Access at ", f.Level, " not found as expected")
			}

			// no else because need to prop twice to see that conflict

			// } else {
			// 	return nil, errors.New(fmt.Sprint("propagating a variable that is already assigned:", propVar))
			// }
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

	// if very few clauses, then use DLIS instead of VSIDS
	CountFuncMod := CountFunc
	if CountFuncMod == 4 && len(f.LearnedClauses) < 4 {
		CountFuncMod = 1
	}

	switch CountFuncMod {
	case DLCS, RDLCS:
		sort.Slice(keys, func(i, j int) bool {
			iCounts := f.VarCount[keys[i]]
			jCounts := f.VarCount[keys[j]]
			if (iCounts.NegCount + iCounts.PosCount) > (jCounts.NegCount + jCounts.PosCount) {
				return true
			}
			return keys[i] < keys[j]
			// counts stored in struct with NegCount and PosCount
		})
	case DLIS, RDLIS, VSIDS:
		sort.Slice(keys, func(i, j int) bool {
			iCounts := f.VarCount[keys[i]]
			jCounts := f.VarCount[keys[j]]
			if max(iCounts.NegCount, iCounts.PosCount) > max(jCounts.NegCount, jCounts.PosCount) {
				return true
			}
			return keys[i] < keys[j]
		})
		// fmt.Println("sorted keys", keys)
	default:
		for _, clause := range f.Clauses {
			for variable := range clause {
				return uint(abs(variable)), true
			}
		}
	}
	switch CountFunc {
	case DLCS, DLIS, VSIDS:
		for _, k := range keys {
			if f.Vars[uint(abs(k))] == Unassigned {
				return uint(abs(k)), f.VarCount[k].PosCount > f.VarCount[k].NegCount
			}
		}
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

// history := make([]uint, 1)
// history[0] = f.BranchingHist[f.Level]
// history = append(history, f.PropagateHist[f.Level]...)

// historyOld := make([]uint, len(history))
// copy(historyOld, history)

// slices.Reverse(history)

// ASSERTION: history is reversed
// for i, varCurr := range history {
// 	if varCurr != historyOld[len(historyOld)-i-1] {
// 		log.Panic("history not reversed correctly!")
// 	}
// }

// log.Println("History for level ", f.Level, history)
