package db

import (
	"fmt"
	"time"

	"github.com/naoina/kocha"

	"gopkg.in/mgo.v2"
)

var DatabaseMap = kocha.DatabaseMap{
	"mongodb": {
		Driver: kocha.SettingEnv("LGTM_DB_DRIVER", "mongodb"),
		DSN:    kocha.SettingEnv("LGTM_DB_DSN", "lgtmmaker"),
	},
}

var dbMap = make(map[string]*mgo.Database)

func Get(name string) *mgo.Database {
	return dbMap[name]
}

func init() {
	for name, dbconf := range DatabaseMap {
		var d *mgo.DialInfo
		switch dbconf.Driver {
		case "mongodb":
			d = &mgo.DialInfo{
				Addrs: []string{
					"192.168.59.103:27017",
				},
				Timeout:  10 * time.Second,
				Database: dbconf.DSN,
			}
		default:
			panic(fmt.Errorf("unsupported DB type: %v", dbconf.Driver))
		}
		session, err := mgo.DialWithInfo(d)
		if err != nil {
			panic(err)
		}
		dbMap[name] = session.DB("")
	}
}
