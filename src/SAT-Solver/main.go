package main

import (
	"SAT-Solver/pkg"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	// log.SetOutput(ioutil.Discard)
	if len(os.Args) != 2 && len(os.Args) != 3 {
		log.Fatalf("Usage: ./solver <file name> or Usage: ./solver <file name> <heuristic>")
	}
	inputFile := os.Args[1]
	filename := filepath.Base(inputFile)
	if len(os.Args) == 3 {
		pkg.CountFunc, _ = strconv.Atoi(os.Args[2])
	}

	instance, err := pkg.ParseCNFFile(inputFile)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("initial equations", instance)
	start := time.Now()
	newInstance, isSAT := pkg.DPLL(instance)
	duration := time.Since(start)

	if isSAT {
		// TODO: write to a file
		fmt.Println("final output", newInstance)
		// fmt.Println("final truth values", newInstance.Vars)
		fmt.Printf("{\"Instance\": \"%s\", \"Time\": %.2f, \"Result\": \"SAT\", \"Solution\": \"%v\"}\n", filename, duration.Seconds(), mapToString(newInstance.Vars))
	} else {
		fmt.Printf("{\"Instance\": \"%s\", \"Time\": %.2f, \"Result\": \"UNSAT\"}\n", filename, duration.Seconds())
	}
}

func mapToString(vars map[int]int) string {
	keys := pkg.SortedKeys(vars)
	var result string
	for _, key := range keys {
		switch vars[key] {
		case pkg.True:
			result += fmt.Sprintf("%d true ", key)
		case pkg.False:
			result += fmt.Sprintf("%d false ", key)
		case pkg.Unassigned:
			// TODO: maybe throw an error
			fmt.Printf("%d is left unassigned \n", key)
			result += fmt.Sprintf("%d false ", key)
		}
		if vars[key] == pkg.True {
		}
	}
	return result[:len(result)-1]
}
