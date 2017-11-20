package main

import (
	"log"
	"os"

	"github.com/goruby/goruby/repl"
	"github.com/goruby/readline"
)

func main() {
	exit := startRepl()
	os.Exit(exit)
}

// Readline returns a readline enabled REPL
func startRepl() int {
	// Configure input
	rm := rawMode{StdinFd: int(os.Stdin.Fd())}
	config := &readline.Config{
		InterruptPrompt:   "^C",
		EOFPrompt:         "\n",
		HistorySearchFold: true,
		FuncMakeRaw:       rm.enter,
		FuncExitRaw:       rm.exit,
	}
	l, err := readline.NewEx(config)
	if err != nil {
		log.Printf("Error initializing readlines: %v\n", err)
		return 1
	}
	defer l.Close()
	lNoInterrupt := &ignoreInterrupt{l}

	r := repl.New(lNoInterrupt, lNoInterrupt, lNoInterrupt)
	err = r.Start()
	if err != nil {
		log.Printf("Error within repl: %v\n", err)
		return 1
	}

	return 0
}

type ignoreInterrupt struct {
	*readline.Instance
}

func (i *ignoreInterrupt) Readline() (string, error) {
	line, err := i.Instance.Readline()
	if err == readline.ErrInterrupt {
		return line, nil
	}
	return line, err
}

// rawMode is a helper for entering and exiting raw mode.
type rawMode struct {
	StdinFd int

	state *readline.State
}

// enter is used to put the terminal in raw mode
func (r *rawMode) enter() (err error) {
	r.state, err = readline.MakeRaw(r.StdinFd)
	return err
}

// exit restores the terminal's previous state
func (r *rawMode) exit() error {
	if r.state == nil {
		return nil
	}

	return readline.Restore(r.StdinFd, r.state)
}
