package config

import (
	"context"
	"github.com/spf13/viper"
	etcdclientv3 "go.etcd.io/etcd/client/v3"
	"strings"
	"time"
)

func ViperInit() {
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	// 配置中心获取配置
	configCenterConfig := viper.GetStringMapString("config_center")
	if configCenterConfig["host"] != "" && configCenterConfig["port"] != "" {
		client, err := etcdclientv3.New(etcdclientv3.Config{
			Endpoints:         []string{configCenterConfig["host"] + ":" + configCenterConfig["port"]},
			DialTimeout:       3 * time.Second,
			DialKeepAliveTime: 3 * time.Second,
		})
		if err != nil {
			panic(err)
		}

		resp, err := client.Get(context.Background(), "/services/"+viper.GetString("service_name")+"/config/app.json")
		if err != nil {
			panic(err)
		}
		for _, ev := range resp.Kvs {
			err = viper.MergeConfig(strings.NewReader(string(ev.Value)))
			if err != nil {
				panic(err)
			}
		}
	}
}
