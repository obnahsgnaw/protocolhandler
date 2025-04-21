package config

import (
	"errors"
	"fmt"
	"github.com/obnahsgnaw/application/pkg/logging/logger"
	"github.com/obnahsgnaw/goutils/configutil"
	"github.com/obnahsgnaw/protocolhandler/version"
	"os"
)

const (
	Name       = "demo"    // 名称
	ModelNames = "M300RTK" // 原始协议的模型名称, 多个以逗号分割
	ActionId   = 1399      // ActionId
)

type Config struct {
	configutil.BaseConfig
	Application *Application
	Log         *logger.Config
	Register    *Register
	Server      *Server
	Version     bool   `short:"v" long:"version" description:"show version"`
	IniFile     string `short:"c" long:"conf" description:"Ini file"`
}
type Application struct {
	Id         string `long:"cluster-id" description:"Cluster id, if not set the server will be run as a independent server"`
	Name       string `long:"cluster-name" description:"Cluster name, default with the id"`
	InternalIp string `long:"internal-ip" description:"Server ip address"`
	Debug      bool   `long:"debug" description:"Enable debug"`
}

type Register struct {
	Driver    string   `long:"register-driver" description:"register driver, etcd, local"`
	Endpoints []string `long:"register-point" description:"register endpoint, multi able , e.g etcd:2379"`
	Timeout   int64    `long:"register-timeout" description:"register timeout, second" default:"5"`
	RegTtl    int64    `long:"register-ttl" description:"register ttl, second" default:"5"`
}

type Server struct {
	Port int `long:"server-port" default:"9011" description:"Port to listen on"`
}

func (c *Config) validate() error {
	if c.Application.Id == "" {
		return errors.New("config error: application id is required")
	}
	if c.Application.Name == "" {
		return errors.New("config error: application name is required")
	}
	if c.Application.InternalIp == "" {
		c.Application.InternalIp = "127.0.0.1"
	}
	if c.Server.Port == 0 {
		return errors.New("config error: server port is required")
	}
	return nil
}

func Parse() *Config {
	c := &Config{}
	if err := c.ParseFlag(c); err != nil {
		println(err.Error())
		os.Exit(1)
	}

	if c.Version {
		fmt.Println(version.Info().String())
		os.Exit(0)
	}

	if err := c.ParseFile(c, c.IniFile); err != nil {
		println(err.Error())
		os.Exit(2)
	}

	if err := c.validate(); err != nil {
		println(err.Error())
		os.Exit(3)
	}

	return c
}
