package main

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/atomi-ai/atomi/models"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	host     = "atomi-mysql.mysql.database.azure.com"
	database = "atomi_dev"
	user     = "devuser"
	password = "p6FGFvLcQ2sm"
)

func main() {
	// Initialize connection string.
	rootCertPool := x509.NewCertPool()
	pem, err := ioutil.ReadFile(".config/DigiCertGlobalRootCA.crt.pem")
	if err != nil {
		log.Fatal(err)
	}
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		log.Fatal("Failed to append PEM.")
	}
	mysqlDriver.RegisterTLSConfig("custom", &tls.Config{RootCAs: rootCertPool})
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true&allowNativePasswords=true&tls=custom", user, password, host, database)

	// Open database with standard MySQL driver
	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	// Pass sql.DB instance to GORM
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Create
	gormDB.Create(&models.Product{Name: "D42", Price: 100})

	// Read
	var product models.Product
	gormDB.First(&product, 1)                 // find product with integer primary key
	gormDB.First(&product, "code = ?", "D42") // find product with code D42

	// Update - update product's price to 200
	gormDB.Model(&product).Update("Price", 200)

	// Delete - delete product
	gormDB.Delete(&product)
}
