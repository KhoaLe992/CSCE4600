package builtins

import (
	"fmt"
	"io"
	"strings"
)

// Echo prints the specified text to the console.
func Echo(w io.Writer, args ...string) error {
	text := strings.Join(args, " ")
	_, err := fmt.Fprintln(w, text)
	return err
}
