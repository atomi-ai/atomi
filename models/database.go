package models

import (
	"fmt"
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
	dbUrl := "appuser:p6FGFvLcQ2sm@tcp(127.0.0.1:3306)/atomi_exp?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dbUrl), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Can set LogLevel here
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	autoMigrate(db, &Config{}, &User{}, &Product{}, &Store{}, &ProductStore{})

	return db
}
