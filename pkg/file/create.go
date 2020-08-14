package file

import "fmt"

func (m *manager) CreateDirectory(path string, ps ...PermissionSetter) error {
	perms := setDefaults(0700)
	for _, p := range ps {
		p(&perms)
	}

	pathExists, err := m.FilepathExists(path)
	switch {
	case pathExists:
		isFile, fileErr := m.IsFile(path)
		switch {
		case isFile:
			return fmt.Errorf("a file already exists at %s - cannot add a directory", path)
		case fileErr != nil:
			return fmt.Errorf("could not determine the type for %s - %w", path, err)
		}

		return m.chmod(path, perms.mode)
	case !pathExists:
		mkdirErr := m.mkdir(path, perms.mode)
		if mkdirErr != nil {
			return fmt.Errorf("could not create directory at %s - %w", path, err)
		}
	case err != nil:
		return fmt.Errorf("could not create directory %s - %w", path, err)
	}

	return nil
}

func (m *manager) Touch(path string, ps ...PermissionSetter) error {
	return m.Write(path, nil, ps...)
}
