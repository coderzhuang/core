package core

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"google.golang.org/grpc"
	"testing"
)

func TestNew(t *testing.T) {
	var err error
	container := dig.New()
	_ = container.Provide(New)
	// http
	_ = container.Provide(New, dig.Group("server"))
	_ = container.Provide(func() *Option {
		return &Option{
			Mode:           "debug",
			TrustedProxies: []string{"0.0.0.0"},
			Addr:           "0.0.0.0:8080",
		}
	})
	_ = container.Provide(func() Middle {
		return func(e *gin.Engine) {
			e.GET("/", func(c *gin.Context) {
				c.String(200, "hello")
			})
		}
	}, dig.Group("middle"))

	// grpc
	_ = container.Provide(New, dig.Group("server"))
	_ = container.Provide(func() *Option {
		return &Option{
			Addr: "0.0.0.0:8080",
		}
	})
	_ = container.Provide(func() RpcServer {
		return func(e grpc.ServiceRegistrar) {
			// Register
		}
	}, dig.Group("grpc_server"))

	err = container.Invoke(func(app *Application) {
		app.Start()
	})
	if err != nil {
		panic(err.Error())
	}
}
