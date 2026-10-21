package main

import "fmt"

func handlerLogin(s *state, cmd command) error {
	if len(cmd.argument) != 1 {
		return fmt.Errorf("username usage: login <username>")
	}

	if err := s.cfg.SetUser(cmd.argument[0]); err != nil {
		return fmt.Errorf("error setting username: %w", err)
	}

	fmt.Printf("user has been set: %s\n", cmd.argument[0])

	return nil
}
