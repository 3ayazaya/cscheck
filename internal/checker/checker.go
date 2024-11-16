package checker

import (
	"cscheck/pkg/shell"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wojas/go-healthz"
)

const (
	StatusDown = "down"
)

type Teamserver struct {
	address  string
	port     uint
	user     string
	password string
	script   string
	dir      string
	bind     string
}

func New(
	address string,
	port uint,
	user string,
	password string,
	script string,
	dir string,
	bind string,
) *Teamserver {
	return &Teamserver{
		address:  address,
		port:     port,
		user:     user,
		password: password,
		script:   script,
		dir:      dir,
		bind:     bind,
	}
}

func (t Teamserver) Check() error {
	out, err := shell.New(
		"/bin/sh",
		"agscript",
		t.address,
		fmt.Sprint(t.port),
		t.user,
		t.password,
		t.script,
	).WithDir(t.dir).Run()
	if err != nil {
		logrus.Warnf("agscript server error: %s", err)
		return healthz.Warn(StatusDown)
	}

	logrus.Infof("agscript server: %s", out)
	return nil
}

func (t *Teamserver) Start() error {
	logrus.WithField("address", t.bind).Info("HTTP stats server enabled")
	healthz.Register("teamserver", time.Second*5, t.Check)
	http.Handle("/healthz", healthz.Handler())
	return http.ListenAndServe(t.bind, nil)
}

// func Start(a string, t *Teamserver) {

// 	if a == "" {
// 		logrus.Info("HTTP stats server disabled")
// 		return
// 	}
// 	logrus.WithField("address", a).Info("HTTP stats server enabled")
// 	healthz.Register("teamserver", time.Second*5, t.Check)
// 	http.Handle("/healthz", healthz.Handler())
// 	err := http.ListenAndServe(a, nil)
// 	logrus.Fatalf("HTTP server error: %v", err)
// }
