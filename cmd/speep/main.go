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
	// peeper subcommand "get"
	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	// Flags for "get" subcommand
	// ("flag name", default, "Flag explanation")
	getAll := getCmd.Bool("all", false, "Get full shortcut list")
	getByName := getCmd.String("name", "", "Find a shortcut with a name e.g. \"speep get -name Copy\" ")

	// peeper subcommand "add"
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	// Flags(inputs) for "add" subcommand
	addName := addCmd.String("name", "", "Name of the shortcut: Use \"\" for more than one word e.g. -name \"to the back of the line\"")
	addKey := addCmd.String("key", "", "Key of the shortcut key\n e.g. \"speep add -name copy -shortcut Ctrl+C\"")

	// peeper subcommand "delete"
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	// Flags for "delete" subcommand
	deleteAll := deleteCmd.Bool("all", false, "Delete whole shortcut list with \"speep delete -all\"")
	deleteByName := deleteCmd.String("name", "", "Delete a shortcut by name e.g. \"speep delete -name Copy\"")

	if len(os.Args) < 2 {
		return errors.New("Please check how to use the app by typing either \"speep get\", \"speep add\" or \"speep delete\" for more details.")
	}

	switch os.Args[1] {
	case "get":
		return HandleGet(getCmd, getAll, getByName)
	case "add":
		return HandleAdd(addCmd, addName, addKey)
	case "delete":
		return HandleDelete(deleteCmd, deleteAll, deleteByName)
	default: // When user chose command which doesn't exist
		fmt.Println("Please check how to use the app by typing either \"speep get\", \"speep add\" or \"speep delete\" for more details.")

	}
	return nil
}

func HandleGet(getCmd *flag.FlagSet, all *bool, name *string) error {
	getCmd.Parse(os.Args[2:])

	if !*all && *name == "" {
		getCmd.PrintDefaults()
		return nil
	}

	if *all {
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
	addCmd.Parse(os.Args[2:])
	if *name == "" || *key == "" {
		addCmd.PrintDefaults()
		return errors.New("name and the shortcut key are required to add a shortcut key")
	}
	return nil
}

func HandleAdd(addCmd *flag.FlagSet, name *string, newShortcut *string) error {
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

func HandleDelete(deleteCmd *flag.FlagSet, all *bool, name *string) error {
	deleteCmd.Parse(os.Args[2:])
	if !*all && *name == "" {
		deleteCmd.PrintDefaults()
		return nil
	}

	if *all {
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
