package srv

import (
	"io"
	"io/fs"
)

// InstanceConfigurer is the same interface as github.com/melato.org/cloudconfig.BaseConfigurer
type InstanceConfigurer interface {
	SetLogWriter(w io.Writer)
	// RunScript runs sh with the given input
	RunScript(input string) error

	// RunCommand runs program args[0], with args args
	RunCommand(args ...string) error

	// WriteFile writes a file, like os.WriteFile.  It should not try to create any directories.
	WriteFile(path string, data []byte, perm fs.FileMode) error

	// AppendFile appends to a file.  It should not try to create any directories.
	AppendFile(path string, data []byte, perm fs.FileMode) error

	// FileExists checks if the file exists.\
	// It should return an error only if an unexpected error occurs,
	// not if the file simply does not exist.
	FileExists(path string) (bool, error)
}
