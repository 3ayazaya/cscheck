package checker

import (
	"context"
	"cscheck/pkg/shell"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexliesenfeld/health"
	"github.com/sirupsen/logrus"
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

func (t Teamserver) Check(ctx context.Context) error {
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
		return fmt.Errorf(StatusDown)
	}
	if !strings.Contains(out, "Hello World") {
		logrus.Warnf("agscript server error: %s", out)
		return fmt.Errorf(StatusDown)
	}

	logrus.Infof("agscript server: %s", out)
	return nil
}

func (t *Teamserver) Start() error {
	logrus.WithField("address", t.bind).Info("HTTP stats server enabled")
	checker := health.NewChecker(
		health.WithCacheDuration(1*time.Second),
		health.WithTimeout(10*time.Second),
		health.WithPeriodicCheck(10*time.Second, 3*time.Second, health.Check{
			Name:  "teamserver",
			Check: t.Check,
		}),
	)
	http.Handle("/healthz", health.NewHandler(checker))
	return http.ListenAndServe(t.bind, nil)
}
