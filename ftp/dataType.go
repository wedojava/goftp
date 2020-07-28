package ftp

type dataType int

const (
	ascii dataType = iota
	binary
)

// setDataType was written, caused by Go variables are always initialized to their zero values, every new ftp.Conn we create begins life with dataType = 0, ascii, and we provide a function for the client to change that if necessary.
func (c *Conn) setDataType(args []string) {
	if len(args) == 0 {
		c.respond(status501)
	}

	// FTP clients vary, but the typical example included with GNU Inetutils responds to the commands ascii and image by sending TYPE A or TYPE I to the server, informing the server of the datatype it expects to receive.
	switch args[0] {
	case "A":
		c.dataType = ascii
	case "I": // image/binary
		c.dataType = binary
	default:
		c.respond(status504)
		return
	}
	c.respond(status200)
}
