package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/nairobi-gophers/fupisha/api"
	"github.com/nairobi-gophers/fupisha/internal/config"
)

var (
	logger       = log.New(os.Stdout, "", log.LstdFlags|log.LUTC)
	useEnvConfig = flag.Bool("e", false, "use environment variables as config")
)

func main() {
	flag.Usage = help
	flag.Parse()

	// router := chi.NewRouter()
	api, err := api.NewServer()
	if err != nil {
		log.Fatal(err)
	}

	cmds := map[string]func(){
		"start": api.Start,
		"key":   config.GenKey,
		"help":  help,
	}

	if cmdFunc, ok := cmds[flag.Arg(0)]; ok {
		cmdFunc()
	} else {
		help()
		os.Exit(1)
	}
}

func help() {
	fmt.Fprintln(os.Stderr, `
	Usage: 
	  fupisha start		- start the server
	  fupisha key			- generate a random 32-byte hex-encoded key         
	 `)
}
