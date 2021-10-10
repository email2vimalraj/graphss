package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/email2vimalraj/graphss/pkg/domains"
	"github.com/email2vimalraj/graphss/pkg/http"
)

var (
	version string
	commit  string
)

func main() {
	domains.Version = strings.TrimPrefix(version, "")
	domains.Commit = commit

	// // Setup signal handlers.
	// ctx, cancel := context.WithCancel(context.Background())
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt)
	// go func() { <-c; cancel() }()

	os.Exit(run())
}

var serverFlagSet = flag.NewFlagSet("server", flag.ContinueOnError)

func run() int {
	var (
		configFilePath = serverFlagSet.String("config", "", "path to the configuration file")
		// pidFilePath    = serverFlagSet.String("pidfile", "", "path to the pid file")
		v = serverFlagSet.Bool("v", false, "print the version and exit")
	)

	if err := serverFlagSet.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	if *v {
		fmt.Printf("Version: %s, Commit: %s\n", version, commit)
		return 0
	}

	if err := runServer(*configFilePath); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	return 0
}

func runServer(configFilePath string) error {
	// Setup signal handlers.
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() { <-c; cancel() }()

	// Initialize the server.
	server, err := http.InitializeServer(configFilePath)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", server.Cfg)

	if err := server.Open(); err != nil {
		return err
	}

	// Wait for CTRL-C.
	<-ctx.Done()

	return nil
}
