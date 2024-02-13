package main

import (
	"SAT-Solver/pkg"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// log.SetOutput(ioutil.Discard)
	if len(os.Args) != 2 {
		log.Fatalf("Usage: ./solver <file name>")
	}
	inputFile := os.Args[1]
	filename := filepath.Base(inputFile)

	start := time.Now()

	instance, err := pkg.ParseCNFFile(inputFile)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// newClause := make(map[int]bool, 0)
	// newClause[3] = false
	// newClause[-1] = false
	// instance.AddClause(newClause)
	// newClause = make(map[int]bool, 0)
	// newClause[1] = false
	// newClause[-2] = false
	// newClause[-3] = false
	// instance.AddClause(newClause)
	fmt.Println("initial equations", instance)
	newInstance, isSAT := pkg.DPLL(instance)
	duration := time.Since(start)

	if isSAT {
		fmt.Println("final output", newInstance)
		fmt.Println("final truth values", newInstance.Vars)
		fmt.Printf("{\"Instance\": \"%s\", \"Time\": %.2f, \"Result\": \"SAT\"}\n", filename, duration.Seconds())
	} else {
		fmt.Printf("{\"Instance\": \"%s\", \"Time\": %.2f, \"Result\": \"UNSAT\"}\n", filename, duration.Seconds())
	}
}
