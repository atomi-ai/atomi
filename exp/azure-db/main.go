package main

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/go-sql-driver/mysql"

	_ "github.com/go-sql-driver/mysql"
)

const (
	host     = "atomi-mysql.mysql.database.azure.com"
	database = "atomi_dev"
	user     = "devuser"
	password = "p6FGFvLcQ2sm"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var err error

	// Initialize connection string.
	rootCertPool := x509.NewCertPool()
	pem, _ := ioutil.ReadFile(".config/DigiCertGlobalRootCA.crt.pem")
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		log.Fatal("Failed to append PEM.")
	}
	mysql.RegisterTLSConfig("custom", &tls.Config{RootCAs: rootCertPool})
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true&tls=custom", user, password, host, database)
	db, _ := sql.Open("mysql", connectionString)

	checkError(err)
	defer db.Close()

	err = db.Ping()
	checkError(err)
	fmt.Println("Successfully created connection to database.")

	// Drop previous table of same name if one exists.
	_, err = db.Exec("DROP TABLE IF EXISTS inventory;")
	checkError(err)
	fmt.Println("Finished dropping table (if existed).")

	// Create table.
	_, err = db.Exec("CREATE TABLE inventory (id serial PRIMARY KEY, name VARCHAR(50), quantity INTEGER);")
	checkError(err)
	fmt.Println("Finished creating table.")

	// Insert some data into table.
	sqlStatement, err := db.Prepare("INSERT INTO inventory (name, quantity) VALUES (?, ?);")
	res, err := sqlStatement.Exec("banana", 150)
	checkError(err)
	rowCount, err := res.RowsAffected()
	fmt.Printf("Inserted %d row(s) of data.\n", rowCount)

	res, err = sqlStatement.Exec("orange", 154)
	checkError(err)
	rowCount, err = res.RowsAffected()
	fmt.Printf("Inserted %d row(s) of data.\n", rowCount)

	res, err = sqlStatement.Exec("apple", 100)
	checkError(err)
	rowCount, err = res.RowsAffected()
	fmt.Printf("Inserted %d row(s) of data.\n", rowCount)
	fmt.Println("Done.")
}
