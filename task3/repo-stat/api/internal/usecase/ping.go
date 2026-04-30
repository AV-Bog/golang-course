package usecase

import (
	"context"
	"sync"
	"time"
)

type Pinger interface {
	Ping(ctx context.Context) error
}

type ServiceStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type PingResult struct {
	Status   string          `json:"status"`
	Services []ServiceStatus `json:"services"`
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

func (u *Ping) Execute(ctx context.Context) (PingResult, int) {
	services := []ServiceStatus{
		{Name: "processor", Status: "up"},
		{Name: "subscriber", Status: "up"},
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
		if err := u.processorClient.Ping(ctx); err != nil {
			services[0].Status = "down"
		}
	}()

	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
		if err := u.subscriberClient.Ping(ctx); err != nil {
			services[1].Status = "down"
		}
	}()

	wg.Wait()

	overall := "ok"
	statusCode := 200
	for _, s := range services {
		if s.Status == "down" {
			overall = "degraded"
			statusCode = 503
			break
		}
	}

	return PingResult{
		Status:   overall,
		Services: services,
	}, statusCode
}
