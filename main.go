package main

import (
	"fmt"
	"log"

	"github.com/ajk024/blog_rss/internal/config"
)

//const username = "Dustin"

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config 1: %v", err)
	}

	s := state{cfg: &cfg}
	c := commands{registeredCommands: make(map[string]func(*state, command) error)}

	c.register("login", handlerLogin)

	/*
		err = cfg.SetUser(username)
		if err != nil {
			log.Fatalf("error setting username: %v", err)
		}*/

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("error reading config 2: %v", err)
	}

	fmt.Println(cfg)
}

type state struct {
	cfg *config.Config
}

type command struct {
	name     string
	argument []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.argument) != 1 {
		return fmt.Errorf("username not entered\n")
	}

	if err := s.cfg.SetUser(cmd.argument[0]); err != nil {
		return fmt.Errorf("error setting username: %v\n", err)
	}

	fmt.Printf("user has been set: %s\n", cmd.argument[0])

	return nil
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.registeredCommands[cmd.name]
	if ok == false {
		return fmt.Errorf("handler not registered")
	}
	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}
