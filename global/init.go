package global

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/ic3network/mccs-alpha/internal/pkg/l"
	"github.com/spf13/viper"
)

var (
	once       = new(sync.Once)
	configName = flag.String(
		"config",
		"development",
		"config file name, default is development",
	)
	ShowVersionInfo = flag.Bool("v", false, "show version info or not")
)

func Init() {
	once.Do(func() {
		isTest := false
		for _, arg := range os.Args {
			if strings.Contains(arg, "test") {
				isTest = true
				break
			}
		}

		if !isTest && !flag.Parsed() {
			flag.Parse()
		}

		if err := initConfig(); err != nil {
			panic(fmt.Errorf("initconfig failed: %s", err))
		}
		watchConfig()

		l.Init(viper.GetString("env"))
	})
}

func initConfig() error {
	if err := setConfigNameAndType(); err != nil {
		return fmt.Errorf("setting config name and type failed: %w", err)
	}

	addConfigPaths()

	if err := setupEnvironmentVariables(); err != nil {
		return fmt.Errorf("setting up environment variables failed: %w", err)
	}

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("reading config failed: %w", err)
	}

	return nil
}

func setConfigNameAndType() error {
	viper.SetConfigName(*configName)
	viper.SetConfigType("yaml")
	return nil
}

func addConfigPaths() {
	viper.AddConfigPath("configs")
	viper.AddConfigPath(App.RootDir + "/configs")
}

func setupEnvironmentVariables() error {
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	return nil
}

func watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed: %s", e.Name)
	})
}
