package consul

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/consul/api"
	consul "github.com/hashicorp/consul/api"
)

type Registry struct{
	client *consul.Client
}

func NewRegistory(addr, serviceName string)  (*Registry, error) {
	config := consul.DefaultConfig()
	config.Address = addr
	client, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Registry{client}, nil
}

func(r *Registry) Register(ctx context.Context, instanceID, serverName, hostPort string) error {
	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		return errors.New("invalied host:port format. eg: localhost:8081v")
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}
	host := parts[0]

	return r.client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		ID: instanceID,	
		Address: host,
		Port: port,
		Name: serverName,
		Check: &consul.AgentServiceCheck{
			CheckID: instanceID,
			TLSSkipVerify: true,
			TTL: "5s",
			DeregisterCriticalServiceAfter: "10s",
		},
	})
}

func(r *Registry) DeRegister(ctx context.Context, instanceID, serverName string) error {
	log.Printf("Deregistering service %s", instanceID)
	return r.client.Agent().CheckDeregister(instanceID)
}

func(r *Registry) HealthCheck(instanceID, serverName string) error {
	return r.client.Agent().UpdateTTL(instanceID, "online", api.HealthPassing)
}

func(r *Registry) Descover(ctx context.Context, serviceName string) ([]string, error) {
	entries, _, err := r.client.Health().Service(serviceName, "", true, nil )
	if err != nil {
		return nil ,err
	}

	var instance []string
	for _, entry := range entries {
		instance = append(instance, fmt.Sprintf("%s:%d", entry.Service.Address, entry.Service.Port))
	}

	return instance, nil
}
