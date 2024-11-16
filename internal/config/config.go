package config

import (
	"errors"
	"flag"
)

var (
	ErrEmptyPassword = errors.New("empty password")
)

type Config struct {
	CsIP       *string
	CsPort     *uint
	CsUser     *string
	CsPassword *string
	CnaScript  *string

	Dir *string

	Bind *string
}

func New() (*Config, error) {
	c := &Config{}
	c.CsIP = flag.String("ip", "127.0.0.1", "Cobalt-Strike admin IP")
	c.CsPort = flag.Uint("port", 50050, "Cobalt-Strike admin port")
	c.CsUser = flag.String("user", "user", "Cobalt-Strike admin username")
	c.CsPassword = flag.String("password", "", "Cobalt-Strike admin password")
	c.CnaScript = flag.String("cna", "checkConnection.cna", "Agressor cna script")
	c.Bind = flag.String("bind", "127.0.0.1:8000", "Bind checker web address")
	c.Dir = flag.String("dir", "", "Directory with agscript")
	flag.Parse()

	if *c.CsPassword == "" {
		return nil, ErrEmptyPassword
	}
	return c, nil
}
