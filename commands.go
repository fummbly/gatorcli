package main

import (
	"errors"
	"fmt"
)

type input struct {
	Name string
	Args []string
}

type command struct {
	handler func(*state, input) error
	descr   string
}

type commands struct {
	registeredCommands map[string]command
}

func (c *commands) register(name string, cmd command) {
	c.registeredCommands[name] = cmd
}

func (c *commands) help() {
	for cmdName, cmd := range c.registeredCommands {
		fmt.Printf("%s: %s\n", cmdName, cmd.descr)
	}
}

func (c *commands) run(s *state, in input) error {
	cmd, ok := c.registeredCommands[in.Name]
	if !ok {
		return errors.New("command not found")
	}
	return cmd.handler(s, in)
}
