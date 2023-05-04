package models

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func autoMigrate(db *gorm.DB, entities ...interface{}) {
	for _, entity := range entities {
		if err := db.AutoMigrate(entity); err != nil {
			panic("AutoMigrate error:" + err.Error())
		}
	}
}

func InitDB() *gorm.DB {
	enableCustomTls := viper.GetBool("enableCustomTls")
	if enableCustomTls {
		// Initialize connection string.
		rootCertPool := x509.NewCertPool()
		pem, err := ioutil.ReadFile(viper.GetString("pemPath"))
		if err != nil {
			log.Fatal(err)
		}
		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			log.Fatal("Failed to append PEM.")
		}
		mysqlDriver.RegisterTLSConfig("custom", &tls.Config{RootCAs: rootCertPool})
	}
	dsn := viper.GetString("dsn")

	// Open database with standard MySQL driver
	sqlDb, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDb.Close()

	// Pass sql.DB instance to GORM
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDb,
	}), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)}) // Can set LogLevel here
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	autoMigrate(db, &Config{}, &User{}, &Product{}, &Store{}, &ProductStore{},
		&UserStore{}, &UserAddress{}, &Order{}, &OrderItem{})

	return db
}
