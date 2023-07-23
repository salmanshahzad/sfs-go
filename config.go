package main

import (
	"errors"
	"os"
	"regexp"
)

const ipv4_regex = `^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4}$`

type config struct {
	dir  string
	host string
	port int
}

func (c *config) validate() error {
	if stat, err := os.Stat(c.dir); err != nil || !stat.IsDir() {
		return errors.New("Invalid directory")
	}

	if matched, err := regexp.Match(ipv4_regex, []byte(c.host)); err != nil || !matched {
		return errors.New("Invalid host")
	}

	if c.port < 1024 || c.port > 65535 {
		return errors.New("Invalid port")
	}

	return nil
}
