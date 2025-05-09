package commands

import (
	"errors"

	"github.com/yashranjan1/gator/internal/command"
	"github.com/yashranjan1/gator/internal/state"
)

type Commands struct {
	Callback map[string]func(*state.State, command.Command) error
}

func (c *Commands) Run(s *state.State, cmd command.Command) error {
	fn, exists := c.Callback[cmd.Name]
	if !exists {
		return errors.New("Error: Command does not exist")
	}

	err := fn(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *Commands) Register(name string, fn func(*state.State, command.Command) error) {
	c.Callback[name] = fn
}
