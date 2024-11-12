package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/fummbly/gatorcli/internal/config"
	"github.com/fummbly/gatorcli/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	config, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", config.DBUrl)
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}

	dbQueries := database.New(db)

	programState := &state{
		cfg: &config,
		db:  dbQueries,
	}

	cmds := commands{
		registeredCommands: make(map[string]command),
	}

	cmds.register("login", command{handler: handlerLogin, descr: "login in a user. Usage: login <username>"})
	cmds.register("register", command{handler: handlerRegister, descr: "register a user. Usage: register <username>"})
	cmds.register("reset", command{handler: handlerReset, descr: "reset the database"})
	cmds.register("users", command{handler: handlerUsers, descr: "get a list of all users"})
	cmds.register("agg", command{handler: handlerAgg, descr: "aggregate all the feeds"})
	cmds.register("addfeed", command{handler: middlewareLoggedIn(handlerAddFeed), descr: "add a feed Usage: addfeed <feed-name> <feed-url>"})
	cmds.register("feeds", command{handler: handlerGetFeeds, descr: "get all feeds registered"})
	cmds.register("follow", command{handler: middlewareLoggedIn(handlerFollow), descr: "follow a feed for a user Usage: follow <feed-url>"})
	cmds.register("following", command{handler: middlewareLoggedIn(handlerFollowing), descr: "get a list of all feeds the current user is following"})
	cmds.register("unfollow", command{handler: middlewareLoggedIn(handlerUnfollow), descr: "unfollow a feed for current user Usage: unfollow <feed-url>"})
	cmds.register("browse", command{handler: middlewareLoggedIn(handlerBrowse), descr: "browse posts for the current user"})
	cmds.register("removefeed", command{handler: handlerRemoveFeed, descr: "remove a feed from the list Usage: removefeed <feed-url>"})

	if len(os.Args) < 2 {
		fmt.Println("Usage cli <command> [args...]")
		return
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

  if cmdName == "help" {
    cmds.help()
    return
  }

	err = cmds.run(programState, input{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}

}

func middlewareLoggedIn(handler func(s *state, cmd input, user database.User) error) func(*state, input) error {
	return func(s *state, cmd input) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUsername)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}

}

