package lxops_incus

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	incus "github.com/lxc/incus/v6/client"
	"github.com/lxc/incus/v6/shared/api"
)

type InstanceConfigurer struct {
	Server      incus.InstanceServer
	Log         io.Writer
	instance    string
	createdDirs map[string]struct{}
}

// NewInstanceConfigurer creates a BaseConfigurer for an instance.
// The configurer should not be reused for other instances.
func NewInstanceConfigurer(server incus.InstanceServer, instance string) *InstanceConfigurer {
	t := &InstanceConfigurer{Server: server, instance: instance}
	t.createdDirs = make(map[string]struct{})
	return t
}

func (t *InstanceConfigurer) SetLogWriter(w io.Writer) {
	t.Log = w
}

func (t *InstanceConfigurer) RunScript(script string) error {
	return t.exec(script, "/bin/sh")

}

func (t *InstanceConfigurer) RunCommand(args ...string) error {
	return t.exec("", args...)

}

func (t *InstanceConfigurer) FileExists(file string) (bool, error) {
	reader, _, err := t.Server.GetInstanceFile(t.instance, file)
	if err != nil {
		/*
			fmt.Printf("%v (%T)\n", err, err)
			switch e := err.(type) {
			case api.StatusError:
				fmt.Printf("status: %d\n", e.Status())
			}
		*/
		// The error could be:
		// 	Not Found
		// 	Instance not found
		//  ... and possibly others if the communication with the server fails
		// Rather than trying to distinguish between different types of errors,
		// return that the file does not exist.
		return false, nil
	}
	reader.Close()
	return true, nil
}

func exitCode(serverOp incus.Operation) (int, error) {
	op := serverOp.Get()
	/*
		data, _ := json.MarshalIndent(op, "", " ")
		os.Stdout.Write(data)
	*/
	returnValue := op.Metadata["return"]
	if returnValue == nil {
		return 1, fmt.Errorf("missing return code")
	}
	switch code := returnValue.(type) {
	case float64:
		return int(code), nil
	case int:
		return code, nil
	default:
		return 1, fmt.Errorf("unexpected return type: %T", returnValue)
	}
}

func (t *InstanceConfigurer) exec(input string, execArgs ...string) error {
	if len(execArgs) == 0 {
		return fmt.Errorf("empty command")
	}
	var post api.InstanceExecPost
	post.Command = execArgs
	post.WaitForWS = true

	var args incus.InstanceExecArgs
	if t.Log != nil {
		args.Stderr = NopWriteCloser(t.Log)
	} else {
		args.Stderr = NopWriteCloser(os.Stderr)
	}
	if t.Log != nil {
		args.Stdout = NopWriteCloser(t.Log)
	}

	if input != "" {
		args.Stdin = io.NopCloser(strings.NewReader(input))
	}
	op, err := t.Server.ExecInstance(t.instance, post, &args)
	if err != nil {
		return fmt.Errorf("%s: %w", t.instance, err)
	}
	err = op.Wait()
	if err != nil {
		return fmt.Errorf("%s: %w", t.instance, err)
	}
	exitCode, err := exitCode(op)
	if err != nil {
		return err
	}
	if exitCode != 0 {
		return fmt.Errorf("exit code: %d", exitCode)
	}
	return nil
}

func (t *InstanceConfigurer) ensureDirExists(dir string) error {
	if dir == "/" || dir == "." {
		return nil
	}
	_, exists := t.createdDirs[dir]
	if exists {
		return nil
	}
	err := t.exec("", "mkdir", "-p", dir)
	if err != nil {
		return err
	}
	for d := dir; !(d == "." || d == "/"); d = filepath.Dir(d) {
		t.createdDirs[d] = struct{}{}
	}
	return nil
}

func (t *InstanceConfigurer) writeOrAppendFile(path string, data []byte, perm fs.FileMode, writeMode string) error {
	var args incus.InstanceFileArgs
	args.Mode = int(perm)
	args.WriteMode = writeMode
	args.Content = bytes.NewReader(data)
	err := t.Server.CreateInstanceFile(t.instance, path, args)
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	return nil
}

func (t *InstanceConfigurer) WriteFile(path string, data []byte, perm fs.FileMode) error {
	dir := filepath.Dir(path)
	err := t.ensureDirExists(dir)
	if err != nil {
		return err
	}
	return t.writeOrAppendFile(path, data, perm, "overwrite")
}

// AppendFile appends to an existing file
// The file must already exist
// If we want to make append work with a new file,
// we would have to make more calls to touch the file first.
// Need to check what cloud-init does
func (t *InstanceConfigurer) AppendFile(path string, data []byte, perm fs.FileMode) error {
	return t.writeOrAppendFile(path, data, perm, "append")
}

// NopWriteCloser returns a WriteCloser with a no-op Close method wrapping
// the provided Writer
func NopWriteCloser(w io.Writer) io.WriteCloser {
	return nopWriteCloser{w}
}

type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error { return nil }
