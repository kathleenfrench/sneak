package files

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"

	"github.com/phayes/permbits"
)

// WriteOptions are options to set when writing files
type WriteOptions struct {
}

// Manager is an interface for methods that determine information about files
type Manager interface {
	// FileExists(path string) bool
	// DirectoryExists(path string) bool
	// CreateDir(path string) error
	// IsFile(path string) (bool, error)
	// IsDirectory(path string) (bool, error)
	Readable(path string, uid int, gid int) (bool, error)
	Executable(path string, uid int, gid int) (bool, error)
	Writable(path string, uid int, gid int) (bool, error)
	// Delete(path string) error
	// ReadFile(path string) ([]byte, error)
	GetOwnership(path string) (userID int, groupID int, err error)
	GetUserPerms(path string) (read bool, write bool, execute bool, err error)
	GetGroupPerms(path string) (read bool, write bool, execute bool, err error)
	GetOthersPerms(path string) (read bool, write bool, execute bool, err error)
	SetOwnership(path string, uid int, gid int) error
	// SetUserPerms(path string, mod os.FileMode) error
	// Append(path string, lines []string, opts ...WriteOptions) error
	// Touch(path string, opts ...WriteOptions) error
	// Write(path string, data []byte, opts ...WriteOptions) error
	// CopyFile(srcPath string, targetPath string) error
	// CopyDirectory(srcPath string, targetPath string) error
	// CopySymlink(srcPath string, targetPath string) error
}

type manager struct {
	stat            func(name string) (os.FileInfo, error)
	read            func(filename string) ([]byte, error)
	rm              func(path string) error
	open            func(path string) (*os.File, error)
	create          func(path string) (*os.File, error)
	chmod           func(path string, mod os.FileMode) error
	chown           func(path string, uid int, gid int) error
	write           func(path string, data []byte, perm os.FileMode) error
	mkdir           func(path string, perm os.FileMode) error
	dirname         func(path string) string
	dirinfo         func(path string) ([]os.FileInfo, error)
	symlink         func(oldPath string, newPath string) error
	readlink        func(path string) (string, error)
	doesNotExistErr func(err error) bool
}

// NewManager instantiates a new interface for methods that handle working with files on the file system
func NewManager() Manager {
	return &manager{
		stat:            os.Stat,
		read:            ioutil.ReadFile,
		rm:              os.RemoveAll,
		open:            os.Open,
		create:          os.Create,
		chmod:           os.Chmod,
		chown:           os.Chown,
		write:           ioutil.WriteFile,
		mkdir:           os.MkdirAll,
		dirname:         filepath.Dir,
		dirinfo:         ioutil.ReadDir,
		symlink:         os.Symlink,
		readlink:        os.Readlink,
		doesNotExistErr: os.IsNotExist,
	}
}

type accessible string

const (
	readable   accessible = "readable"
	writable   accessible = "writable"
	executable accessible = "executable"
)

func (m *manager) GetOwnership(path string) (userID int, groupID int, err error) {
	fi, err := m.stat(path)
	if err != nil {
		return 0, 0, err
	}

	stat, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return 0, 0, nil
	}

	return int(stat.Uid), int(stat.Gid), nil
}

func (m *manager) SetOwnership(path string, uid int, gid int) error {
	return m.chown(path, uid, gid)
}

func (m *manager) GetUserPerms(path string) (read bool, write bool, execute bool, err error) {
	fi, err := m.stat(path)
	if err != nil {
		return false, false, false, permissionsGetError(err)
	}

	perms := permbits.FileMode(fi.Mode())
	return perms.UserRead(), perms.UserWrite(), perms.UserExecute(), nil
}

func (m *manager) GetGroupPerms(path string) (read bool, write bool, execute bool, err error) {
	fi, err := m.stat(path)
	if err != nil {
		return false, false, false, permissionsGetError(err)
	}

	perms := permbits.FileMode(fi.Mode())
	return perms.GroupRead(), perms.GroupWrite(), perms.GroupExecute(), nil
}

func (m *manager) GetOthersPerms(path string) (read bool, write bool, execute bool, err error) {
	fi, err := m.stat(path)
	if err != nil {
		return false, false, false, permissionsGetError(err)
	}

	perms := permbits.FileMode(fi.Mode())
	return perms.OtherRead(), perms.OtherWrite(), perms.OtherExecute(), nil
}

func permissionsGetError(err error) error {
	return fmt.Errorf("could not determine permissions: %w", err)
}

func (m *manager) Readable(path string, uid int, gid int) (bool, error) {
	return m.canAccess(path, uid, gid, readable)
}

func (m *manager) Executable(path string, uid int, gid int) (bool, error) {
	return m.canAccess(path, uid, gid, executable)
}

func (m *manager) Writable(path string, uid int, gid int) (bool, error) {
	return m.canAccess(path, uid, gid, writable)
}

func (m *manager) canAccess(path string, uid int, gid int, accessibility accessible) (bool, error) {
	ownerUID, ownerGID, err := m.GetOwnership(path)
	if err != nil {
		return false, err
	}

	hasAccess := false

	switch accessibility {
	case readable:
		hasAccess, _, _, err = m.GetOthersPerms(path)
	case writable:
		_, hasAccess, _, err = m.GetOthersPerms(path)
	case executable:
		_, _, hasAccess, err = m.GetOthersPerms(path)
	}

	if err != nil {
		return false, err
	}

	if hasAccess {
		return true, nil
	}

	if gid == ownerGID {
		hasAccess := false

		switch accessibility {
		case readable:
			hasAccess, _, _, err = m.GetGroupPerms(path)
		case writable:
			_, hasAccess, _, err = m.GetGroupPerms(path)
		case executable:
			_, _, hasAccess, err = m.GetGroupPerms(path)
		}

		if err != nil {
			return false, err
		}

		if hasAccess {
			return true, nil
		}
	}

	if uid == ownerUID {
		hasAccess := false

		switch accessibility {
		case readable:
			hasAccess, _, _, err = m.GetUserPerms(path)
		case writable:
			_, hasAccess, _, err = m.GetUserPerms(path)
		case executable:
			_, _, hasAccess, err = m.GetUserPerms(path)
		}

		if err != nil {
			return false, err
		}

		if hasAccess {
			return true, nil
		}
	}

	return false, nil
}
