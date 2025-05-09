package main

import (
	"fmt"
	"os"

	"github.com/yashranjan1/gator/internal/command"
	"github.com/yashranjan1/gator/internal/commands"
	"github.com/yashranjan1/gator/internal/config"
	"github.com/yashranjan1/gator/internal/state"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	s := state.State{
		Config: &conf,
	}

	cmds := commands.Commands{
		Callback: make(map[string]func(*state.State, command.Command) error),
	}

	cmds.Register("login", handlerLogin)

	args := os.Args

	if len(args) < 2 {
		fmt.Println("Error: Not enough arguments\n")
		fmt.Println("Usage:\n")
		fmt.Println("gator <COMMAND> [optional]")
		os.Exit(1)
	}

	var cmdArgs []string

	if len(args) > 2 {
		cmdArgs = args[2:]
	}

	cmd := command.Command{
		Name: args[1],
		Args: cmdArgs,
	}

	err = cmds.Run(&s, cmd)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
