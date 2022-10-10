package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type shortcut struct {
	Name        string
	ShortcutKey string
}

func getShortcuts() (shortcuts []shortcut) {
	// Check if the json file exists and create new one if it doesn't exist
	_, err := os.Open("./shortcuts.json")
	if err != nil {
		err = os.WriteFile("./shortcuts.json", nil, 0644)
		if err != nil {
			panic(err)
		}
	}

	fileBytes, err := os.ReadFile("./shortcuts.json")
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(fileBytes, &shortcuts)
	if err != nil {
		fmt.Println(err)
	}

	return shortcuts
}

func saveShortcuts(shortcuts []shortcut) {
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
		var mapStruct map[string]json.RawMessage
		if err := json.Unmarshal([]byte(shortcut), &mapStruct); err != nil {
			panic(err)
		}
		if strings.Contains(string(mapStruct[string]), name) {
			delete(shortcuts, name)
		}
	}
}
