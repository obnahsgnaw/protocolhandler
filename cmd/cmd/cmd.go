package cmd

import (
	"errors"
	"fmt"
	"github.com/obnahsgnaw/application"
	"github.com/obnahsgnaw/application/service/regCenter"
	"github.com/obnahsgnaw/goutils/runtimeutil"
	protocolhandlerapplication "github.com/obnahsgnaw/protocolhandler/application"
	"github.com/obnahsgnaw/protocolhandler/config"
	"github.com/obnahsgnaw/protocolhandler/version"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var RootCmd = &cobra.Command{
	Use:     "protocol-handler",
	Short:   "protocol-handler",
	Long:    `protocol-handler service`,
	Version: version.Version(),
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	RootCmd.Flags().StringP("conf", "c", "", "config file path")
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		PrintError(err)
		os.Exit(1)
	}
}

func run() {
	defer runtimeutil.HandleRecover(func(err, stack string) {
		PrintError(errors.New(config.Name + ":" + err))
		os.Exit(1)
	})

	cnf := config.Parse()

	reg, err := regCenter.NewEtcdRegister(cnf.Register.Endpoints, time.Duration(cnf.Register.Timeout)*time.Second)
	if err != nil {
		PrintError(err)
		os.Exit(1)
	}
	app := application.New("protocol-handler:"+config.Name,
		application.CusCluster(application.NewCluster(cnf.Application.Id, cnf.Application.Name)),
		application.Register(reg, cnf.Register.RegTtl),
		application.Logger(cnf.Log),
		application.Debug(func() bool {
			return cnf.Application.Debug
		}),
	)
	defer app.Release()

	protocolhandlerapplication.New(app, protocolhandlerapplication.IpPort(cnf.Application.InternalIp, cnf.Server.Port))

	app.Run(func(err error) {
		PrintError(err)
		os.Exit(2)
	})

	app.Wait()
}

func PrintError(err error) {
	_, _ = fmt.Fprintln(os.Stderr, err)
}
