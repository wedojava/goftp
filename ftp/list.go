package ftp

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// The client-side ls path/to/file command reaches the server as LIST path/to/file, and it will come as no surprise that we have a (c *Conn) list handler function to match.
func (c *Conn) list(args []string) {
	target := filepath.Join(c.rootDir, c.workDir)
	if len(args) > 0 {
		target = filepath.Join(target, args[0])
	}
	f, err := os.Open(target)
	if err != nil {
		log.Print(err)
		c.respond(status550)
		return
	}
	c.respond(status150)
	w, err := c.dataConnect()
	if err != nil {
		log.Print(err)
		c.respond(status425)
		return
	}
	defer w.Close()
	stat, err := f.Stat()
	if err != nil {
		log.Print(err)
		c.respond(status450)
		return
	}
	if stat.IsDir() {
		filenames, err := f.Readdirnames(0)
		if err != nil {
			log.Print(err)
			c.respond(status550)
			return
		}
		for _, filename := range filenames {
			_, err = fmt.Fprint(w, filename, c.EOL())
			if err != nil {
				log.Print(err)
				c.respond(status426)
				return
			}
		}
		c.respond(status226)
		return
	}
	rel, err := filepath.Rel(c.rootDir, target)
	if err != nil {
		log.Print(err)
		c.respond(status550)
		return
	}
	_, err = fmt.Fprint(w, rel, c.EOL())
	if err != nil {
		log.Print(err)
		c.respond(status426)
		return
	}
	c.respond(status226)
	return
}
