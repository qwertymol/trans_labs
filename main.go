package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

const emptyString = "~"

var inputFile = flag.String("i", "gram.txt", "input file (grammatics")
var inputSequence = flag.String("s", "S", "sequence for input")
var maxDepth = flag.Int("d", 4, "max depth (-1 means no limit)")
var maxLength = flag.Int("max", -1, "max length of output sequence (-1 means no limit)")
var minLength = flag.Int("min", -1, "min length of output sequence (-1 means no limit)")
var verbose = flag.Bool("v", false, "verbose on/of")
var searchOutput = flag.Bool("o", false, "search path (instead of find all outputs from -s)")

type Transition struct {
	from, to string
}

func parseLine(line string) []*Transition {
	buf := strings.Split(line, "->")
	if len(buf) != 2 {
		return nil
	}

	var res []*Transition

	for _, val := range strings.Split(buf[1], "|") {
		if val == emptyString {
			val = ""
		}

		res = append(res, &Transition{from: buf[0], to: val})
	}

	return res
}

func parseFile(fileName string) ([]*Transition, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	var res []*Transition
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()

		// comment
		if strings.HasPrefix(text, "//") {
			continue
		}

		tr := parseLine(text)
		if tr == nil {
			continue
		}

		res = append(res, tr...)
	}

	return res, nil
}

func checkNonTerminal(seq string) bool {
	return strings.ToLower(seq) != seq
}

func checkTerminal(seq string) bool {
	return strings.ToUpper(seq) != seq
}

func checkStart(seq string) bool {
	return seq != "S"
}

func stringIndexFrom(s, sub string, from int) int {
	if from > len(s) {
		from = len(s)
	}

	ind := strings.Index(s[from:], sub)
	if ind != -1 {
		ind += from
	}

	return ind
}

func replaceFrom(s, old, new string, from int) string {
	return s[:from] + strings.Replace(s[from:], old, new, 1)
}

func addBrackets(s string, from, to int) string {
	return s[:from] + "[" + s[from:to] + "]" + s[to:]
}

func checkMinLen(seq string) bool {
	if *minLength == -1 {
		return true
	}

	return len(seq) >= *minLength
}

func checkMaxLen(seq string) bool {
	if *maxLength == -1 {
		return true
	}

	return len(seq) <= *maxLength
}

func checkDepth(depth int) bool {
	if *maxDepth == -1 {
		return true
	}

	return depth <= *maxDepth
}

func findSequences(transitions *[]*Transition, sequence string, depth int, path []string) {
	if !checkNonTerminal(sequence) {
		if !checkMaxLen(sequence) || !checkMinLen(sequence) {
			return
		}

		if *verbose {
			fmt.Print(strings.Join(path, ",") + " ")
		}

		fmt.Println("'" + sequence + "'")
		return
	}

	if depth == 0 {
		return
	}

	// проверяем все переходы
	for _, tr := range *transitions {
		last := 0
		for last = stringIndexFrom(sequence, tr.from, last); last != -1; last = stringIndexFrom(sequence, tr.from, last) {
			newSeq := replaceFrom(sequence, tr.from, tr.to, last)
			newPath := fmt.Sprintf("('%s'=>'%s')", addBrackets(sequence, last, last+len(tr.from)), tr.to)
			findSequences(transitions, newSeq, depth-1, append(path, newPath))
			last++
		}
	}
}

func findPath(transitions *[]*Transition, sequence string, depth int, path []string) {
	if !checkStart(sequence) {
		if *verbose {
			fmt.Print(strings.Join(path, ",") + " ")
		}

		fmt.Println("'" + *inputSequence + "'")
		return
	}

	if depth == 0 {
		return
	}

	for _, tr := range *transitions {
		last := 0
		for last = stringIndexFrom(sequence, tr.to, last); last != -1; last = stringIndexFrom(sequence, tr.to, last) {
			newSeq := replaceFrom(sequence, tr.to, tr.from, last)
			newPath := fmt.Sprintf("('%s'=>'%s')", addBrackets(newSeq, last, last+len(tr.from)), tr.to)
			findPath(transitions, newSeq, depth-1, append([]string{newPath}, path...))
			last++
		}
	}
}

func init() {
	flag.Parse()

	if *maxDepth == -1 && *maxLength == -1 {
		fmt.Println("-max or -d will be set > -1")
		os.Exit(-1)
	}
}

func main() {
	trs, err := parseFile(*inputFile)
	if err != nil {
		panic(err)
	}

	fmt.Println("Input sequence:", *inputSequence)
	fmt.Println("Transitions:")
	for _, tr := range trs {
		fmt.Printf("  %s->%s\n", tr.from, tr.to)
	}
	fmt.Println()

	fmt.Println("Variants:")
	if *searchOutput {
		findPath(&trs, *inputSequence, *maxDepth, []string{})
	} else {
		findSequences(&trs, *inputSequence, *maxDepth, []string{})
	}
}
