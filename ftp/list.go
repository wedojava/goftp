package ftp

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
)

// The client-side ls path/to/file command reaches the server as LIST path/to/file, and it will come as no surprise that we have a (c *Conn) list handler function to match.
func (c *Conn) list(args []string) {
	var target string
	if len(args) > 0 {
		target = filepath.Join(c.rootDir, c.workDir, args[0])
	} else {
		target = filepath.Join(c.rootDir, c.workDir)
	}
	files, err := ioutil.ReadDir(target)
	if err != nil {
		log.Print(err)
		c.respond(status550)
		return
	}
	c.respond(status150)
	// When sending anything other than statuses, the server must establish a second, temporary connection to the client, known as the data connection, or dataConn.
	// Moreover, the connection must be made to a specific port that the FTP client has selected in advance.
	// How is this achieved? Before sending the LIST directive to the server, the client sends another command behind the scenes: PORT. PORT has a six-byte argument, corresponding to the four parts of an IP address, plus two bytes to represent a port number up to five digits long, e.g., PORT [127,0,0,1,245,1]. Be aware that your client may use subtly different formats but, if your server has proper logging, you’ll soon spot the difference.
	dataConn, err := c.dataConnect()
	if err != nil {
		log.Print(err)
		c.respond(status425)
		return
	}
	defer dataConn.Close()

	for _, file := range files {
		_, err := fmt.Fprint(dataConn, file.Name(), c.EOL())
		if err != nil {
			log.Print(err)
			c.respond(status426)
		}
	}
	_, err = fmt.Fprintf(dataConn, c.EOL())
	if err != nil {
		log.Print(err)
		c.respond(status426)
	}
	c.respond(status226)
}
