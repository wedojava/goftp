package ftp

import "net"

// Conn represents a connection to the ftp server
type Conn struct {
	conn         net.Conn
	dataType     dataType
	dataPort     *dataPort
	rootDir      string
	workDir      string
	pasvListener net.Listener
	cmdErr       error // saved command connection write error.
}

// NewConn returns a new ftp connection
func NewConn(conn net.Conn, rootDir string) *Conn {
	return &Conn{
		conn:    conn,
		rootDir: rootDir,
		workDir: "/",
	}
}
