package builtins

import (
	"fmt"
	"os"
)

// MakeDirectory creates a new directory with the specified name.
func MakeDirectory(args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("mkdir: missing operand")
	}
	dir := args[0]
	return os.Mkdir(dir, os.ModePerm)
}
