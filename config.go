package core

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
	"path/filepath"
)

var Conf = Config{}

type Common struct {
	Debug bool `yaml:"debug" json:"debug"` // debug
}

type HttpServer struct {
	Switch         bool     `yaml:"switch" json:"switch"`                   // 开关
	Name           string   `yaml:"name" json:"name"`                       // 服务名称
	Addr           string   `yaml:"addr" json:"addr"`                       // 服务地址
	Mode           string   `yaml:"mode" json:"mode"`                       // gin Mode
	TrustedProxies []string `yaml:"trusted_proxies" json:"trusted_proxies"` // 信任的代理
}

type GrpcServer struct {
	Switch bool   `yaml:"switch" json:"switch"` // 开关
	Name   string `yaml:"name" json:"name"`     // 服务名称
	Addr   string `yaml:"addr" json:"addr"`     // 服务地址
}

type CronServer struct {
	Switch bool `yaml:"switch" json:"switch"` // 开关
}

type Config struct {
	Common     *Common     `yaml:"common" json:"common"`           //
	HttpServer *HttpServer `yaml:"http_server" json:"http_server"` //
	GrpcServer *GrpcServer `yaml:"grpc_server" json:"grpc_server"` //
	CronServer *CronServer `yaml:"cron_server" json:"cron_server"` //
}

func InitConf(cfg interface{}) {
	workPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	configFile := filepath.Join(workPath, "config", "config.yaml")
	_, err = os.Stat(configFile)
	if !(err == nil || os.IsExist(err)) {
		panic("config file does not exists")
	}
	b, _ := ioutil.ReadFile(configFile)
	_ = yaml.Unmarshal(b, &Conf)
	_ = yaml.Unmarshal(b, cfg)
}
