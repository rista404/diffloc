package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	var diffArg string
	if len(os.Args) > 2 {
		diffArg = os.Args[2]
	}

	cmd := exec.Command("git", "diff", "--numstat", diffArg)

	var out bytes.Buffer
	cmd.Stdout = &out

	fmt.Println("Output from cmd:")
	fmt.Println(out.String())
}
