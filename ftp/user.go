package ftp

import (
	"fmt"
	"strings"
)

// When ftp.Serve encounters the USER command, it passes the arguments to a handler method of ftp.Conn: (c *Conn) user. From there, it’s easy to imagine how the details would be checked against a database of known users and their access permissions. Once authenticated, the username and privileges could be stored as additional fields on the ftp.Conn instance.
// In our case, we’re so blasé about network security that we simply echo the username back to the client with the appropriate success code: 230 User %s logged in, proceed.
func (c *Conn) user(args []string) {
	c.respond(fmt.Sprintf(status230, strings.Join(args, " ")))
}
