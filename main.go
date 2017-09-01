package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/rancher/routes-updater/providers"
	"github.com/rancher/routes-updater/providers/hostgw"
	"github.com/urfave/cli"
)

var VERSION = "v0.0.0-dev"

func main() {
	app := cli.NewApp()
	app.Name = "routes-updater"
	app.Version = VERSION
	app.Usage = "Update L3 routes for per-host-subnet networking"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "listen",
			Value: ":8999",
		},
		cli.BoolFlag{
			Name:   "debug, d",
			EnvVar: "RANCHER_DEBUG",
		},
		cli.StringFlag{
			Name:   "metadata-address",
			Value:  providers.DefaultMetadataAddress,
			EnvVar: "RANCHER_METADATA_ADDRESS",
		},
		cli.StringFlag{
			Name: "provider, p",
		},
	}
	app.Action = func(ctx *cli.Context) {
		if err := appMain(ctx); err != nil {
			logrus.Fatal(err)
		}
	}

	app.Run(os.Args)
}

func appMain(ctx *cli.Context) error {
	if ctx.Bool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	}

	done := make(chan error)

	var p providers.Provider
	switch ctx.String("provider") {
	case hostgw.ProviderName:
		p, err := hostgw.NewInst(ctx.String("metadata-address"))
		if err != nil {
			return err
		}
		p.Start()
	default:
		return errors.New("No provider specified")
	}

	listenPort := ctx.String("listen")
	logrus.Debugf("About to start server and listen on port: %v", listenPort)
	go func() {
		s := &APIServer{P: p}
		done <- s.ListenAndServe(listenPort)
	}()

	return <-done
}
