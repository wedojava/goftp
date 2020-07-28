package ftp

import (
	"fmt"
	"net"
)

type dataPort struct {
	h1, h2, h3, h4 int // host
	p1, p2         int // port
}

// dataPortFromHostPort parses the six-byte IP address format into a struct of its parts that we store on the ftp.Conn instance.
func dataPortFromHostPort(hostPort string) (*dataPort, error) {
	var dp dataPort
	_, err := fmt.Sscanf(hostPort, "%d,%d,%d,%d,%d,%d",
		&dp.h1, &dp.h2, &dp.h3, &dp.h4, &dp.p1, &dp.p2)
	if err != nil {
		return nil, err
	}
	return &dp, nil
}

// toAddress converts that struct to a traditionally formatted IP-address-plus-port that the server can connect to with net.Dial.
func (d *dataPort) toAddress() string {
	if d == nil {
		return ""
	}
	// convert hex port bytes to decimal port
	port := d.p1<<8 + d.p2
	// How do you create a port number up to five digits long from two separate bytes? p1<<8 + p2
	// Start by left-shifting p1 by eight places.
	// If p1 = 00011011, p1<<8 = 0001101100000000.
	// When you add p2, it fills the eight bits left empty by the shift.
	// p2 = 11111111; p1 + p2 = 0001101111111111 = 7167.
	// 1111111111111111 is 65535, is 16 bits
	return fmt.Sprintf("%d.%d.%d.%d:%d", d.h1, d.h2, d.h3, d.h4, port)
}

// dataConnect is a simple wrapper for net.Dial, it establish connection, returns a struct satisfying the net.Conn interface(net.TCPConn), not the custom ftp.Conn.
func (c *Conn) dataConnect() (net.Conn, error) {
	conn, err := net.Dial("tcp4", c.dataPort.toAddress())
	if err != nil {
		return nil, err
	}
	return conn, nil
}
