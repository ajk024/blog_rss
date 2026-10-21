package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ajk024/blog_rss/internal/database"
	"github.com/google/uuid"
)

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
