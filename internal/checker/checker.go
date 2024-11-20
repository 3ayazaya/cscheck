package checker

import (
	"context"
	"cscheck/internal/metrics"
	"cscheck/pkg/shell"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/alexliesenfeld/health"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	logrus.Infof("agscript server: %s", "Teamserver is up")
	return nil
}

func (t Teamserver) GetListeners(
	script string,
	teamserverListenersTotal *metrics.GaugeMetrics,
	teamserverListeners *metrics.GaugeVecMetrics,
) error {
	out, err := shell.New(
		"/bin/sh",
		"agscript",
		t.address,
		fmt.Sprint(t.port),
		"exporter",
		t.password,
		script,
	).WithDir(t.dir).Run()
	if err != nil {
		logrus.Warnf("agscript server error: %s", err)
		return fmt.Errorf(StatusDown)
	}
	logrus.Infof("Current listeners: %s", fmt.Sprint(strings.Count(out, "Listener:")))
	re, _ := regexp.Compile("== Listener: (.*) ==")
	matches := re.FindAllStringSubmatch(out, -1)
	teamserverListeners.Metrics.Reset()
	for _, match := range matches {
		teamserverListeners.Metrics.WithLabelValues(match[1]).Set(1)
	}
	teamserverListenersTotal.Metrics.Set(float64(strings.Count(out, "Listener:")))
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

	teamserverListenersTotal := metrics.NewGaugeMetrics("teamserver_listeners_total", "Cobalt Strike listeners total count.")
	teamserverListeners := metrics.NewGaugeVecMetrics("teamserver_listener", "Cobalt Strike listener.")

	go func() {
		for {
			t.GetListeners(
				"checkListeners.cna",
				teamserverListenersTotal,
				teamserverListeners,
			)
			time.Sleep(10 * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/healthz", health.NewHandler(checker))
	return http.ListenAndServe(t.bind, nil)
}
