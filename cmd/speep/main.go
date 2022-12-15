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

const usage = `Usage of speep: speep <command> <flag1> ... <flag2> ...
	
	Command Options:
		get [--all] | [--name <name>]
		add --category <category> --name <name> --key <key>
		delete [--all] | [--category <category> --name <name>]

	Flag Options:
		-a, --all retrieve all shortcuts
		-c, --category name of the category of the registered shortcut key: e.g. shell
		-n, --name name of the registered shortcut key: Use "" for more than one word e.g. -name \"to the back of the line\"
		-k, --key registered shortcut key`

// TODO:unit test

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(exitFail)
	}
}

func run() error {

	const longFlagName = "Use \"\" for more than one word e.g. --name \"to the back of the line\""
	const shortFlagName = "Use \"\" for more than one word e.g. -n \"to the back of the line\""

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
	addCategory := addCmd.String("category", "", "Specify category to add a shortcut")
	addCategoryShort := addCmd.String("c", "", "Specify category to add a shortcut")
	addName := addCmd.String("name", "", "Name of the shortcut: e.g. \"speep add --name copy --key Ctrl+C\"."+longFlagName)
	addNameShort := addCmd.String("n", "", "Name of the shortcut: e.g. \"speep add -n copy -k Ctrl+C\"."+shortFlagName)
	addKey := addCmd.String("key", "", "Key of the shortcut key\n e.g. \"speep add --name copy --key Ctrl+C\"")
	addKeyShort := addCmd.String("k", "", "Key of the shortcut key\n e.g. \"speep add -n copy -k Ctrl+C\"")

	// peeper subcommand "delete"
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	// Flags for "delete" subcommand
	deleteAll := deleteCmd.Bool("all", false, "Delete whole shortcut list with \"speep delete --all\"")
	deleteAllShort := deleteCmd.Bool("a", false, "Delete whole shortcut list with \"speep delete -a\"")
	deleteCategory := deleteCmd.String("category", "", "Specify category to delete a shortcut by name")
	deleteCategoryShort := deleteCmd.String("c", "", "Specify category to delete a shortcut by name")
	deleteByName := deleteCmd.String("name", "", "Delete a shortcut by name e.g. \"speep delete --name Copy\". "+longFlagName)
	deleteByNameShort := deleteCmd.String("n", "", "Delete a shortcut by name e.g. \"speep delete -n Copy\". "+shortFlagName)

	// Return usage if no command specified
	if len(os.Args) < 2 {
		return fmt.Errorf("no command and flag specified. \n %s", usage)
	}

	// Check user input command and proceed to each process
	command := os.Args[1]
	switch {
	case command == "get":
		getCmd.Parse(os.Args[2:])

		if *getByName == "" && *getByNameShort != "" {
			*getByName = *getByNameShort
		}

		return handleGet(getCmd, getAll, getAllShort, getByName)
	case command == "add":
		addCmd.Parse(os.Args[2:])

		if *addCategory == "" {
			*addCategory = *addCategoryShort
		}
		if *addName == "" {
			*addName = *addNameShort
		}
		if *addKey == "" {
			*addKey = *addKeyShort
		}
		return handleAdd(addCmd, addCategory, addName, addKey)
	case command == "delete":
		deleteCmd.Parse(os.Args[2:])
		if *deleteCategory == "" {
			*deleteCategory = *deleteCategoryShort
		}
		if *deleteByName == "" {
			*deleteByName = *deleteByNameShort
		}
		return handleDelete(deleteCmd, deleteAll, deleteAllShort, deleteCategory, deleteByName)
	case command == "help" || command == "--help" || command == "-h":
		return handlehelp(usage)
	default: // Return usage when user input command doesn't exist
		return fmt.Errorf("speep:'%s' is not a speep command. \n %s", os.Args[1], usage)
	}
}

func handleGet(getCmd *flag.FlagSet, all *bool, allShort *bool, name *string) error {

	if !*all && !*allShort && *name == "" {
		fmt.Println("Flag and/or argument for the command get is missing.")
		fmt.Print(handlehelp(usage))
		return nil
	}

	if *all || *allShort {
		// Return full shortcut list
		shortcuts, err := shortcuts.GetShortcuts()
		if err != nil {
			return fmt.Errorf("please add new shortcuts, no shortcut key registered: %w", err)
		}
		fmt.Println("Category  Name  Shortcut key\n")
		for _, shortcut := range shortcuts {
			fmt.Printf("%v \t %v \t %v \n", shortcut.Category, shortcut.Name, shortcut.ShortcutKey)
		}
	}

	if *name != "" {
		shortcuts, err := shortcuts.GetShortcuts()
		if err != nil {
			return fmt.Errorf("failed to acquire the existing list: %w", err)
		}
		name := *name
		fmt.Println("Category  Name  Shortcut key")
		//TODO: How to remove new line betwewn the print above and the print below
		for _, shortcut := range shortcuts {
			// TODO: Add category filter here too
			if strings.Contains(shortcut.Name, name) {
				fmt.Printf("%v\t %v\t %v\n", shortcut.Category, shortcut.Name, shortcut.ShortcutKey)
			}
		}
	}
	return nil
}

func ValidateNewShortcutKey(category *string, name *string, key *string) error {
	if *category == "" || *name == "" || *key == "" {
		return errors.New("category, name and shortcut key are required to register new shortcut key. Check \"speep help\" for the usage.")
	}
	return nil
}

func handleAdd(addCmd *flag.FlagSet, category *string, name *string, newShortcut *string) error {

	err := ValidateNewShortcutKey(category, name, newShortcut)
	if err != nil {
		return fmt.Errorf("%s", err)
	}

	var allShortcuts []shortcuts.Shortcut
	shortcut := shortcuts.Shortcut{
		Category:    *category,
		Name:        *name,
		ShortcutKey: *newShortcut,
	}

	err = shortcuts.CheckNameDuplication(category, name)
	if err != nil {
		return fmt.Errorf("the name \"%v\" for %v %w", *name, *category, err)
	}

	allShortcuts, _ = shortcuts.GetShortcuts()
	allShortcuts = append(allShortcuts, shortcut)
	err = shortcuts.SaveShortcuts(allShortcuts)
	if err != nil {
		return fmt.Errorf("failed to save the updated list: %w", err)
	}

	fmt.Printf("New shortcut \"%v\" for %v successfully registered", *name, *category)
	return nil
}

func handleDelete(deleteCmd *flag.FlagSet, all *bool, allShort *bool, category *string, name *string) error {

	if !*all && !*allShort {
		if *category == "" && *name == "" {
			return errors.New("category and name are required to delete a shortcut key")
		} else if *category == "" || *name == "" {
			fmt.Println("Flag and/or argument for the command delete missing")
			fmt.Print(handlehelp(usage))
			return nil
		}
	}

	if *all || *allShort {
		// Delete full shortcut list
		err := shortcuts.DeleteShortcuts()
		if err != nil {
			return fmt.Errorf("failed to delete the shortcut list: %w", err)
		}
		fmt.Println("Shortcut list deleted")
	}

	if *category != "" && *name != "" {
		err := shortcuts.DeleteShortcut(*category, *name)
		if err != nil {
			return fmt.Errorf("failed to remove an item from the list: %w", err)
		}
		fmt.Printf("Shortcut %v for %v successfully removed from the Shortcut key list", *name, *category)
		// Add return error of the target not exist if the name doesn't exist in the list
	}

	return nil
}

func handlehelp(usage string) error {
	return errors.New(usage)
}
