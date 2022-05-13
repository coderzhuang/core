package core

import (
	"github.com/coderzhuang/core/application"
	"github.com/coderzhuang/core/provider/http_service"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"testing"
)

func TestNew(t *testing.T) {
	var err error
	container := dig.New()
	_ = container.Provide(application.New)
	_ = container.Provide(http_service.New, dig.Group("server"))
	_ = container.Provide(func() *http_service.Option {
		return &http_service.Option{
			Mode:           "debug",
			TrustedProxies: []string{"0.0.0.0"},
			Addr:           "0.0.0.0:8080",
		}
	})
	_ = container.Provide(func() http_service.Middle {
		return func(e *gin.Engine) {
			e.GET("/", func(c *gin.Context) {
				c.String(200, "hello")
			})
		}
	}, dig.Group("middle"))

	err = container.Invoke(func(app *application.Application) {
		app.Start()
	})
	if err != nil {
		panic(err.Error())
	}
}
