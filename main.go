package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/yashranjan1/gator/internal/command"
	"github.com/yashranjan1/gator/internal/commands"
	"github.com/yashranjan1/gator/internal/config"
	"github.com/yashranjan1/gator/internal/database"
	"github.com/yashranjan1/gator/internal/state"

	_ "github.com/lib/pq"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	db, err := sql.Open("postgres", conf.DBUrl)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	dbQueries := database.New(db)

	s := state.State{
		Config:   &conf,
		DataBase: dbQueries,
	}

	cmds := commands.Commands{
		Callback: make(map[string]func(*state.State, command.Command) error),
	}

	cmds.Register("login", handlerLogin)
	cmds.Register("register", handleRegister)
	cmds.Register("reset", handleReset)
	cmds.Register("users", handleList)
	cmds.Register("agg", handleAggregate)
	cmds.Register("addfeed", middlewareLoggedIn(handleAddFeed))
	cmds.Register("feeds", handleListFeeds)
	cmds.Register("follow", middlewareLoggedIn(handleFollow))
	cmds.Register("following", middlewareLoggedIn(handleFollowing))
	cmds.Register("unfollow", middlewareLoggedIn(handleUnfollow))

	args := os.Args

	if len(args) < 2 {
		fmt.Print("Error: Not enough arguments\n\n")
		fmt.Print("Usage:\n\n")
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
