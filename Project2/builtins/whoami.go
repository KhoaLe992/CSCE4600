package builtins

import (
	"fmt"
	"io"
	"os/user"
)

// WhoAmI prints the current username.
func WhoAmI(w io.Writer) error {
	u, err := user.Current()
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, u.Username)
	return err
}
