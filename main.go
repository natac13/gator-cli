package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/natac13/gator-cli/internal/database"
)

func main() {
	state, err := newState()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", state.config.DbUrl)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	state.db = dbQueries

	cmds := newCommands()
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerFeedCreate))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFeedFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	input := os.Args

	if len(input) < 2 {
		log.Fatal("missing command")
	}

	cmd := command{
		name: input[1],
		args: input[2:],
	}

	err = cmds.run(state, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
