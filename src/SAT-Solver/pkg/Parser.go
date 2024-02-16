package pkg

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ParseCNFFile(fileName string) (*SATInstance, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var satInstance *SATInstance

	numClauses := 0

	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Fields(line)

		// Skip comment lines
		if len(tokens) > 0 && tokens[0] == "c" {
			continue
		}

		// Problem line
		if len(tokens) > 0 && tokens[0] == "p" {
			if len(tokens) < 4 || tokens[1] != "cnf" {
				return nil, fmt.Errorf("invalid DIMACS file format")
			}

			numVars, err := strconv.Atoi(tokens[2])
			if err != nil {
				return nil, fmt.Errorf("invalid number of variables: %w", err)
			}

			numClauses, err = strconv.Atoi(tokens[3])
			if err != nil {
				return nil, fmt.Errorf("invalid number of clauses: %w", err)
			}

			satInstance = NewSATInstanceVars(numVars)
			continue
		}

		// Clause line
		clause := make(map[int]bool)
		for _, token := range tokens {
			literal, err := strconv.Atoi(token)
			if err != nil {
				return nil, fmt.Errorf("invalid literal: %w", err)
			}
			if literal == 0 {
				break
			}
			clause[literal] = true
			satInstance.addVariable(literal)
		}
		if len(clause) == 0 {
			continue
		}
		satInstance.AddClause(clause)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning file: %w", err)
	}

	if numClauses != satInstance.NumClauses {
		panic("num clauses in parsing doesn't match")
	}
	return satInstance, nil
}
