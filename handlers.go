package main

import (
	"errors"
	"fmt"

	"github.com/yashranjan1/gator/internal/command"
	"github.com/yashranjan1/gator/internal/state"
)

func handlerLogin(s *state.State, cmd command.Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("the login handler expects a single argument, the username")
	}

	err := s.Config.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Println("User has been set")
	return nil
}
