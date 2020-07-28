package ftp

import (
	"bytes"
	"fmt"
	"log"
	"net"
)

func (c *Conn) pasv(args []string) {
	if len(args) > 0 {
		c.respond(status501)
		return
	}
	var firstError error
	storeFirstError := func(err error) {
		if firstError == nil {
			firstError = err
		}
	}
	var err error
	c.pasvListener, err = net.Listen("tcp4", "")
	storeFirstError(err)
	_, port, err := net.SplitHostPort(c.pasvListener.Addr().String())
	storeFirstError(err)
	ip, _, err := net.SplitHostPort(c.conn.LocalAddr().String())
	storeFirstError(err)
	addr, err := dataPortFromHostPort(fmt.Sprintf("%s:%s", ip, port))
	storeFirstError(err)
	c.dataPort = addr
	if firstError != nil {
		c.pasvListener.Close()
		c.pasvListener = nil
		c.log(logPairs{"cmd": "PASV", "err": err})
		c.respond(status451)
		return
	}
	// DJB recommends putting an extra character before the address.
	c.respond("227" + c.dataPort.toAddress())
}

type logPairs map[string]interface{}

func (c *Conn) log(pairs logPairs) {
	b := &bytes.Buffer{}
	fmt.Fprintf(b, "addr=%s", c.conn.RemoteAddr().String())
	for k, v := range pairs {
		fmt.Fprintf(b, " %s=%s", k, v)

	}
	log.Print(b.String())

}
