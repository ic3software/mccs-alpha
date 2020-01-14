package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ic3network/mccs-alpha/global"
	"github.com/ic3network/mccs-alpha/internal/app/http"
	"github.com/ic3network/mccs-alpha/internal/app/service/balancecheck"
	"github.com/ic3network/mccs-alpha/internal/app/service/dailyemail"
	"github.com/ic3network/mccs-alpha/internal/migration"
	"github.com/ic3network/mccs-alpha/internal/pkg/l"
	"github.com/ic3network/mccs-alpha/internal/pkg/version"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
)

func init() {
	global.Init()
}

func main() {
	// Flushes log buffer, if any.
	defer l.Logger.Sync()

	if *global.ShowVersionInfo {
		versionInfo := version.Get()
		marshalled, err := json.MarshalIndent(&versionInfo, "", "  ")
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(marshalled))
		return
	}

	go ServeBackGround()
	go RunMigration()

	http.AppServer.Run(viper.GetString("port"))
}

// ServeBackGround performs the background activities.
func ServeBackGround() {
	c := cron.New()
	viper.SetDefault("daily_email_schedule", "0 0 7 * * *")
	c.AddFunc(viper.GetString("daily_email_schedule"), func() {
		l.Logger.Info("[ServeBackGround] Running daily email schedule. \n")
		dailyemail.Run()
	})
	viper.SetDefault("balance_check_schedule", "0 0 * * * *")
	c.AddFunc(viper.GetString("balance_check_schedule"), func() {
		l.Logger.Info("[ServeBackGround] Running balance check schedule. \n")
		balancecheck.Run()
	})
	c.Start()
}

func RunMigration() {
	// Runs at 2019-08-20
	migration.SetUserActionCategory()
}
