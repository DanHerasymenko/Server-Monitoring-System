package service

import (
	"Server-Monitoring-System/internal/config"
	"Server-Monitoring-System/internal/constants"
	"Server-Monitoring-System/internal/logger"
	pb "Server-Monitoring-System/proto"
	"context"
	"fmt"
	"github.com/kardianos/service"
	"log"
	"os"
	"time"
)

type Service struct {
	cfg     *config.Config
	context context.Context
	cancel  context.CancelFunc
	client  pb.MonitoringService_StreamMetricsClient
}

func NewService(cfg *config.Config, ctx context.Context, cancel context.CancelFunc, client pb.MonitoringService_StreamMetricsClient) *Service {
	return &Service{
		cfg:     cfg,
		context: ctx,
		cancel:  cancel,
		client:  client,
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
			collectedMetrics, err := s.CollectMetrics()
			if err != nil {
				logger.Error(s.context, fmt.Errorf("error collecting metrics: %v", err))
				continue
			}

			err = s.SendMetrics(collectedMetrics, s.client)
			if err != nil {
				logger.Error(s.context, err)
				continue
			}

			time.Sleep(time.Duration(s.cfg.CollectMetricsInterval) * time.Second)
		}
	}
}

// RunAgentService – stream initialization
func (s *Service) RunAgentService() {

	// add server_services and agent_services IP to context
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
