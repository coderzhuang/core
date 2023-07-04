package core

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"net/http"
)

type HttpService struct {
	e *gin.Engine
	h *http.Server
}

type Middle func(e *gin.Engine)

type Option struct {
	Mode           string
	TrustedProxies []string
	Addr           string
}

type OptionGroup struct {
	dig.In

	Option *Option
	Middle []Middle `group:"middle"`
}

func NewHttp(opts OptionGroup) Server {
	gin.SetMode(opts.Option.Mode)
	e := gin.New()
	_ = e.SetTrustedProxies(opts.Option.TrustedProxies)

	for _, opt := range opts.Middle {
		opt(e)
	}
	e.Use(gin.Recovery())

	server := &HttpService{e: e}
	server.h = &http.Server{
		Addr:    opts.Option.Addr,
		Handler: e,
	}
	return server
}

func (s *HttpService) Run() {
	go func() {
		_ = s.h.ListenAndServe()
	}()
}

func (s *HttpService) Shutdown() {
	_ = s.h.Shutdown(context.Background())
}
