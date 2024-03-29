package core

import (
	"github.com/robfig/cron/v3"
)

type CronService struct {
	c *cron.Cron
}

type CronClosure func(*cron.Cron)

func NewCron(fn CronClosure) Server {
	c := cron.New()
	fn(c)
	return &CronService{c: c}
}

func (s *CronService) Run() {
	s.c.Start()
}

func (s *CronService) Shutdown() {
	ctx := s.c.Stop()
	<-ctx.Done()
}
