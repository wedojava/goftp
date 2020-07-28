package ftp

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

// retr downloading files from the server
// the get command you enter in your client is sent to the server as RETR
func (c *Conn) retr(args []string) {
	if len(args) != 1 {
		c.respond(status501)
		return
	}

	// build the path for the file we want.
	path := filepath.Join(c.rootDir, c.workDir, args[0])
	file, err := os.Open(path)
	if err != nil {
		log.Print(err)
		c.respond(status550)
	}
	c.respond(status150)

	// If the file does exist, we open a new data connection with (c *Conn) dataConnect
	dataConn, err := c.dataConnect()
	if err != nil {
		log.Print(err)
		c.respond(status425)
	}
	defer dataConn.Close()

	// copy file content to dataConn
	// This works well for small files, but files on the gigabyte scale could kill the server, especially when you consider that we’re serving multiple clients concurrently.
	_, err = io.Copy(dataConn, file)
	if err != nil {
		log.Print(err)
		c.respond(status426)
		return
	}
	io.WriteString(dataConn, c.EOL())
	c.respond(status226)
}
