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
	getByKeyword := getCmd.String("keyword", "", "Find a shortcut with keyword ")

	// peeper subcommand "add"
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	// Flags(inputs) for "add" subcommand
	addName := addCmd.String("name", "", "Name of the shortcut")
	addShortcut := addCmd.String("shortcut", "", "Shortcut")

	// peeper subcommand "delete"
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	// Flags for "delete" subcommand
	deleteAll := deleteCmd.Bool("all", false, "Delete whole shortcut list")
	deleteByName := deleteCmd.String("name", "", "Delete a shortcut by name")

	if len(os.Args) < 2 {
		return errors.New("exepected either 'get', 'add', 'delete' subcommands")
	}

	switch os.Args[1] {
	case "get":
		return HandleGet(getCmd, getAll, getByKeyword)
	case "add":
		return HandleAdd(addCmd, addName, addShortcut)
	case "delete":
		return HandleDelete(deleteCmd, deleteAll, deleteByName)
		// TODO: Do we need to add default and return error for the case of the command doesn't exist?
	}
	return nil
}

func HandleGet(getCmd *flag.FlagSet, all *bool, keyword *string) error {
	getCmd.Parse(os.Args[2:])

	if !*all && *keyword == "" {
		getCmd.PrintDefaults()
		return errors.New("specify the target shortcut with all or keyword flag")
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

	if *keyword != "" {
		shortcuts, err := shortcuts.GetShortcuts()
		if err != nil {
			return fmt.Errorf("failed to acquire the existing list: %w", err)
		}
		keyword := *keyword
		for _, shortcut := range shortcuts {
			if strings.Contains(shortcut.ShortcutKey, keyword) {
				fmt.Println("Name \t Shortcut key")
				fmt.Printf("%v \t %v \n", shortcut.Name, shortcut.ShortcutKey)
			}
		}
	}
	return nil
}

func ValidateNewShortcutKey(addCmd *flag.FlagSet, name *string, shortcut *string) error {
	addCmd.Parse(os.Args[2:])
	if *name == "" || *shortcut == "" {
		addCmd.PrintDefaults()
		return errors.New("name and the shortcut key are required to add a shortcut key")
	}
	return nil
}

func HandleAdd(addCmd *flag.FlagSet, name *string, newShortcut *string) error {
	ValidateNewShortcutKey(addCmd, name, newShortcut)

	var allShortcuts []shortcuts.Shortcut
	shortcut := shortcuts.Shortcut{
		Name:        *name,
		ShortcutKey: *newShortcut,
	}

	allShortcuts, _ = shortcuts.GetShortcuts()
	allShortcuts = append(allShortcuts, shortcut)
	err := shortcuts.SaveShortcuts(allShortcuts)
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
		return errors.New("specify the target shortcut with all or name flag")
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
	}
	return nil
}
