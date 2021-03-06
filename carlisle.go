package main

import (
	"fmt"
	"log"
	"os"
	"sort"
)

type Command interface {
	Execute(args []string) error
	Help() string
}

var (
	commands     = map[string]Command{} // Filled in in init()s
	commandNames []string
)

func usage(status int) {
	fmt.Printf(`Usage:

    %s COMMAND [arg1] [arg2] ...

where COMMAND is one of %v
(Type '%[1]s COMMAND help' to see information about a specific command.)
`, os.Args[0], commandNames)
	os.Exit(status)
}

func main() {
	for name := range commands {
		commandNames = append(commandNames, name)
	}
	sort.Strings(commandNames)

	if len(os.Args) < 2 {
		usage(-1)
	}
	switch os.Args[1] {
	case "-h", "-help", "--help", "help":
		usage(0)
	}
	command, ok := commands[os.Args[1]]
	if !ok {
		usage(-1)
	}
	if len(os.Args) >= 3 {
		switch os.Args[2] {
		case "-h", "-help", "--help", "help":
			fmt.Println(command.Help())
			os.Exit(0)
		}
	}
	if err := command.Execute(os.Args[2:]); err != nil {
		log.Fatal(err)
	}
}
