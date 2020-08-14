package file

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
)

func (m *manager) CopyFile(srcPath string, targetPath string) error {
	newFile, err := m.create(targetPath)
	if err != nil {
		return fmt.Errorf("could not create file at %s - %w", targetPath, err)
	}

	defer func() {
		closeErr := newFile.Close()
		if err == nil && closeErr != nil {
			err = fmt.Errorf("error copying file from %s to %s - %w", srcPath, targetPath, closeErr)
		}
	}()

	input, err := m.open(srcPath)

	defer func() {
		closeErr := input.Close()
		if err == nil && closeErr != nil {
			err = fmt.Errorf("error copying file from %s to %s - %w", srcPath, targetPath, closeErr)
		}
	}()

	if err != nil {
		return fmt.Errorf("error copying file from %s to %s - %w", srcPath, targetPath, err)
	}

	_, err = io.Copy(newFile, input)
	if err != nil {
		return fmt.Errorf("error copying file from %s to %s - %w", srcPath, targetPath, err)
	}

	return nil
}

func (m *manager) CopyDirectory(srcPath string, targetPath string) error {
	contents, err := m.readDir(srcPath)
	if err != nil {
		return fmt.Errorf("could not read the contents of directory %s - %w", srcPath, err)
	}

	for _, c := range contents {
		var isSymLink bool
		from := filepath.Join(srcPath, c.Name())
		to := filepath.Join(targetPath, c.Name())
		info, err := os.Stat(from)
		if err != nil {
			return fmt.Errorf("could not get information about %s - %w", from, err)
		}

		stat, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("could not get stats for %s", from)
		}

		switch info.Mode() & os.ModeType {
		case os.ModeSymlink:
			isSymLink = true
			err = m.CopySymlink(from, to)
			if err != nil {
				return fmt.Errorf("could not copy symlink from %s to %s - %w", from, to, err)
			}
		case os.ModeDir:
			err = m.CreateDirectory(to, SetPermissions(0700))
			if err != nil {
				return fmt.Errorf("could not create directory at %s - %w", to, err)
			}
		default:
			err = m.CopyFile(from, to)
			if err != nil {
				return fmt.Errorf("could not copy file %s to %s - %w", from, to, err)
			}
		}

		err = os.Lchown(from, int(stat.Uid), int(stat.Gid))
		if err != nil {
			return fmt.Errorf("could not change ownership - %w", err)
		}

		if !isSymLink {
			err = os.Chmod(to, c.Mode())
			if err != nil {
				return fmt.Errorf("could not set correct permissions on the copied file - %w", err)
			}
		}
	}

	return nil
}

func (m *manager) CopySymlink(srcPath string, targetPath string) error {
	return nil
}
