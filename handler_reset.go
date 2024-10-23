package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error deleting all users: %v", err)
	}

	fmt.Println("Reseting database")

	return nil
}
