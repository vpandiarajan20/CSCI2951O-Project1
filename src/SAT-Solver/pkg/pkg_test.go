package pkg

import (
	"fmt"
	"testing"
)

func TestUnitProp(t *testing.T) {
	instance, err := ParseCNFFile("../toy_simple.cnf")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	newClause := make(map[int]bool, 0)
	newClause[3] = true
	instance.Clauses = append(instance.Clauses, newClause)
	fmt.Println(instance)
	UnitPropagate(instance)
	fmt.Println(instance)
	fmt.Println(instance.Vars)
	if len(instance.Clauses) != 0 {
		t.Errorf("Unit Propagation failed")
	}
}
func TestUnitProp2(t *testing.T) {
	instance, err := ParseCNFFile("../toy_simple.cnf")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	newClause := make(map[int]bool, 0)
	newClause[-3] = true
	instance.Clauses = append(instance.Clauses, newClause)
	// fmt.Println(instance)
	UnitPropagate(instance)
	// fmt.Println(instance)
	// fmt.Println(instance.Vars)
	if len(instance.Clauses) != 1 {
		t.Errorf("Unit Propagation failed")
	}
}

func TestUnitPropStop(t *testing.T) {
	instance, err := ParseCNFFile("../toy_simple.cnf")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// fmt.Println(instance)
	UnitPropagate(instance)
	// fmt.Println(newInstance)
	// fmt.Println(TValues)
	if len(instance.Clauses) != 2 {
		t.Errorf("Unit Propagation failed to stop")
	}
}

func TestPureLiteralElim(t *testing.T) {
	instance, err := ParseCNFFile("../toy_simple.cnf")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// fmt.Println(instance)
	PureLiteralElim(instance)
	// fmt.Println(newInstance)
	if len(instance.Clauses) != 0 {
		t.Errorf("Pure Literal Elim failed")
	}
}

func TestPureLiteralElimStop(t *testing.T) {
	instance, err := ParseCNFFile("../toy_simple.cnf")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	newClause := make(map[int]bool, 0)
	newClause[3] = false
	newClause[-1] = false
	instance.Clauses = append(instance.Clauses, newClause)
	// fmt.Println("initial equations", instance)
	PureLiteralElim(instance)
	// fmt.Println("final output", newInstance)
	// fmt.Println("final truth values", TValues)
	if len(instance.Clauses) != 2 {
		t.Errorf("Pure Literal Elim failed")
	}
}

func TestDPLL(t *testing.T) {
	instance, err := ParseCNFFile("../toy_simple.cnf")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	newInstance, isSAT := DPLL(instance)
	if !isSAT {
		t.Errorf("DPLL fail")
	}
	if !newInstance.Vars[1] && newInstance.Vars[3] {
		t.Errorf("DPLL fail")
	}
}

func TestDPLL2(t *testing.T) {
	instance, err := ParseCNFFile("../toy_simple.cnf")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	newClause := make(map[int]bool, 0)
	newClause[3] = false
	newClause[-1] = false
	instance.Clauses = append(instance.Clauses, newClause)

	newInstance, isSAT := DPLL(instance)
	if !isSAT {
		t.Errorf("DPLL fail")
	}
	if !newInstance.Vars[1] && newInstance.Vars[3] {
		t.Errorf("DPLL fail")
	}
}

func TestDPLL3(t *testing.T) {
	instance, err := ParseCNFFile("../toy_lecture.cnf")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	_, isSAT := DPLL(instance)
	if !isSAT {
		t.Errorf("DPLL fail")
	}
}

func TestDPLL4Fail(t *testing.T) {
	instance, err := ParseCNFFile("../toy_infeasible.cnf")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	_, isSAT := DPLL(instance)
	if isSAT {
		t.Errorf("DPLL fail")
	}
}

func TestDPLL4(t *testing.T) {
	instance, err := ParseCNFFile("../toy_solveable.cnf")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	_, isSAT := DPLL(instance)
	if !isSAT {
		t.Errorf("DPLL fail")
	}
}
