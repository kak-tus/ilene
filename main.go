package main

import (
	"git.aqq.me/go/app/appconf"
	"git.aqq.me/go/app/launcher"
	"github.com/iph0/conf/envconf"
	"github.com/iph0/conf/fileconf"
	"github.com/kak-tus/ilene/api"
)

func init() {
	fileLdr := fileconf.NewLoader("etc", "/etc")
	envLdr := envconf.NewLoader()

	appconf.RegisterLoader("file", fileLdr)
	appconf.RegisterLoader("env", envLdr)

	appconf.Require("file:ilene.yml")
	appconf.Require("env:^ILENE_")
}

func main() {
	launcher.Run(func() error {
		go api.GetAPI().Start()
		return nil
	})
}
