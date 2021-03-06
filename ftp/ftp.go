package ftp

import (
	"bufio"
	"log"
	"strings"
)

// Serve scans incoming requests for valid commands and routes them to handler functions.
func Serve(c *Conn) {
	c.respond(status220)
	s := bufio.NewScanner(c.conn)
	for s.Scan() {
		input := strings.Fields(s.Text())
		if len(input) == 0 {
			continue
		}
		command, args := input[0], input[1:]
		// log.Printf("<< %s %v", command, args)
		log.Printf("<< %s %v %s", command, args, c.dataPort.toAddress())

		switch strings.ToUpper(command) {
		case "PASV":
			c.pasv(args)
		case "OPTS":
			c.opts(args)
		case "PWD":
			c.pwd(args)
		case "SYST":
			c.syst(args)
		case "CWD": // cd
			c.cwd(args)
		case "LIST": //ls
			c.list(args)
		case "PORT":
			c.port(args)
		case "USER":
			c.user(args)
		case "QUIT": // close
			c.respond(status221)
			return
		case "RETR": // get
			c.retr(args)
		case "TYPE":
			// when the TYPE command is identified, its argument is passed to setDataType
			c.setDataType(args)
			// then the connection’s dataType is updated
		default:
			c.respond(status502)
		}
		// Cleanup PASV listeners if they go unused.
		if command != "PASV" && c.pasvListener != nil {
			c.pasvListener.Close()
			c.pasvListener = nil
		}
	}
	if s.Err() != nil {
		log.Print(s.Err())
	}
}
