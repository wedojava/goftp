package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/wedojava/goftp/ftp"
)

var (
	port int
	root string
)

func init() {
	flag.IntVar(&port, "port", 8000, "port number")
	flag.StringVar(&root, "root", "public", "root Directory")
	flag.Parse()
	// Create root dir if it is not exist.
	_, err := os.Stat(root)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(root, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	server := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp4", server)
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
	absPath, err := filepath.Abs(root)
	if err != nil {
		log.Fatal(err)
	}
	ftp.Serve(ftp.NewConn(c, absPath))
}
