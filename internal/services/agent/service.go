package agent

import (
	"Server-Monitoring-System/internal/config"
	"Server-Monitoring-System/internal/constants"
	"Server-Monitoring-System/internal/logger"
	"context"
	"fmt"
	"github.com/kardianos/service"
	"log"
)

type Service struct {
	cfg     *config.Config
	context context.Context
	cancel  context.CancelFunc
}

func NewService(cfg *config.Config, ctx context.Context) *Service {
	ctx, cancel := context.WithCancel(ctx)
	return &Service{
		cfg:     cfg,
		context: ctx,
		cancel:  cancel,
	}
}

func (s *Service) Start(svc service.Service) error {
	s.context, s.cancel = context.WithCancel(context.Background())

	go s.run()
	return nil
}

func (s *Service) Stop(svc service.Service) error {
	logger.Info(s.context, "Agent stop...")

	s.cancel()
	return nil
}

func (s *Service) run() {
	logger.Info(s.context, "Agent started in background...")

	for {
		select {
		case <-s.context.Done():
			logger.Info(s.context, "Agent stopped")
			return
		default:
			s.CollectMetrics()
		}
	}
}

// RunAgentService – service initialization
func (s *Service) RunAgentService() {

	// add server and agent IP to context
	s.context = logger.SetServerIP(s.context, s.cfg.ServerIP)
	s.context = logger.SetAgentIP(s.context, s.cfg.AgentIP)

	// configure service
	svcConfig := &service.Config{
		Name:        constants.ServiceName,
		DisplayName: constants.ServiceDisplayName,
		Description: constants.ServiceDescription,
	}

	svc, err := service.New(s, svcConfig)
	if err != nil {
		log.Fatalf("Error creating a service: %v", err)
	}

	// service commands handling
	if len(svcConfig.Arguments) > 1 {
		switch svcConfig.Arguments[1] {
		case "install":
			err = svc.Install()
			if err == nil {
				fmt.Println("Service successfully installed!")
			}
		case "uninstall":
			err = svc.Uninstall()
			if err == nil {
				fmt.Println("Service successfully uninstalled!")
			}
		case "start":
			err = svc.Start()
			if err == nil {
				fmt.Println("Service started!")
			}
		case "stop":
			err = svc.Stop()
			if err == nil {
				fmt.Println("Service stopped!")
			}
		default:
			fmt.Println("Available commands: install, uninstall, start, stop")
		}

		if err != nil {
			fmt.Println("Error:", err)
		}
		return
	}

	// start service
	err = svc.Run()
	if err != nil {
		log.Fatalf("Error starting the service: %v", err)
	}
}
