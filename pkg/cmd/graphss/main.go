package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/email2vimalraj/graphss/pkg/config"
)

var (
	version string
	commit  string
)

func main() {
	os.Exit(Run())
}

var serverFlagSet = flag.NewFlagSet("server", flag.ContinueOnError)

func Run() int {
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
	cfg, err := config.NewCfg(configFilePath)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", cfg)
	return nil
}
