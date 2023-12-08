package builtins

import (
	"fmt"
	"io"
	"os"
)

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
