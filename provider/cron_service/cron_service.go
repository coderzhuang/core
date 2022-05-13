package cron_service

import (
	"github.com/coderzhuang/core/application"
	"github.com/robfig/cron/v3"
)

type CronService struct {
	c *cron.Cron
}

type CronClosure func(*cron.Cron)

func New(fn CronClosure) application.Service {
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
