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
	//serviceLogger.Info(s.context, "Agent stop...")

	s.cancel()
	return nil
}

func (s *Service) Start(svc service.Service) error {
	logger.Info(s.context, "🚀 Agent starting...")

	go func() {
		time.Sleep(1 * time.Second) // ✅ Give Windows a small delay before running
		s.run()
	}()

	return nil // ✅ Ensure `Start()` exits quickly
}

func (s *Service) run() {

	logger.Info(s.context, "✅ Agent is running...")

	for {
		select {
		case <-s.context.Done():
			logger.Info(s.context, "Agent stopped")
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

	// service commands handling
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "status":
			status, err := svc.Status()
			if err == nil {
				serviceLogger.Infof("Service status: %v", status)
			}
		case "install":
			err = svc.Install()
			if err == nil {
				serviceLogger.Info("Service successfully installed!")
			}
		case "uninstall":
			err = svc.Stop()
			if err == nil {
				serviceLogger.Info("Service successfully stopped!")
			}
			err = svc.Uninstall()
			if err == nil {
				serviceLogger.Info("Service successfully uninstalled!")
			}
		case "start":
			err = svc.Start()
			if err == nil {
				serviceLogger.Info(" Service successfully started!")
			}
			// Start service execution
			err = svc.Run()
			if err != nil {
				log.Fatalf("❌ Error running service: %v", err)
			}

		case "stop":
			err = svc.Stop()
			if err == nil {
				serviceLogger.Info("Service successfully stopped!")
			}
		default:
			serviceLogger.Info("Unknown command, available commands: install, uninstall, start, stop")
		}

		if err != nil {
			serviceLogger.Infof("Error: %v", err)
		}

		return
	}

}
