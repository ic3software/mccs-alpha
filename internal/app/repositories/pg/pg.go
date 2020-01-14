package pg

import (
	"fmt"
	"log"
	"time"

	"github.com/ic3network/mccs-alpha/global"
	"github.com/ic3network/mccs-alpha/internal/app/types"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
)

var db *gorm.DB

func init() {
	global.Init()
	// TODO: set up test docker environment.
	if viper.GetString("env") == "test" {
		return
	}
	db = New()
}

// New returns an initialized DB instance.
func New() *gorm.DB {
	db, err := gorm.Open("postgres", connectionInfo())
	if err != nil {
		panic(err)
	}

	for {
		err := db.DB().Ping()
		if err != nil {
			log.Printf("PostgreSQL connection error: %+v \n", err)
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}

	autoMigrate(db)

	return db
}

func connectionInfo() string {
	password := viper.GetString("psql.password")
	host := viper.GetString("psql.host")
	port := viper.GetString("psql.port")
	user := viper.GetString("psql.user")
	dbName := viper.GetString("psql.db")

	if password == "" {
		return fmt.Sprintf("host=%s port=%s user=%s dbname=%s "+
			"sslmode=disable", host, port, user, dbName)
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s "+
		"dbname=%s sslmode=disable", host, port, user,
		password, dbName)
}

// AutoMigrate will attempt to automatically migrate all tables
func autoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&types.Account{},
		&types.BalanceLimit{},
		&types.Journal{},
		&types.Posting{},
	).Error
	if err != nil {
		panic(err)
	}
}
