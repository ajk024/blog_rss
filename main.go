package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/ajk024/blog_rss/internal/config"
	"github.com/ajk024/blog_rss/internal/database"
	_ "github.com/lib/pq"
)

//const username = "Dustin"

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config 1: %v", err)
	}

	s := state{cfg: &cfg}
	cmds := commands{registeredCommands: make(map[string]func(*state, command) error)}

	cmds.register("login", handlerLogin)

	args := os.Args
	if len(args) < 2 {
		log.Fatal("not enough arguments")
	}

	cmd := command{name: args[1], argument: args[2:]}

	//run command
	if err := cmds.run(&s, cmd); err != nil {
		log.Fatalf("error running command: %v", err)
	}

	db, err := sql.Open("postgres", s.cfg.DBURL)
	dbQueries := database.New(db)
	s.db = dbQueries
}

type state struct {
	cfg *config.Config
	db  *database.Queries
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
		return fmt.Errorf("username usage: login <username>")
	}

	if err := s.cfg.SetUser(cmd.argument[0]); err != nil {
		return fmt.Errorf("error setting username: %w", err)
	}

	fmt.Printf("user has been set: %s\n", cmd.argument[0])

	return nil
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.registeredCommands[cmd.name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.name)
	}
	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}
