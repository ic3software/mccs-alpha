package global

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/ic3network/mccs-alpha/internal/pkg/l"
	"github.com/spf13/viper"
)

var (
	once            = new(sync.Once)
	configName      = flag.String("config", "development", "config file name, default is development")
	ShowVersionInfo = flag.Bool("v", false, "show version info or not")
)

func Init() {
	once.Do(func() {
		if !flag.Parsed() {
			flag.Parse()
		}
		if err := initConfig(); err != nil {
			panic(fmt.Errorf("initconfig failed: %s \n", err))
		}
		watchConfig()

		l.Init(viper.GetString("env"))
	})
}

func initConfig() error {
	viper.SetConfigName(*configName)
	viper.AddConfigPath("configs")
	viper.AddConfigPath(App.RootDir + "/configs")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

func watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed: %s", e.Name)
	})
}
