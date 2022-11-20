package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	shortcuts "github.com/snkzt/shortcut-peeper"
)

const (
	exitFail = 1
)

// TODO: add language flag diversion, unit test

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(exitFail)
	}
}

func run() error {

	const usage = `Usage of speep: speep <command> <flag1> ... <flag2> ...
	
	Command Options:
		get [--all] | [--name <name>]
		add --name <name> --key <key>
		delete  [--all] | [--name <name>]

	Flag Options:
	  -a, --all retrieve all shortcuts
	  -n, --name name of the registered shortcut key: Use \"\" for more than one word e.g. -name \"to the back of the line\".
	  -k, --key registered shortcut key`

	const longFlagName = "Use \"for more than one word e.g. --name \"to the back of the line\""
	const shortFlagName = "Use \"for more than one word e.g. -n \"to the back of the line\""

	// peeper subcommand "get"
	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	// Flags for "get" subcommand
	// ("flag name", default, "Flag explanation")
	getAll := getCmd.Bool("all", false, "Get full shortcut list")
	getAllShort := getCmd.Bool("a", false, "Get full shortcut list")
	getByName := getCmd.String("name", "", "Find a shortcut with a name e.g. \"speep get --name Copy\"."+longFlagName)
	getByNameShort := getCmd.String("n", "", "Find a shortcut with a name e.g. \"speep get -n Copy\"."+shortFlagName)

	// peeper subcommand "add"
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	// Flags(inputs) for "add" subcommand
	addName := addCmd.String("name", "", "Name of the shortcut: e.g. \"speep add --name copy --key Ctrl+C\"."+longFlagName)
	addNameShort := addCmd.String("n", "", "Name of the shortcut: e.g. \"speep add -n copy -k Ctrl+C\"."+shortFlagName)
	addKey := addCmd.String("key", "", "Key of the shortcut key\n e.g. \"speep add --name copy --key Ctrl+C\"")
	addKeyShort := addCmd.String("k", "", "Key of the shortcut key\n e.g. \"speep add -n copy -k Ctrl+C\"")

	// peeper subcommand "delete"
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	// Flags for "delete" subcommand
	deleteAll := deleteCmd.Bool("all", false, "Delete whole shortcut list with \"speep delete --all\"")
	deleteAllShort := deleteCmd.Bool("a", false, "Delete whole shortcut list with \"speep delete -a\"")
	deleteByName := deleteCmd.String("name", "", "Delete a shortcut by name e.g. \"speep delete --name Copy\"."+longFlagName)
	deleteByNameShort := deleteCmd.String("n", "", "Delete a shortcut by name e.g. \"speep delete -n Copy\"."+shortFlagName)

	// Return usage if no command specified
	if len(os.Args) < 2 {
		return fmt.Errorf("No command and flag specified. \n %s", usage)
	}

	// Check user input command and proceed to each process
	command := os.Args[1]
	switch {
	case command == "get":
		return HandleGet(getCmd, getAll, getAllShort, getByName, getByNameShort)
	case command == "add":
		return HandleAdd(addCmd, addName, addNameShort, addKey, addKeyShort)
	case command == "delete":
		return HandleDelete(deleteCmd, deleteAll, deleteAllShort, deleteByName, deleteByNameShort)
	case command == "help" || command == "--help" || command == "-h":
		return Handlehelp(usage)
	default: // Return usage when user input command doesn't exist
		return fmt.Errorf("speep:'%s' is not a speep command. \n %s", os.Args[1], usage)
	}
}

func HandleGet(getCmd *flag.FlagSet, all *bool, allShort *bool, name *string, nameShort *string) error {
	if *name == "" {
		*name = *nameShort
	}

	getCmd.Parse(os.Args[2:])

	if !*all && !*allShort && *name == "" {
		getCmd.PrintDefaults()
		return nil
	}

	if *all || *allShort {
		// Return full shortcut list
		shortcuts, err := shortcuts.GetShortcuts()
		if err != nil {
			return fmt.Errorf("please add new shortcuts, no shortcut key registered: %w", err)
		}
		fmt.Println("Name \t Shortcut key \n")
		for _, shortcut := range shortcuts {
			fmt.Printf("%v \t %v \n", shortcut.Name, shortcut.ShortcutKey)
		}
	}

	if *name != "" {
		shortcuts, err := shortcuts.GetShortcuts()
		if err != nil {
			return fmt.Errorf("failed to acquire the existing list: %w", err)
		}
		name := *name
		for _, shortcut := range shortcuts {
			if strings.Contains(shortcut.ShortcutKey, name) {
				fmt.Println("Name \t Shortcut key")
				fmt.Printf("%v \t %v \n", shortcut.Name, shortcut.ShortcutKey)
			}
		}
	}
	return nil
}

func ValidateNewShortcutKey(addCmd *flag.FlagSet, name *string, key *string) error {
	if *name == "" || *key == "" {
		addCmd.PrintDefaults()
		return errors.New("name and the shortcut key are required to add a shortcut key")
	}
	return nil
}

func HandleAdd(addCmd *flag.FlagSet, name *string, nameShort *string, newShortcut *string, newShortcutShort *string) error {
	addCmd.Parse(os.Args[2:])

	if *name == "" {
		*name = *nameShort
	}
	if *newShortcut == "" {
		*newShortcut = *newShortcutShort
	}

	err := ValidateNewShortcutKey(addCmd, name, newShortcut)
	if err != nil {
		return nil
	}

	var allShortcuts []shortcuts.Shortcut
	shortcut := shortcuts.Shortcut{
		Name:        *name,
		ShortcutKey: *newShortcut,
	}

	err = shortcuts.CheckNameDuplication(name)
	if err != nil {
		return fmt.Errorf("the name \"%s\" %w", *name, err)
	}

	allShortcuts, _ = shortcuts.GetShortcuts()
	allShortcuts = append(allShortcuts, shortcut)
	err = shortcuts.SaveShortcuts(allShortcuts)
	if err != nil {
		return fmt.Errorf("failed to save the updated list: %w", err)
	}

	fmt.Printf("New shortcut \"%v\" successfully registered", *name)
	return nil
}

func HandleDelete(deleteCmd *flag.FlagSet, all *bool, allShort *bool, name *string, nameShort *string) error {
	deleteCmd.Parse(os.Args[2:])

	if *name == "" {
		*name = *nameShort
	}
	if !*all && !*allShort && *name == "" {
		deleteCmd.PrintDefaults()
		return nil
	}

	if *all || *allShort {
		// Delete full shortcut list
		err := shortcuts.DeleteShortcuts()
		if err != nil {
			return fmt.Errorf("failed to delete the shortcut list: %w", err)
		}
		fmt.Println("Shortcut list deleted")
	}

	if *name != "" {
		err := shortcuts.DeleteShortcut(*name)
		if err != nil {
			return fmt.Errorf("failed to remove an item from the list: %w", err)
		}
		fmt.Printf("Shortcut %v successfully removed from the Shortcut key list", *name)
		// Add return error of the target not exist if the name doesn't exist in the list
	}
	return nil
}

func Handlehelp(usage string) error {
	return errors.New(usage)
}
