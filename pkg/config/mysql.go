package config

import (
	"log"
	"os"
	"sync"
)

var doOnce sync.Once
var config *MysqlCfg

type MysqlCfg struct {
	Host      string
	User      string
	Password  string
	UriArgs   string
	DefaultDb string
}

func Get() *MysqlCfg {
	doOnce.Do(func() {
		config = &MysqlCfg{}
		config.Host = MustEnv("MYSQL_HOST")
		config.User = MustEnv("MYSQL_USER")
		config.Password = MustEnv("MYSQL_PASSWORD")
		config.UriArgs = MustEnv("MYSQL_URI_ARGS")
		config.DefaultDb = MustEnv("MYSQL_DEFAULT_DB")
	})
	return config
}

func MustEnv(name string) string {
	value, found := os.LookupEnv(name)
	if !found {
		log.Fatalf("environment variable %s is missing", name)
	}
	return value
}
