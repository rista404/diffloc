package main

import (
	"bufio"
	"bytes"
	"fmt"
	. "github.com/logrusorgru/aurora"
	"golang.org/x/text/message"
	"io"
	"os"
	"os/exec"
	"strconv"
)

// counts number of lines for a reader
func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func getDiffForNewFiles() int {
	// sum of all new files
	sum := 0

	// command for listing out new file names
	cmd := exec.Command("git", "ls-files", "--others", "--exclude-standard")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()

	// scanner for untracked file names
	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		// get the file name
		fn := scanner.Text()

		// open the file
		f, err := os.Open(fn)
		if err != nil {
			fmt.Printf("Cant't open file %s: %s", fn, err)
			os.Exit(2)
		}

		// read number of lines
		add, err := lineCounter(f)
		if err != nil {
			fmt.Printf("Could not read new files diff: %s\n", err)
			os.Exit(2)
		}

		// add to the sum
		sum += add

	}

	return sum
}

// getStatForDiff parses git diff output and combines the loc diff for every file
func getStatForDiff() (int, int) {
	diffArgs := []string{"diff", "--numstat"}
	if len(os.Args) > 2 {
		diffArgs = append(diffArgs, os.Args[2])
	}

	cmd := exec.Command("git", diffArgs...)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Could not run git diff: %s", err)
		os.Exit(2)
	}

	scanner := bufio.NewScanner(&out)

	add := 0
	rm := 0

	// go through every line of changed files
	for scanner.Scan() {
		flds := bytes.Fields(scanner.Bytes())

		// first field is the loc added
		currAdd, err := strconv.Atoi(string(flds[0]))
		if err != nil {
			fmt.Println("Problem with git diff output")
			os.Exit(2)
		}
		// add to counter
		add += currAdd

		// second field is the loc removed
		currRm, err := strconv.Atoi(string(flds[1]))
		if err != nil {
			fmt.Println("Problem with git diff output")
			os.Exit(2)
		}
		// add to counter
		rm += currRm
	}

	return add, rm
}

func main() {
	addDiff, rm := getStatForDiff()
	addNew := getDiffForNewFiles()

	add := addDiff + addNew

	// make an english format printer so we format the thousands
	p := message.NewPrinter(message.MatchLanguage("en"))
	addStr := Bold(Green("+" + p.Sprint(add)))
	rmStr := Bold(Red("-" + p.Sprint(rm)))
	// print the diff in github style
	fmt.Printf("%s %s\n", addStr, rmStr)
}
