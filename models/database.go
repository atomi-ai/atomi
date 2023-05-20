package models

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"log"
	"os"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"

	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func autoMigrateTables(db *gorm.DB, entities ...interface{}) {
	for _, entity := range entities {
		if err := db.AutoMigrate(entity); err != nil {
			panic("AutoMigrate error:" + err.Error())
		}
	}
}

func AutoMigrate(db *gorm.DB) {
	autoMigrateTables(db, &Config{}, &User{}, &Product{}, &Store{}, &ProductStore{}, &UserStore{}, &UserAddress{}, &Order{}, &OrderItem{}, &ManagerStores{}, &DeleteUserRequest{}, &TaxRate{})
}

func InitDB() *gorm.DB {
	enableCustomTLS := viper.GetBool("enableCustomTLS")
	if enableCustomTLS {
		// Initialize connection string.
		rootCertPool := x509.NewCertPool()
		pem, err := os.ReadFile(viper.GetString("pemPath"))
		if err != nil {
			log.Fatal("Failed to read pem file", err)
		}
		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			log.Fatal("Failed to append PEM.")
		}
		if err = mysqlDriver.RegisterTLSConfig("custom", &tls.Config{
			RootCAs:    rootCertPool,
			MinVersion: tls.VersionTLS12,
		}); err != nil {
			log.Fatal("Failed to register TLS config", err)
		}
	}
	dsn := viper.GetString("dsn")

	// Open database with standard MySQL driver
	sqlDb, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// Pass sql.DB instance to GORM
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDb,
	}), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)}) // Can set LogLevel here
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return db
}
