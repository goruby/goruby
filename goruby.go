package main

import (
	"flag"
	"fmt"
	"github.com/goruby/goruby/lexer"
	"github.com/goruby/goruby/parser"
)

func main() {
	var (
		eFlag   = flag.String("e", "", "one line of script. Several -e's allowed. Omit [programfile]")
		verbose = flag.Bool("v", false, "print version number, then turn on verbose mode")
	)

	flag.Parse()

	if *verbose {
		fmt.Printf("Verbose: %t\n", *verbose)
	}
	// cmd := exec.Command("ruby", "-e "+*eFlag, verboseOpt)
	// output, err := cmd.CombinedOutput()
	// if len(output) != 0 {
	// 	fmt.Print(string(output))
	// }
	// if err != nil {
	// 	fmt.Print(err.Error())
	// }
	l := lexer.NewLexer([]byte(*eFlag))
	p := parser.NewParser()
	rslt, err := p.Parse(l)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Result: %+#v", rslt)
}
