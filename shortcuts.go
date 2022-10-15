package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Shortcut struct {
	Name        string
	ShortcutKey string
}

func getShortcuts() (shortcuts []Shortcut) {
	// Check if the json file exists and create new one if it doesn't exist
	_, err := os.Open("./shortcuts.json")
	if err != nil {
		err = os.WriteFile("./shortcuts.json", nil, 0644)
		if err != nil {
			fmt.Println(err)
		}
	}

	fileBytes, err := os.ReadFile("./shortcuts.json")
	if err != nil {
		fmt.Println(err)
	}

	// Set the decoded data of fileBytes to shortcuts struct
	err = json.Unmarshal(fileBytes, &shortcuts)
	if err != nil {
		fmt.Println(err)
	}

	return shortcuts
}

func saveShortcuts(shortcuts []Shortcut) {
	shortcutBytes, err := json.Marshal(shortcuts)
	if err != nil {
		fmt.Println(err)
	}

	// Write the new shortcut into the json file and create one if the file not exists
	err = os.WriteFile("./shortcuts.json", shortcutBytes, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func deleteShortcuts() {
	// Remove shortcuts.json file
	err := os.Remove("./shortcuts.json")
	if err != nil {
		fmt.Println(err)
	}
}

func deleteShortcut(name string) {
	// Get existing shortcut key list
	shortcuts := getShortcuts()
	for _, shortcut := range shortcuts {
		var mapStruct map[string]interface{}
		shortcutByte, _ := json.Marshal(shortcut)
		if err := json.Unmarshal(shortcutByte, &mapStruct); err != nil {
			fmt.Println(err)
		}
		for structName, _ := range mapStruct {
			if strings.Contains(structName, name) {
				delete(mapStruct, name)
			}
			fmt.Printf("%v successfully removed", name)
		}

		shortcutStruct, _ := json.Marshal(mapStruct)
		if err := json.Unmarshal(shortcutStruct, &shortcuts); err != nil { // will this be append to existing structs[]?
			fmt.Println(err)
		}
		fmt.Printf("No \"%v\" exists", name)
	}
	saveShortcuts(shortcuts)
}
