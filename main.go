package main

import (
	"fmt"
	"log"

	"github.com/ajk024/blog_rss/internal/config"
)

const username = "Dustin"

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config 1: %v", err)
	}

	err = cfg.SetUser(username)
	if err != nil {
		log.Fatalf("error setting username: %v", err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("error reading config 2: %v", err)
	}

	fmt.Println(cfg)
}
