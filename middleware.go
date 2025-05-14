package main

import (
	"context"
	"database/sql"
	"errors"

	"github.com/yashranjan1/gator/internal/command"
	"github.com/yashranjan1/gator/internal/database"
	"github.com/yashranjan1/gator/internal/state"
)

func middlewareLoggedIn(handler func(s *state.State, cmd command.Command, user database.User) error) func(*state.State, command.Command) error {
	return func(s *state.State, c command.Command) error {
		username := s.Config.CurrentUser
		if len(username) == 0 {
			return errors.New("no user logged in")
		}

		user, err := s.DataBase.GetUserByName(context.Background(), username)

		if err == sql.ErrNoRows {
			return errors.New("this user does not exist")
		} else if err != nil {
			return err
		}
		return handler(s, c, user)
	}
}
