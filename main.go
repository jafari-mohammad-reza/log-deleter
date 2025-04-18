package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	query := buildQuery()
	addr, index := os.Getenv("ADDR"), os.Getenv("INDEX")
	if addr == "" || index == "" {
		log.Fatal("insert both addr and index")
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/_delete_by_query", addr, index), bytes.NewBuffer(query))
	if err != nil {
		log.Fatalf("error while deleting %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("error while deleting %s", err.Error())
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("response:", string(body))
}

func buildQuery() []byte {
	input := os.Getenv("QUERY")
	if input == "" {
		log.Fatal("insert query")
	}
	input = strings.ToLower(input)
	input = strings.ReplaceAll(input, "and not", "+ANDNOT+")
	input = strings.ReplaceAll(input, "and", "+AND+")

	clauses := strings.Split(input, "+")
	mustClauses := []map[string]map[string]string{}
	mustNotClauses := []map[string]map[string]string{}

	for i := 0; i < len(clauses); i++ {
		clause := strings.TrimSpace(clauses[i])
		if clause == "" {
			continue
		}

		switch clause {
		case "AND":
			if i+1 < len(clauses) {
				target := strings.TrimSpace(clauses[i+1])
				splited := strings.SplitN(target, ":", 2)
				if len(splited) == 2 {
					key, val := strings.TrimSpace(splited[0]), strings.TrimSpace(splited[1])
					mustClauses = append(mustClauses, map[string]map[string]string{
						"match": {key: val},
					})
				}
				i++
			}
		case "ANDNOT":
			if i+1 < len(clauses) {
				target := strings.TrimSpace(clauses[i+1])
				splited := strings.SplitN(target, ":", 2)
				if len(splited) == 2 {
					key, val := strings.TrimSpace(splited[0]), strings.TrimSpace(splited[1])
					mustNotClauses = append(mustNotClauses, map[string]map[string]string{
						"match": {key: val},
					})
				}
				i++
			}
		default:
			splited := strings.SplitN(clause, ":", 2)
			if len(splited) == 2 {
				key, val := strings.TrimSpace(splited[0]), strings.TrimSpace(splited[1])
				mustClauses = append(mustClauses, map[string]map[string]string{
					"match": {key: val},
				})
			}
		}
	}
	scrollSize := 5000
	slices := 5
	scrollEnv := os.Getenv("SCROLL_SIZE")
	if scrollEnv != "" {
		scrollSize, _ = strconv.Atoi(scrollEnv)
	}
	sliceEnv := os.Getenv("SLICES")
	if sliceEnv != "" {
		slices, _ = strconv.Atoi(sliceEnv)
	}
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must":     mustClauses,
				"must_not": mustNotClauses,
			},
		},
		"conflicts":   "proceed",
		"scroll_size": scrollSize,
		"slices":      slices,
	}

	queryJSON, err := json.MarshalIndent(query, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal query: %v", err)
	}
	fmt.Println(string(queryJSON))
	return queryJSON
}
