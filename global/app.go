package global

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func init() {
	App.RootDir = "."
	if !viper.InConfig("url") {
		App.RootDir = inferRootDir()
	}
}

func inferRootDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	var infer func(d string) string
	infer = func(d string) string {
		if Exist(d + "/configs") {
			return d
		}

		return infer(filepath.Dir(d))
	}

	return infer(cwd)
}

var App = &app{}

type app struct {
	RootDir string
}

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
