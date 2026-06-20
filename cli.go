package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nikdotcodes/gator/internal/database"
	"github.com/nikdotcodes/gator/internal/rss"
)

type commands struct {
	commandNames map[string]func(*state, command) error
}

type command struct {
	name      string
	arguments []string
}

func (c *commands) run(s *state, cmd command) error {
	cmdFunc, ok := c.commandNames[cmd.name]
	if !ok {
		return errors.New("unknown command: " + cmd.name)
	}

	if err := cmdFunc(s, cmd); err != nil {
		return err
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	_, ok := c.commandNames[name]
	if ok {
		fmt.Println("command " + name + " already registered")
		return
	}

	c.commandNames[name] = f
	return
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.arguments) > 0 {
		return errors.New("arguments are not supported")
	}

	ctx := context.Background()
	if err := s.db.DeleteAllUsers(ctx); err != nil {
		return err
	}

	fmt.Println("users table reset")
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return errors.New("login expects username argument.")
	}

	if err := s.cfg.SetUser(cmd.arguments[0]); err != nil {
		return err
	}

	ctx := context.Background()
	queryUsr := cmd.arguments[0]

	_, err := s.db.GetUserByName(ctx, queryUsr)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("user " + queryUsr + " does not exist - cannot login")
		return err
	}

	fmt.Printf("login successful - User %s has been logged in.\n", cmd.arguments[0])

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return errors.New("register expects username argument.")
	}

	ctx := context.Background()

	queryUsr := cmd.arguments[0]

	_, err := s.db.GetUserByName(ctx, queryUsr)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Printf("registering %s\n", queryUsr)
	} else if err != nil {
		return err
	} else {
		return fmt.Errorf("user %s already registered", queryUsr)
	}

	createUsr := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
	}

	u, err := s.db.CreateUser(ctx, createUsr)
	if err != nil {
		return err
	}
	if err := s.cfg.SetUser(u.Name); err != nil {
		return err
	}
	fmt.Printf("User %s has been registered.\n", u.Name)
	return nil
}

func handlerUsers(s *state, cmd command) error {
	if len(cmd.arguments) > 0 {
		return errors.New("arguments are not supported")
	}

	ctx := context.Background()
	usrs, err := s.db.GetUsers(ctx)
	if err != nil {
		return err
	}
	for _, usr := range usrs {
		if usr == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", usr)
		} else {
			fmt.Printf("* %s\n", usr)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	//if len(cmd.arguments) != 1 {
	//	return errors.New("agg expects RSS feed argument.")
	//}
	//
	//ctx := context.Background()
	//feed, err := rss.FetchFeed(ctx, cmd.arguments[0])
	//if err != nil {
	//	return err
	//}
	//
	////fmt.Println(feed)
	//fmt.Println(feed)
	//return nil

	feedUrl := "https://www.wagslane.dev/index.xml"
	ctx := context.Background()
	feed, err := rss.FetchFeed(ctx, feedUrl)
	if err != nil {
		return err
	}
	fmt.Println(feed)
	return nil
}
