package main

import "fmt"

import "os/user"

import "github.com/yashranjan1/gator/internal/config"

func main() {
	conf, err := config.Read()
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	currentUser, err := user.Current()
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	conf.SetUser(currentUser.Username)

	conf, err = config.Read()
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	fmt.Printf("%s: %s\n", "Url", conf.DBUrl)
	fmt.Printf("%s: %s\n", "Username", conf.CurrentUser)
}
