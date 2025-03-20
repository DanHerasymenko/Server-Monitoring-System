package service

import (
	"Server-Monitoring-System/internal/config"
	"Server-Monitoring-System/internal/constants"
	"Server-Monitoring-System/internal/logger"
	"context"
	"fmt"
	"github.com/kardianos/service"
	"log"
	"os"
)

type Service struct {
	cfg     *config.Config
	context context.Context
	cancel  context.CancelFunc
}

func NewService(cfg *config.Config, ctx context.Context, cancel context.CancelFunc) *Service {
	ctx, cancelFunc := context.WithCancel(ctx)
	return &Service{
		cfg:     cfg,
		context: ctx,
		cancel:  cancelFunc,
	}
}

func (s *Service) Stop(svc service.Service) error {

	s.cancel()
	return nil
}

func (s *Service) Start(svc service.Service) error {

	go func() {
		s.run()
	}()

	return nil
}

func (s *Service) run() {

	logger.Info(s.context, "Agent is running...")

	for {
		select {
		case <-s.context.Done():
			logger.Info(s.context, "Agent is stopped")
			return
		default:
			logger.Info(s.context, "Collecting metrics...")
			s.CollectMetrics()
		}
	}
}

// RunAgentService – stream initialization
func (s *Service) RunAgentService() {

	// add server_service and agent_service IP to context
	s.context = logger.SetServerIP(s.context, s.cfg.ServerIP)
	s.context = logger.SetAgentIP(s.context, s.cfg.AgentIP)

	// configure stream
	svcConfig := &service.Config{
		Name:        constants.ServiceName,
		DisplayName: constants.ServiceDisplayName,
		Description: constants.ServiceDescription,
		//Dependencies: []string{
		//	constants.DependencieNetwork,
		//	constants.DependencieAfter},
	}

	svc, err := service.New(s, svcConfig)
	if err != nil {
		log.Fatalf("Error creating a stream: %v", err)
	}

	serviceLogger, err := svc.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	// якщо є команда (install, uninstall, start, stop), робимо її та робимо return.
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "install":
			err = svc.Install()
			if err != nil {
				serviceLogger.Errorf("Error: %v", err)
			} else {
				serviceLogger.Info("Service installed successfully")
			}
			return
		case "uninstall":
			err = svc.Stop()
			if err != nil {
				serviceLogger.Errorf("Error: %v", err)
			}
			err = svc.Uninstall()
			if err != nil {
				serviceLogger.Errorf("Error: %v", err)
			} else {
				serviceLogger.Info("Service uninstalled successfully")
			}
			return
		case "start":
			err = svc.Start()
			if err != nil {
				serviceLogger.Errorf("Error: %v", err)
			} else {
				serviceLogger.Info("Service started successfully")
			}
		case "stop":
			err = svc.Stop()
			if err != nil {
				serviceLogger.Errorf("Error: %v", err)
			} else {
				serviceLogger.Info("Service stopped successfully")
			}
			return
		default:
			serviceLogger.Info("Unknown command, use install, uninstall, start or stop")
			return
		}

		if err != nil {
			serviceLogger.Errorf("Error: %v", err)
		}

	}
	path, err := os.Getwd()

	fmt.Println(path)
	err = svc.Run()
	if err != nil {
		serviceLogger.Info("error")
	}

}
