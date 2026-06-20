package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/nikdotcodes/gator/internal/config"
	"github.com/nikdotcodes/gator/internal/database"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	c, _ := config.Read()
	db, err := sql.Open("postgres", c.DBUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	s := &state{cfg: c, db: dbQueries}

	cmnds := commands{commandNames: map[string]func(*state, command) error{}}
	cmnds.register("login", handlerLogin)
	cmnds.register("register", handlerRegister)
	cmnds.register("reset", handlerReset)
	cmnds.register("users", handlerUsers)
	cmnds.register("agg", handlerAgg)

	if len(os.Args) < 2 {
		fmt.Println("No command specified")
		os.Exit(1)
	}

	ran := command{name: os.Args[1], arguments: os.Args[2:]}
	if err := cmnds.run(s, ran); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
