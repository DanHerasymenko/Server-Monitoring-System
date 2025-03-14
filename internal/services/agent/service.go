package agent

import (
	"Server-Monitoring-System/internal/config"
	"Server-Monitoring-System/internal/constants"
	"Server-Monitoring-System/internal/logger"
	"context"
	"github.com/kardianos/service"
	"log"
	"os"
	"time"
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

	logger.Info(s.context, "✅ Agent is running...")

	for {
		select {
		case <-s.context.Done():
			logger.Info(s.context, "Agent is stopped")
			return
		default:
			logger.Info(s.context, "Collecting metrics...")
			s.CollectMetrics()
		}
		time.Sleep(time.Duration(s.cfg.CollectMetricsInterval) * time.Second)
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
		//Dependencies: []string{
		//	constants.DependencieNetwork,
		//	constants.DependencieAfter},
	}

	svc, err := service.New(s, svcConfig)
	if err != nil {
		log.Fatalf("Error creating a service: %v", err)
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
			return
		case "uninstall":
			err = svc.Stop()
			err = svc.Uninstall()
			return
		case "start":
			serviceLogger.Info("1111111111111111111")
			logger.Info(s.context, "111111111111111111")

			err = svc.Start()
		case "stop":
			err = svc.Stop()
			return
		default:
			serviceLogger.Info("Unknown command")
			return
		}

		if err != nil {
			serviceLogger.Errorf("Error: %v", err)
		}

	}

	err = svc.Run()
	if err != nil {
		serviceLogger.Info("error")
	}

}
