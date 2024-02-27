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
	isSAT, err := pkg.DPLL(instance)
	if err != nil {
		log.Panic(err.Error())
	}
	duration := time.Since(start)

	if isSAT {
		// TODO: write to a file
		fmt.Println("final output", instance)
		// fmt.Println("final truth values", newInstance.Vars)
		// TODO: will probably have to remove 0 from map
		fmt.Printf("{\"Instance\": \"%s\", \"Time\": %.2f, \"Result\": \"SAT\", \"Solution\": \"%v\"}\n", filename, duration.Seconds(), mapToString(instance.Vars))
	} else {
		fmt.Printf("{\"Instance\": \"%s\", \"Time\": %.2f, \"Result\": \"UNSAT\"}\n", filename, duration.Seconds())
	}
}

func mapToString(vars map[uint]int) string {
	keys := pkg.SortedKeysUnsigned(vars)
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
	}
	return result[:len(result)-1]
}
