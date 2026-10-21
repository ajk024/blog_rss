package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ajk024/blog_rss/internal/database"
	"github.com/google/uuid"
)

// login to database with username
func handlerLogin(s *state, cmd command) error {
	if len(cmd.argument) != 1 {
		return fmt.Errorf("error... username usage: login <username>")
	}

	//check if username exists in database
	_, err := s.db.GetUser(context.Background(), cmd.argument[0])
	if err != nil {
		return fmt.Errorf("error... username not in database: %w", err)
	}

	//set user in config
	if err := s.cfg.SetUser(cmd.argument[0]); err != nil {
		return fmt.Errorf("error setting username 1: %w", err)
	}
	fmt.Printf("user has been set: %s\n", cmd.argument[0])
	return nil
}

// register a new username in database
func handlerRegister(s *state, cmd command) error {
	if len(cmd.argument) != 1 {
		return fmt.Errorf("error... username usage: register <username>")
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.argument[0],
	})
	if err != nil {
		return fmt.Errorf("could not create user: %w", err) //Error code 1 if duplicate user is entered
	}

	//set user in config
	if err := s.cfg.SetUser(cmd.argument[0]); err != nil {
		return fmt.Errorf("error setting username 2: %w", err)
	}

	fmt.Println("user was created!")
	fmt.Printf("%v\n", user)

	return nil
}

// reset users database by deleting all users
func handlerReset(s *state, cmd command) error {
	if err := s.db.DeleteUsers(context.Background()); err != nil {
		return fmt.Errorf("error deleting all users: %w", err)
	}

	fmt.Println("all users deleted from database")
	return nil
}

// get list of usernames in database
func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error getting users in database: %w", err)
	}

	//fmt.Println(users)
	for i := 0; i < len(users); i++ {
		fmt.Printf("* %s", users[i])

		if users[i] == s.cfg.CurrentUserName {
			fmt.Print(" (current)")
		}
		fmt.Println()

	}

	return nil
}
