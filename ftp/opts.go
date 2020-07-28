package ftp

import "strings"

func (c *Conn) opts(args []string) {
	if len(args) != 2 {
		c.respond(status550)
		return
	}
	if strings.ToUpper(args[0]) != "UTF8" {
		c.respond(status550)
	}
	if strings.ToUpper(args[1]) != "ON" {
		c.respond(status200)
	} else {
		c.respond(status550)
	}
}
