package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yashranjan1/gator/internal/command"
	"github.com/yashranjan1/gator/internal/database"
	"github.com/yashranjan1/gator/internal/state"
)

func handlerLogin(s *state.State, cmd command.Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("the login handler expects a single argument, the username")
	}

	_, err := s.DataBase.GetUserByName(context.Background(), cmd.Args[0])

	if err == sql.ErrNoRows {
		return errors.New("user does not exist")
	} else if err != nil {
		return err
	}

	err = s.Config.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Println("User has been set")
	return nil
}

func handleRegister(s *state.State, cmd command.Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("the register handler expects a single argument, the username")
	}

	_, err := s.DataBase.GetUserByName(context.Background(), cmd.Args[0])

	if err == nil {
		return errors.New("a user with this name already exists")
	}

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
	}

	userCreated, err := s.DataBase.CreateUser(context.Background(), params)

	err = s.Config.SetUser(userCreated.Name)

	if err != nil {
		return err
	}

	fmt.Printf("SUCCESS: user \"%s\" created, exiting....", userCreated.Name)

	return nil

}
