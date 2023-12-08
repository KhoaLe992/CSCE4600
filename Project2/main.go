package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/KhoaLe992/CSCE4600/Project2/builtins"
)

// Echo prints the specified text to the console.
func Echo(w io.Writer, args ...string) error {
	text := strings.Join(args, " ")
	_, err := fmt.Fprintln(w, text)
	return err
}

// PrintWorkingDirectory prints the current working directory.
func PrintWorkingDirectory(w io.Writer) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, wd)
	return err
}

// ListDirectory lists the files and directories in the specified directory.
func ListDirectory(w io.Writer, args ...string) error {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		_, err := fmt.Fprintln(w, file.Name())
		if err != nil {
			return err
		}
	}
	return nil
}

// MakeDirectory creates a new directory with the specified name.
func MakeDirectory(args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("mkdir: missing operand")
	}
	dir := args[0]
	return os.Mkdir(dir, os.ModePerm)
}

// WhoAmI prints the current username.
func WhoAmI(w io.Writer) error {
	u, err := user.Current()
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, u.Username)
	return err
}

func main() {
	exit := make(chan struct{}, 2) // buffer this so there's no deadlock.
	runLoop(os.Stdin, os.Stdout, os.Stderr, exit)
}

func runLoop(r io.Reader, w, errW io.Writer, exit chan struct{}) {
	var (
		input    string
		err      error
		readLoop = bufio.NewReader(r)
	)
	for {
		select {
		case <-exit:
			_, _ = fmt.Fprintln(w, "exiting gracefully...")
			return
		default:
			if err := printPrompt(w); err != nil {
				_, _ = fmt.Fprintln(errW, err)
				continue
			}
			if input, err = readLoop.ReadString('\n'); err != nil {
				_, _ = fmt.Fprintln(errW, err)
				continue
			}
			if err = handleInput(w, input, exit); err != nil {
				_, _ = fmt.Fprintln(errW, err)
			}
		}
	}
}

func printPrompt(w io.Writer) error {
	// Get current user.
	// Don't prematurely memoize this because it might change due to `su`?
	u, err := user.Current()
	if err != nil {
		return err
	}
	// Get current working directory.
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// /home/User [Username] $
	_, err = fmt.Fprintf(w, "%v [%v] $ ", wd, u.Username)

	return err
}

func handleInput(w io.Writer, input string, exit chan<- struct{}) error {
	// Remove trailing spaces.
	input = strings.TrimSpace(input)

	// Split the input separate the command name and the command arguments.
	args := strings.Split(input, " ")
	name, args := args[0], args[1:]

	// Check for built-in commands.
	// New builtin commands should be added here. Eventually this should be refactored to its own func.
	switch name {
	case "cd":
		return builtins.ChangeDirectory(args...)
	case "env":
		return builtins.EnvironmentVariables(w, args...)
	case "exit":
		exit <- struct{}{}
		return nil
	case "echo":
		return builtins.Echo(w, args...)
	case "pwd":
		return builtins.PrintWorkingDirectory(w)
	case "ls":
		return builtins.ListDirectory(w, args...)
	case "mkdir":
		return builtins.MakeDirectory(args...)
	case "whoami":
		return builtins.WhoAmI(w)
	}

	return executeCommand(name, args...)
}

func executeCommand(name string, arg ...string) error {
	// Otherwise prep the command
	cmd := exec.Command(name, arg...)

	// Set the correct output device.
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	// Execute the command and return the error.
	return cmd.Run()
}
