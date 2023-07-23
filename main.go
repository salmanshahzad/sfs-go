package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
)

func main() {
	dir := flag.String("d", ".", "Directory")
	host := flag.String("h", "127.0.0.1", "Host")
	port := flag.Int("p", 1024, "Port")
	flag.Parse()

	cfg := config{
		dir:  *dir,
		host: *host,
		port: *port,
	}
	if err := cfg.validate(); err != nil {
		bail(err)
	}

	server, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.host, cfg.port))
	if err != nil {
		bail(err)
	}

	for {
		conn, err := server.Accept()
		if err != nil {
			stderrPrint("Error accepting connection", err)
		} else {
			go handleConnection(conn, cfg.dir)
		}
	}
}

func handleConnection(c net.Conn, dir string) {
	defer func() {
		if err := c.Close(); err != nil {
			stderrPrint("Error closing connection", err)
		}
	}()

	buf, err := readAll(c)
	if err != nil {
		stderrPrint("Error reading from connection", err)
		return
	}

	req, err := parseRequest(buf)
	var res response
	if err != nil {
		res = responseWithStatus(http.StatusBadRequest)
	} else {
		res = handleRequest(&req, dir)
	}

	if err := writeAll(c, res.getBytes()); err != nil {
		stderrPrint("Error writing to connection", err)
	}
}

func readAll(c net.Conn) ([]byte, error) {
	full := make([]byte, 0)
	for {
		buf := make([]byte, 1024)
		n, err := c.Read(buf)
		if err != nil {
			return full, err
		}

		full = append(full, buf[:n]...)
		if n < len(buf) {
			return full, nil
		}
	}
}

func writeAll(c net.Conn, buf []byte) error {
	written := 0
	for written < len(buf) {
		n, err := c.Write(buf[written:])
		if err != nil {
			return err
		}
		written += n
	}
	return nil
}

func bail(err error) {
	stderrPrint(err.Error())
	os.Exit(1)
}

func stderrPrint(a ...any) {
	fmt.Fprintln(os.Stderr, a...)
}
