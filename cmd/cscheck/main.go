package main

import (
	"cscheck/internal/checker"
	"cscheck/internal/config"

	"github.com/sirupsen/logrus"
)

func main() {
	_config, err := config.New()
	if err != nil {
		logrus.Fatalln(err)
	}

	_checker := checker.New(
		*_config.CsIP,
		*_config.CsPort,
		*_config.CsUser,
		*_config.CsPassword,
		*_config.CnaScript,
		*_config.Dir,
		*_config.Bind,
	)

	if err := _checker.Start(); err != nil {
		logrus.Fatalf("HTTP server error: %s", err)
	}
	// t := checker.New("127.0.0.1", 50050, "admin", "admin")
	// checker.Start(address, t)
}
