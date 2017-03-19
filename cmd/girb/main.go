package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"os/user"
	"syscall"

	"github.com/goruby/goruby/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Ruby programming language in Go!\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")
	start(os.Stdin, os.Stdout)
}

func start(in io.Reader, out io.Writer) {
	printChan := make(chan string)
	sigChan := make(chan os.Signal, 4)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGHUP, os.Kill)
	go repl.Start(in, printChan)
	for {
		select {
		case evaluated, ok := <-printChan:
			fmt.Fprintf(out, "%s", evaluated)
			if !ok {
				return
			}
		case sig := <-sigChan:
			fmt.Fprintln(out)
			if sig != syscall.SIGINT {
				return
			}
		}
	}
}
