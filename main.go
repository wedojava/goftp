package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"path/filepath"

	"github.com/wedojava/goftp/ftp"
)

var (
	port    int
	rootDir string
)

func init() {
	flag.IntVar(&port, "port", 8000, "port number")
	flag.StringVar(&rootDir, "rootDir", "public", "root Directory")
	flag.Parse()
}

func main() {
	server := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", server)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(c net.Conn) {
	defer c.Close()
	absPath, err := filepath.Abs(rootDir)
	if err != nil {
		log.Fatal(err)
	}
	ftp.Serve(ftp.NewConn(c, absPath))
}
