package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/ajk024/blog_rss/internal/config"
	"github.com/ajk024/blog_rss/internal/database"
	_ "github.com/lib/pq"
)

//const username = "Dustin"

func main() {
	//read config file
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config 1: %v", err)
	}

	//open database
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}
	dbQueries := database.New(db)

	//build state
	s := state{cfg: &cfg, db: dbQueries}

	//make command map
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
