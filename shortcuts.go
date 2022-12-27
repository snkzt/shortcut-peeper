package shortcuts

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Shortcut struct {
	Category    string
	Name        string
	ShortcutKey string
}

func createDirectory(fileStoragePath string) (string, error) {
	speepDir := fileStoragePath + "/speep"
	_, err := os.Stat(fileStoragePath)
	_, errS := os.Stat(speepDir)

	switch {
	//case !strings.Contains(fileStoragePath, "/.config"):
	case os.IsNotExist(err):
		err = os.Mkdir(fileStoragePath, 0750)
		if err != nil {
			return "", fmt.Errorf("failed to create directory \".config\": %w", err)
		}
		fallthrough
	case os.IsNotExist(errS):
		err = os.Mkdir(speepDir, 0750)
		if err != nil {
			return "", fmt.Errorf("failed to create directory \".config/speep\": %w", err)
		}
		return speepDir, nil
	default:
		return speepDir, nil
	}
}

func getEnvPath() (string, error) {
	fileStoragePath := os.Getenv("XDG_CONFIG_HOME")
	if fileStoragePath == "" {
		fileStoragePath = os.Getenv("HOME")
	}

	fileStoragePath += "/.config"
	fileStoragePath, err := createDirectory(fileStoragePath)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve env path: %w", err)
	}

	fileStoragePath += "/shortcuts.json"
	return fileStoragePath, nil
}

func GetShortcuts() ([]Shortcut, error) {
	var shortcuts []Shortcut

	// Check if the json file exists and create new one if it doesn't exist
	fileStoragePath, err := getEnvPath()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve shortcut list: %w", err)
	}

	_, err = os.Open(fileStoragePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open shortcut list: %w", err)
	}

	fileBytes, err := os.ReadFile(fileStoragePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read the shortcut list: %w", err)
	}

	// Set the decoded data of fileBytes to shortcuts struct
	err = json.Unmarshal(fileBytes, &shortcuts)
	if err != nil {
		return nil, fmt.Errorf("failed to store byte data to shortcuts struct: %w", err)
	}
	return shortcuts, nil
}

func CheckNameDuplication(category *string, name *string) error {

	shortcutList, _ := GetShortcuts()
	for _, shortcut := range shortcutList {
		if shortcut.Category == *category {
			if shortcut.Name == *name {
				return errors.New("already registered")
			}
		}
	}
	return nil
}

func SaveShortcuts(shortcuts []Shortcut) error {
	shortcutBytes, err := json.Marshal(shortcuts)
	if err != nil {
		return fmt.Errorf("failed to store list as JSON : %w", err)
	}

	fileStoragePath, err := getEnvPath()
	if err != nil {
		return fmt.Errorf("failed to save shortcut list: %w", err)
	}

	// Write the new shortcut into the existing list
	// and create one if no list exists
	err = os.WriteFile(fileStoragePath, shortcutBytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to save the new shortcut: %w", err)
	}
	return nil
}

func DeleteShortcuts() error {
	// Remove shortcuts.json file
	fileStoragePath, err := getEnvPath()
	if err != nil {
		return fmt.Errorf("failed to delete shortcut list: %w", err)
	}

	err = os.Remove(fileStoragePath)
	if err != nil {
		return fmt.Errorf("failed to delete the existing list: %w", err)
	}
	return nil
}

func DeleteShortcut(category string, name string) error {
	// Get existing shortcut key list
	shortcuts, err := GetShortcuts()
	if err != nil {
		return fmt.Errorf("failed to acquire existing list: %w", err)
	}

	// // Convert shortcuts []Shortcut into map to range
	// // to take out each one of the items
	// shortcutsByte, _ := json.Marshal(shortcuts)
	// fmt.Printf("shortcutsByte: %T\n", shortcutsByte)
	// if err := json.Unmarshal(shortcutsByte, &shortcuts); err != nil {
	// 	return fmt.Errorf("failed to convert shortcuts struct into map: %w", err)
	// }

	for i, shortcut := range shortcuts {
		if shortcut.Category == category {
			if shortcut.Name == name {
				if (len(shortcuts) - 1) <= 0 {
					shortcuts = nil
					err = SaveShortcuts(shortcuts)
					if err != nil {
						return fmt.Errorf("failed to save the list after removing an item: %w", err)
					}
					return nil
				} else {
					shortcuts[i] = shortcuts[len(shortcuts)-1]
					shortcuts = shortcuts[:len(shortcuts)-1]
					err = SaveShortcuts(shortcuts)
					if err != nil {
						return fmt.Errorf("failed to save the list after removing an item: %w", err)
					}
					return nil
				}
			}
		}
	}
	return fmt.Errorf("the shortcut name %s for the category %s does not exist", name, category)
}
