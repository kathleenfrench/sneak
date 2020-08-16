package file

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// Manager is an interface for methods that work with files on the file system
type Manager interface {
	// files
	Touch(path string, p ...PermissionSetter) error
	Write(path string, data []byte, p ...PermissionSetter) error
	FileExists(path string) (bool, error)
	FilepathExists(path string) (bool, error)
	IsFile(path string) (bool, error)
	CopyFile(srcPath string, targetPath string) error
	AppendToFile(path string, data []byte, ps ...PermissionSetter) error
	ReadFile(path string) ([]byte, error)
	// directories
	DirectoryExists(path string) (bool, error)
	CreateDirectory(path string, p ...PermissionSetter) error
	IsDirectory(path string) (bool, error)
	CopyDirectory(srcPath string, targetPath string) error
	// symlinks
	CopySymlink(srcPath string, targetPath string) error
	// removes a file or directory of the provided path
	Remove(path string) error
	// current working directories
	CWD() (string, error)
	Basename(path string) string
}

type manager struct {
	mkdir       func(path string, perm os.FileMode) error
	rm          func(path string) error
	write       func(name string, data []byte, mode os.FileMode) error
	dirpath     func(path string) string
	readDir     func(name string) ([]os.FileInfo, error)
	read        func(name string) ([]byte, error)
	stat        func(name string) (os.FileInfo, error)
	open        func(name string) (*os.File, error)
	openFile    func(name string, flag int, mod os.FileMode) (*os.File, error)
	create      func(name string) (*os.File, error)
	notExistErr func(err error) bool
	chmod       func(name string, mod os.FileMode) error
	readlink    func(name string) (string, error)
	symlink     func(from, to string) error
	cwd         func() (string, error)
	basename    func(path string) string
}

// NewManager instantiates a new manager interface for working with files
func NewManager() Manager {
	return &manager{
		mkdir:       os.MkdirAll,
		rm:          os.RemoveAll,
		write:       ioutil.WriteFile,
		dirpath:     filepath.Dir,
		readDir:     ioutil.ReadDir,
		read:        ioutil.ReadFile,
		stat:        os.Stat,
		open:        os.Open,
		openFile:    os.OpenFile,
		create:      os.Create,
		notExistErr: os.IsNotExist,
		chmod:       os.Chmod,
		readlink:    os.Readlink,
		symlink:     os.Symlink,
		cwd:         os.Getwd,
		basename:    filepath.Base,
	}
}
