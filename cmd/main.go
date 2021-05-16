package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nairobi-gophers/fupisha/api"
	"github.com/nairobi-gophers/fupisha/config"
)

func main() {
	flag.Usage = help
	flag.Parse()

	api, err := api.NewServer()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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
	  fupisha key		- generate a random 32-byte hex-encoded key         
	 `)
}
