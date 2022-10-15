package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// TODO: add language flag diversion, unit unit test

func main() {
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
		fmt.Println("Exepected either 'get', 'add', 'delete' subcommands")
	}

	switch os.Args[1] {
	case "get":
		HandleGet(getCmd, getAll, getByKeyword)
	case "add":
		HandleAdd(addCmd, addName, addShortcut)
	case "delete":
		HandleDelete(deleteCmd, deleteAll, deleteByName)
		// TODO: Do we need to add default and return error for the case of the command doesn't exist?
	}
}

func HandleGet(getCmd *flag.FlagSet, all *bool, keyword *string) {
	getCmd.Parse(os.Args[2:])

	if *all == false && *keyword == "" {
		fmt.Println("Specify the target shortcut with all or keyword flag")
		getCmd.PrintDefaults()
		os.Exit(1)
	}

	if *all {
		// Return full shortcut list
		shortcuts := getShortcuts()
		fmt.Println("Name \t Shortcut key \n")
		for _, shortcut := range shortcuts {
			fmt.Printf("%v \t %v \n", shortcut.Name, shortcut.ShortcutKey)
		}
		return
	}

	if *keyword != "" {
		shortcuts := getShortcuts()
		keyword := *keyword
		for _, shortcut := range shortcuts {
			if strings.Contains(shortcut.ShortcutKey, keyword) {
				fmt.Println("Name \t Shortcut key \n")
				fmt.Printf("%v \t %v \n", shortcut.Name, shortcut.ShortcutKey)
			}
		}
	}
}

func ValidateNewShortcutKey(addCmd *flag.FlagSet, name *string, shortcut *string) {
	addCmd.Parse(os.Args[2:])
	if *name == "" || *shortcut == "" {
		fmt.Print("Name and the shortcut key are required to add a shortcut key")
		addCmd.PrintDefaults()
		os.Exit(1)
	}
}

func HandleAdd(addCmd *flag.FlagSet, name *string, newShortcut *string) {
	ValidateNewShortcutKey(addCmd, name, newShortcut)

	shortcut := Shortcut{
		Name:        *name,
		ShortcutKey: *newShortcut,
	}

	shortcuts := getShortcuts()
	shortcuts = append(shortcuts, shortcut)
	saveShortcuts(shortcuts)

	fmt.Printf("New shortcut %v successfully added to the Shortcut key list", *name)
}

func HandleDelete(deleteCmd *flag.FlagSet, all *bool, name *string) {
	deleteCmd.Parse(os.Args[2:])
	if !*all && *name == "" {
		fmt.Println("Specify the target shortcut with all or name flag")
		deleteCmd.PrintDefaults()
		os.Exit(1)
	}

	if *all {
		// Delete full shortcut list
		deleteShortcuts()
		fmt.Println("Shortcut list deleted")
	}

	if *name != "" {
		deleteShortcut(*name)
		fmt.Printf("Shortcut %v successfully removed from the Shortcut key list", *name)
	}
}
