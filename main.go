package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/goruby/goruby/interpreter"
	"github.com/pkg/errors"
)

type multiString []string

func (m multiString) String() string {
	out := ""
	for _, s := range m {
		out += s
	}
	return out
}
func (m *multiString) Set(s string) error {
	*m = append(*m, s)
	return nil
}

var onelineScripts multiString

func main() {
	flag.Var(&onelineScripts, "e", "one line of script. Several -e's allowed. Omit [programfile]")
	flag.Parse()
	interpreter := interpreter.New()
	if len(onelineScripts) != 0 {
		input := strings.Join(onelineScripts, "\n")
		_, err := interpreter.Interpret("", input)
		if err != nil {
			fmt.Printf("%v\n", errors.Cause(err))
			os.Exit(1)
		}
		return
	}
	args := flag.Args()
	if len(args) == 0 {
		log.Println("No program files specified")
		os.Exit(1)
	}
	fileBytes, err := ioutil.ReadFile(args[0])
	if err != nil {
		log.Printf("Error while opening program file: %T:%v\n", err, err)
		os.Exit(1)
	}
	_, err = interpreter.Interpret(args[0], fileBytes)
	if err != nil {
		fmt.Printf("%v\n", errors.Cause(err))
		os.Exit(1)
	}
	return
}
