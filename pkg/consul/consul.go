package consul

import (
	"fmt"
	"time"

	api "github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
)

const (
	TTLDeregisterCriticalServiceAfter = time.Second * 30
)

type ConsulClient struct {
	client        *api.Client
	serviceID     string
	serviceName   string
	servicePort   int
	serviceAddr   string
	checkInterval string
	checkTimeout  string
	deregisterTTL string
	checkerPort   string
}

func NewConsulClient(serviceName string, serviceID string, serviceAddr string, servicePort int, checkerPort string) (*ConsulClient, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}

	return &ConsulClient{
		client:        client,
		serviceID:     serviceID,
		serviceName:   serviceName,
		serviceAddr:   serviceAddr,
		servicePort:   servicePort,
		checkerPort:   checkerPort,
		checkInterval: "10s",                                      // Default health check interval
		checkTimeout:  "1s",                                       // Default health check timeout
		deregisterTTL: TTLDeregisterCriticalServiceAfter.String(), // Default deregister TTL
	}, nil
}

func (c *ConsulClient) Register() error {
	service := &api.AgentServiceRegistration{
		ID:      c.serviceID,
		Name:    c.serviceName,
		Port:    c.servicePort,
		Address: c.serviceAddr,
		Check: &api.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("http://%s:%s/healthz", c.serviceAddr, c.checkerPort),
			Interval:                       c.checkInterval,
			Timeout:                        c.checkTimeout,
			DeregisterCriticalServiceAfter: c.deregisterTTL,
		},
	}

	logrus.Infof("Registering service: %s", c.serviceName)
	return c.client.Agent().ServiceRegister(service)
}

func (c *ConsulClient) Deregister() error {
	logrus.Infof("Deregistering service: %s", c.serviceName)
	return c.client.Agent().ServiceDeregister(c.serviceID)
}
