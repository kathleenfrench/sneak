package htb

import (
	"fmt"
	"io/ioutil"

	"github.com/kathleenfrench/common/fs"
	"github.com/kathleenfrench/sneak/internal/config"
	"github.com/kathleenfrench/sneak/pkg/file"
)

func saveNoteFile(boxName string, note string) error {
	notesPath := fmt.Sprintf("%s/%s", config.GetNotesDirectory(), boxName)

	// create note directory for that box if it doesn't exist
	err := fs.CreateDir(notesPath)
	if err != nil {
		return fmt.Errorf("there was an error creating the notes directory for %s - %w", boxName, err)
	}

	notesFilePath := fmt.Sprintf("%s/main.md", notesPath)
	err = ioutil.WriteFile(notesFilePath, []byte(note), 0644)
	if err != nil {
		return err
	}

	return nil
}

func checkForNoteFile(boxName string) (string, error) {
	notesPath := fmt.Sprintf("%s/%s", config.GetNotesDirectory(), boxName)
	notesFilePath := fmt.Sprintf("%s/main.md", notesPath)

	// if the notes file already exists, read the text from the file and return it as a string to set a s adefault
	if fs.FileExists(notesFilePath) {
		note, err := ioutil.ReadFile(notesFilePath)
		if err != nil {
			return "", err
		}

		return string(note), nil
	}

	// create the directory if it doesn't exist yet
	err := fs.CreateDir(notesPath)
	if err != nil {
		return "", err
	}

	err = fs.CreateFile(notesFilePath)
	if err != nil {
		return "", err
	}

	fm := file.NewManager()
	err = fm.AppendToFile(notesFilePath, []byte(fmt.Sprintf("## %s\n\n", boxName)))
	if err != nil {
		return "", err
	}

	return "", nil
}
