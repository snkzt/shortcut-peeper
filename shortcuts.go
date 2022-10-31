package shortcuts

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

func getEnvPath() string {
	fileStoragePath := os.Getenv("XDG_CONFIG_HOME")
	if fileStoragePath == "" {
		fileStoragePath := os.Getenv("HOME")
		fileStoragePath += "/.config/speep/shortcuts.json"
		return fileStoragePath
	}
	fileStoragePath += "/speep/shortcuts.json"
	return fileStoragePath
}

func GetShortcuts() ([]Shortcut, error) {
	var shortcuts []Shortcut
	// Check if the json file exists and create new one if it doesn't exist
	fileStoragePath := getEnvPath()
	_, err := os.Open(fileStoragePath)
	if err != nil {
		return nil, fmt.Errorf("no list exists: %w", err)
	}

	fileBytes, err := os.ReadFile(fileStoragePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read the existing list: %w", err)
	}

	// Set the decoded data of fileBytes to shortcuts struct
	err = json.Unmarshal(fileBytes, &shortcuts)
	if err != nil {
		return nil, fmt.Errorf("Failed to store byte data to shortcuts struct: %w", err)
	}

	return shortcuts, nil
}

func SaveShortcuts(shortcuts []Shortcut) error {
	shortcutBytes, err := json.Marshal(shortcuts)
	if err != nil {
		return fmt.Errorf("failed to store list as JSON : %w", err)
	}

	// Write the new shortcut into the json file and create one if the file not exists
	fileStoragePath := getEnvPath()
	err = os.WriteFile(fileStoragePath, shortcutBytes, 0644)
	if err != nil {
		return fmt.Errorf("Failed to save the new shortcut: %w", err)
	}
	return nil
}

func DeleteShortcuts() error {
	// Remove shortcuts.json file
	fileStoragePath := getEnvPath()
	err := os.Remove(fileStoragePath)
	if err != nil {
		return fmt.Errorf("failed to delete the existing list: %w", err)
	}
	return nil
}

func DeleteShortcut(name string) error {
	// Get existing shortcut key list
	shortcuts, err := GetShortcuts()
	if err != nil {
		return fmt.Errorf("failed to acquire existing list: %w", err)
	}

	for _, shortcut := range shortcuts {
		// Convert shortcuts struct into map to delete a specific item
		var mapStruct map[string]interface{}
		shortcutByte, _ := json.Marshal(shortcut)
		if err := json.Unmarshal(shortcutByte, &mapStruct); err != nil {
			return fmt.Errorf("Failed to convert shortcuts struct into map: %w", err)
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

	err = SaveShortcuts(shortcuts)
	if err != nil {
		return fmt.Errorf("failed to save the list after removing an item: %w", err)
	}
	return nil
}
