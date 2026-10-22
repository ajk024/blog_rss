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

// retrieve content from url, aggregator
func handlerAggregator(s *state, cmd command) error {
	const url = "https://www.wagslane.dev/index.xml"
	feed, err := fetchFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error in fetchFeed: %w", err)
	}

	fmt.Println(feed)
	return nil
}

// add feed to feeds table
func handlerAddFeed(s *state, cmd command) error {
	//takes 2 arguments (in cmd), name and url
	if len(cmd.argument) != 2 {
		return fmt.Errorf("error adding feed to feeds table... usage: addfeed <feed_name feed_url>")
	}

	//look up the current user to attach the feed to
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error getting current user: %w", err)
	}

	feed, err := s.db.AddFeed(context.Background(), database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.argument[0], //feed name
		Url:       cmd.argument[1], //feed url
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error adding feed: %w", err)
	}

	fmt.Println(feed)

	return nil
}

// retrieve all feeds in feeds table
func handlerFeeds(s *state, cmd command) error {
	if len(cmd.argument) > 0 {
		return fmt.Errorf("incorrect input. usage: feeds")
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feeds: %w", err)
	}
	//feeds returns  url name, url, user_id (uuid)

	//retrieve username from uuid
	for i := range len(feeds) {
		username, err := s.db.GetName(context.Background(), feeds[i].UserID)
		if err != nil {
			return fmt.Errorf("error getting username from uuid: %w", err)
		}
		//print url name, url, and username from uuid
		fmt.Printf("%s\t%s\t%s\n", feeds[i].Name, feeds[i].Url, username)
	}

	return nil
}
