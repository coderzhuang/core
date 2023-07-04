package core

import (
	"fmt"
	"go.uber.org/dig"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var services []func(*dig.Container)

type Server interface {
	Run()
	Shutdown()
}

type ServerGroup struct {
	dig.In

	Servers []Server `group:"server"`
}

type Application struct {
	Servers []Server
}

func New(sg ServerGroup) *Application {
	return &Application{Servers: sg.Servers}
}

func (a *Application) Start() {
	if len(a.Servers) == 0 {
		log.Println("There is no Servers")
		return
	}

	for _, server := range a.Servers {
		go server.Run()
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	for _, server := range a.Servers {
		server.Shutdown()
	}
	log.Println("Servers Shut down")
}

func RegistryService(fn func(*dig.Container)) {
	services = append(services, fn)
}

func InitProvider() *dig.Container {
	container := dig.New()
	// 以下为系统级别服务
	_ = container.Provide(New)
	if Conf.HttpServer.Switch {
		_ = container.Provide(func() *Option {
			return &Option{
				Mode:           Conf.HttpServer.Mode,
				TrustedProxies: Conf.HttpServer.TrustedProxies,
				Addr:           Conf.HttpServer.Addr,
			}
		})
		_ = container.Provide(NewHttp, dig.Group("server"))
		//_ = container.Provide(router.InitRoute, dig.Group("middle"))
	}
	if Conf.GrpcServer.Switch {
		_ = container.Provide(func() *OptionRpc {
			return &OptionRpc{
				Addr: Conf.GrpcServer.Addr,
			}
		})
		_ = container.Provide(NewRpc, dig.Group("server"))
	}
	if Conf.CronServer.Switch {
		_ = container.Provide(NewCron, dig.Group("server"))
		//_ = container.Provide(cron.InitCron)
	}

	return container
}

func Run() {
	InitConf()
	container := InitProvider()
	// 加载业务相关服务
	for _, fn := range services {
		fn(container)
	}
	err := container.Invoke(func(app *Application) {
		app.Start()
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}
