package model

import (
	"strings"

	"git.aqq.me/go/app/appconf"
	"git.aqq.me/go/app/applog"
	"git.aqq.me/go/app/event"
	"github.com/go-redis/redis"
	"github.com/iph0/conf"
)

var obj *Type

func init() {
	event.Init.AddHandler(
		func() error {
			cnfMap := appconf.GetConfig()["model"]

			var cnf configType
			err := conf.Decode(cnfMap, &cnf)
			if err != nil {
				return err
			}

			addrs := strings.Split(cnf.Redis.Addrs, ",")

			var parsed []string
			var pass string

			for _, a := range addrs {
				opt, err := redis.ParseURL(a)

				if err == nil {
					// New format
					parsed = append(parsed, opt.Addr)
					pass = opt.Password
				} else {
					// Backward compatible
					// TODO FIX remove
					parsed = append(parsed, a)
				}
			}

			rDB := redis.NewClusterClient(&redis.ClusterOptions{
				Addrs:    parsed,
				Password: pass,
			})

			obj = &Type{
				cnf: cnf,
				log: applog.GetLogger().Sugar(),
				rDB: rDB,
			}

			obj.log.Info("Started model")

			return nil
		},
	)
	event.Stop.AddHandler(
		func() error {
			obj.log.Info("Stop model")
			obj.log.Info("Stopped model")
			return nil
		},
	)
}

// Get return instance
func Get() *Type {
	return obj
}
