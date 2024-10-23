package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/fummbly/gatorcli/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	name := cmd.Args[0]
	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("User does not exist: %w", err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("Error setting user: %v", err)
	}

	fmt.Printf("User switched to %s\n", name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	name := cmd.Args[0]
	user, _ := s.db.GetUser(context.Background(), name)
	if user.Name != "" {
		fmt.Println("User already exists")
		os.Exit(1)
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        int32(uuid.New().ID()),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
	})
	if err != nil {
		return fmt.Errorf("Error creating user: %v", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("Error setting user: %v", err)
	}

	user, err = s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("Error getting user info: %v", err)
	}

	printUser(user)

	return nil
}

func printUser(user database.User) {
	fmt.Printf(" * ID:       %v\n", user.ID)
	fmt.Printf(" * Name:     %v\n", user.Name)
}
