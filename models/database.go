package models

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

func autoMigrate(db *gorm.DB, entities ...interface{}) {
	for _, entity := range entities {
		if err := db.AutoMigrate(entity); err != nil {
			panic("AutoMigrate error:" + err.Error())
		}
	}
}

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(viper.GetString("dbUrl")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Can set LogLevel here
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	autoMigrate(db, &Config{}, &User{}, &Product{}, &Store{}, &ProductStore{}, &UserStore{}, &UserAddress{})

	return db
}
