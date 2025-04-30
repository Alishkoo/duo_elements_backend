package graceful

type Service interface {
	Name() string
	Stop()
}

func NewService(name string, stop func()) Service {
	return &service{
		name: name,
		stop: stop,
	}
}

type service struct {
	name string
	stop func()
}

func (s *service) Name() string {
	return s.name
}

func (s *service) Stop() {
	s.stop()
}
