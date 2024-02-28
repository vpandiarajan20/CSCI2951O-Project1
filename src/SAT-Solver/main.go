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
	if len(os.Args) != 2 && len(os.Args) != 3 {
		log.Fatalf("Usage: ./solver <file name> or Usage: ./solver <file name> <output path>")
	}
	inputFile := os.Args[1]
	filename := filepath.Base(inputFile)
	outputFile := "output_assignments.txt"
	if len(os.Args) == 3 {
		outputFile = os.Args[2]
	}

	// if len(os.Args) == 3 {
	// 	pkg.CountFunc, _ = strconv.Atoi(os.Args[2])
	// }

	instance, err := pkg.ParseCNFFile(inputFile)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("initial equations", instance)
	start := time.Now()
	isSAT, err := pkg.CDCL(instance)
	duration := time.Since(start)
	if err != nil {
		log.Panic(err.Error())
	}
	file, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if isSAT {
		// TODO: write to a file
		fmt.Fprintf(file, "final truth values %v\n", instance.Vars)
		fmt.Printf("{\"Instance\": \"%s\", \"Time\": %.2f, \"Result\": \"SAT\", \"Solution\": \"%v\"}\n", filename, duration.Seconds(), mapToString(instance.Vars))
	} else {
		fmt.Printf("{\"Instance\": \"%s\", \"Time\": %.2f, \"Result\": \"UNSAT\"}\n", filename, duration.Seconds())
	}
}

func mapToString(vars map[uint]int) string {
	keys := pkg.SortedKeysUint(vars)
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

// func mapToStringOld(vars map[int]bool) string {
// 	keys := pkg.SortedKeys(vars)
// 	var result string
// 	for _, key := range keys {
// 		result += fmt.Sprintf("%d %v ", key, vars[key])
// 	}
// 	return result[:len(result)-1]
// }
