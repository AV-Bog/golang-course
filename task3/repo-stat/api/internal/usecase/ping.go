package usecase

import (
	"context"
	"sync"
)

type Pinger interface {
	Ping(ctx context.Context) error
}

type Ping struct {
	processorClient  Pinger
	subscriberClient Pinger
}

func NewPing(processorClient, subscriberClient Pinger) *Ping {
	return &Ping{
		processorClient:  processorClient,
		subscriberClient: subscriberClient,
	}
}

type ServiceStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type PingResult struct {
	Status   string          `json:"status"`
	Services []ServiceStatus `json:"services"`
}

func (u *Ping) Execute(ctx context.Context) (PingResult, int) {
	services := []ServiceStatus{
		{Name: "processor", Status: "up"},
		{Name: "subscriber", Status: "up"},
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := u.processorClient.Ping(ctx); err != nil {
			services[0].Status = "down"
		}
	}()

	go func() {
		defer wg.Done()
		if err := u.subscriberClient.Ping(ctx); err != nil {
			services[1].Status = "down"
		}
	}()

	wg.Wait()

	overall := "ok"
	for _, s := range services {
		if s.Status == "down" {
			overall = "degraded"
			break
		}
	}

	statusCode := 200
	if overall == "degraded" {
		statusCode = 503
	}

	return PingResult{
		Status:   overall,
		Services: services,
	}, statusCode
}
