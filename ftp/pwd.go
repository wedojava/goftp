package ftp

import "fmt"

func (c *Conn) pwd(args []string) {
	if len(args) > 0 {
		c.respond(status501)
	}
	c.respond(fmt.Sprintf("257 %q is current directory.", c.workDir))
}
