package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/cristianoliveira/ergo/commands"
	"github.com/cristianoliveira/ergo/proxy"
)

// VERSION of ergo
// When ergo is build without a proper tag/release it is named as `unofficial version`.
// For instance, installing through `go get github.com/cristianoliveira/ergo` or `go build`.
var VERSION = "unofficial version"

//USAGE details the usage for ergo proxy.
const USAGE = `
Ergo proxy.
The management apps running on different ports made easy through custom local domains.

Usage:
  ergo run [options]
  ergo list [options]
  ergo list-names [options]
  ergo url <name> [options]
  ergo setup [linux-gnome|osx|windows] [-remove] [options]
  ergo add <service-name> <host:port> [options]
  ergo remove <service-name|host:port> [options]

Options:
  -h      Shows this message.
  -v      Shows ergo's version.
  -config     Set the config file to the proxy.
  -domain     Set a custom domain for services.
  -p          Set ports to proxy.
  -V          Set verbosity on output.

setup:
  -remove     Set remove proxy configurations.
`

func prepareSubCommand(args []string) (commands.Command, *proxy.Config) {
	// Fail fast if we didn't receive a command argument
	if len(args) == 1 {
		return nil, nil
	}

	config := proxy.NewConfig()

	command := flag.NewFlagSet(args[1], flag.ExitOnError)
	command.StringVar(&config.ConfigFile, "config", "./.ergo", "Set the services file")
	command.StringVar(&config.Domain, "domain", ".dev", "Set a custom domain for services")
	command.StringVar(&config.Port, "p", "2000", "Set port to the proxy")
	command.BoolVar(&config.Verbose, "V", false, "Set verbosity on proxy output")

	switch args[1] {
	case "list":
		command.Parse(args[2:])
		return commands.ListCommand{}, config

	case "list-names":
		command.Parse(args[2:])
		return commands.ListNameCommand{}, config

	case "setup":
		if len(args) <= 2 {
			return nil, nil
		}

		system := args[2]
		setupCommand := commands.SetupCommand{System: system}

		command.BoolVar(&setupCommand.Remove, "remove", false, "Set remove proxy configurations.")
		command.Parse(args[3:])

		return setupCommand, config

	case "url":
		if len(args) < 3 {
			return nil, nil
		}

		name := args[2]
		command.Parse(args[3:])

		return commands.URLCommand{FilterName: name}, config

	case "run":
		command.Parse(args[2:])

		return commands.RunCommand{}, config

	case "add":
		if len(args) < 4 {
			return nil, nil
		}

		name := args[2]
		url := args[3]
		service := proxy.Service{Name: name, URL: url}

		command.Parse(args[4:])

		return commands.AddServiceCommand{Service: service}, config

	case "remove":
		if len(args) <= 2 {
			return nil, nil
		}

		nameOrURL := args[2]
		service := proxy.Service{Name: nameOrURL, URL: nameOrURL}

		command.Parse(args[3:])

		return commands.RemoveServiceCommand{Service: service}, config
	}

	return nil, nil
}

func execute(command commands.Command, config *proxy.Config) {
	result, err := command.Execute(config)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}
}

var help = flag.Bool("h", false, "Shows ergo's help.")
var version = flag.Bool("v", false, "Shows ergo's version.")

func main() {
	flag.Parse()

	if *version {
		fmt.Println(VERSION)
		return
	}

	if *help {
		fmt.Println(USAGE)
		return
	}

	command, config := prepareSubCommand(os.Args)
	if command == nil {
		fmt.Println(USAGE)

	} else {

		err := config.LoadServices()
		if err != nil {
			log.Fatalf("Could not load services: %v\n", err)
		}

		execute(command, config)
	}
}
