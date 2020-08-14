package file

import "fmt"

func (m *manager) Write(path string, data []byte, ps ...PermissionSetter) error {
	// give user read/write by default
	perms := setDefaults(0600)
	for _, p := range ps {
		p(&perms)
	}

	fileExists, err := m.FileExists(path)
	switch {
	case err != nil:
		return fmt.Errorf("there was an error writing to %s - %w", path, err)
	case fileExists:
		isDirectory, dirErr := m.IsDirectory(path)
		switch {
		case dirErr != nil:
			return fmt.Errorf("there was an error writing to %s - %w", path, dirErr)
		case isDirectory:
			return fmt.Errorf("%s is a directory", path)
		}
	case !fileExists:
		parentDirectory := m.dirpath(path)
		// give user read/write/execute permissions on the directory by default
		parentDirectoryErr := m.mkdir(parentDirectory, 0700)
		if parentDirectoryErr != nil {
			return fmt.Errorf("could not recursively create the directories for %s - %w", path, err)
		}
	}

	err = m.write(path, data, perms.mode)
	if err != nil {
		return fmt.Errorf("cannot write to file %s - %w", path, err)
	}

	return nil
}
