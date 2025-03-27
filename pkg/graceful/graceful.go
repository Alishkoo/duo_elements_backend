package graceful

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type Connections interface {
	Close()
}

type Logger interface {
	Printf(format string, args ...interface{})
}

type Graceful struct {
	services []Service
}

func New(services ...Service) *Graceful {
	return &Graceful{
		services: services,
	}
}

func (g *Graceful) Shutdown(errs chan error, log Logger, conns Connections) {
	exitCode := 1

	go func() {
		c := make(chan os.Signal, 1)

		signal.Notify(c, syscall.SIGTERM)
		signal.Notify(c, syscall.SIGKILL)
		signal.Notify(c, syscall.SIGQUIT)
		signal.Notify(c, syscall.SIGINT)

		errs <- fmt.Errorf("%v", <-c)
	}()

	err := <-errs
	log.Printf("terminated: %v", err)

	for _, srv := range g.services {
		log.Printf("stopping service: %v", srv.Name())
		srv.Stop()
		log.Printf("stopped service: %v", srv.Name())
	}

	if conns != nil {
		log.Printf("closing connections")
		conns.Close()
	}

	log.Printf("gracefully shutting down")
	log.Printf("terminated: %v", exitCode)

	os.Exit(exitCode)
}
