package ftp

import (
	"log"
	"os"
	"path/filepath"
)

// cwd switch directory when you submit a command such as cd ../parent_folder to your FTP client, it sends that message to the server as CWD ../parent_folder. ftp.Serve passes the file path to cwd, a method of our proprietary ftp.Conn
func (c *Conn) cwd(args []string) {
	if len(args) != 1 {
		c.respond(status501)
		return
	}
	// After checking we have the correct number of arguments, we build an “absolute” path that uses the root directory of the server as its root. In its simplest form, it’s a matter of joining the path arg to the end of the current working directory, then joining the result to the end of rootDir.
	// TODO: In addition, continually appending new path args to an ever-growing workDir is dreadful memory management that will also slow file lookups. Explore how the standard library’s filepath.Clean can help solve this problem.
	workDir := filepath.Join(c.workDir, args[0])
	absPath := filepath.Join(c.rootDir, workDir)
	// Go’s os.Stat returns information about a target file plus an error, which is non-nil if the file can’t be accessed. If that’s the case, we log it and respond 550 Requested action not taken. File unavailable. Otherwise, we update the ftp.Conn’s workDir and respond with a 200 success message.
	_, err := os.Stat(absPath)
	if err != nil {
		log.Print(err)
		c.respond(status550)
		return
	}
	c.workDir = workDir
	// Before responding that the change was successful, however, we need to validate that the directory exists and is accessible. It’s possible that the program, or if you’ve extended the USER handler, this Conn’s user, doesn’t have permission to read the target directory.
	c.respond(status200)
	// If you’d like to challenge yourself by improving this naïve implementation, consider how you would prevent the user from accessing files above the server’s root directory. Currently, there’s nothing to stop them entering cd ../../../../../../.. and getting access to all sorts of things they shouldn’t.
}
