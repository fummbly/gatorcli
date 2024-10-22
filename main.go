package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fummbly/gatorcli/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	config, err := config.Read()
	if err != nil {
		log.Fatal("error reading config: %v", err)
	}

	programState := &state{
		cfg: &config,
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("Usage cli <command> [args...]")
		return
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.run(programState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}

}
