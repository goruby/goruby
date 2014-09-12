package main

import (
	"flag"
	"fmt"
	"os/exec"
)

func main() {
	var (
		eFlag   = flag.String("e", "", "one line of script. Several -e's allowed. Omit [programfile]")
		verbose = flag.Bool("v", false, "print version number, then turn on verbose mode")
	)

	flag.Parse()

	verboseOpt := ""
	if *verbose {
		//verboseOpt = "-v"
	}
	cmd := exec.Command("ruby", "-e "+*eFlag, verboseOpt)
	output, err := cmd.CombinedOutput()
	if len(output) != 0 {
		fmt.Print(string(output))
	}
	if err != nil {
		fmt.Print(err.Error())
	}
}
