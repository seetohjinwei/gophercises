package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Adapted from: https://stackoverflow.com/a/58841827
func readCsvFile(filePath string) [][]string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("unable to read file %q: %v", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("unable to parse file %q: %v", filePath, err)
	}

	return records
}

func shuffleRows(rows [][]string) {
	rand.Shuffle(len(rows), func(i, j int) {
		rows[i], rows[j] = rows[j], rows[i]
	})
}

// isSameAnswer checks if (got == want) With string trimming and case-insensitivity.
func isSameAnswer(got, want string) bool {
	clean := func(x string) string {
		trimmed := strings.TrimSpace(x)
		lowered := strings.ToLower(trimmed)
		return lowered
	}

	gotCleaned := clean(got)
	wantCleaned := clean(want)

	return gotCleaned == wantCleaned
}

func askQuestion(index int, score *int32, rows [][]string, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	if index >= len(rows) {
		return
	}

	row := rows[index]

	question := row[0]
	answer := row[1]

	fmt.Print(question + " ")
	var input string
	fmt.Scanln(&input)

	if isSameAnswer(input, answer) {
		atomic.AddInt32(score, 1)
	}

	askQuestion(index+1, score, rows, wg)
}

func runQuiz(rows [][]string) (score int32) {
	var wg sync.WaitGroup

	go askQuestion(0, &score, rows, &wg)

	timeout := time.After(time.Duration(duration) * time.Second)
	ch := make(chan struct{})
	go func() {
		wg.Wait()
		close(ch)
	}()

	select {
	case <-timeout:
		fmt.Println("Time's up!")
		return
	case <-ch:
		fmt.Println("Everything is answered!")
		return
	}
}

const (
	defaultCsv      = "problem.csv"
	defaultDuration = 30
)

var (
	filePath    string
	duration    int
	shuffleFlag bool
)

func getFlags() {
	flag.StringVar(&filePath, "file", defaultCsv, "file path questions csv")
	flag.IntVar(&duration, "duration", defaultDuration, "duration of quiz")
	flag.BoolVar(&shuffleFlag, "shuffle", false, "shuffle questions")
	flag.Parse()
}

func main() {
	getFlags()

	rows := readCsvFile(filePath)
	if shuffleFlag {
		shuffleRows(rows)
	}

	score := runQuiz(rows)
	total := len(rows)
	fmt.Printf("You scored %d / %d!\n", score, total)
}
