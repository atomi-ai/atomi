package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v74"
	"gorm.io/gorm"
)

func InitStripe(key string) {
	stripe.Key = key
}

func logAllSettings() {
	fmt.Println("All Viper settings:")
	for key, value := range viper.AllSettings() {
		fmt.Printf("%s: %v\n", key, value)
	}
}

func LoadConfig() {
	configFile := os.Getenv("CONFIG_FILE")
	log.Printf("Load config from file: %v", configFile)
	viper.SetConfigFile(configFile)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
	// The code below is debug only
	if viper.GetBool("logAllSettings") {
		logAllSettings()
	}
}

func LoadTaxRates(db *gorm.DB) {
	taxRatesFileDir := viper.GetString("taxRatesFileDir")

	files, err := filepath.Glob(taxRatesFileDir + "/*.csv")
	if err != nil {
		log.Printf("Error reading tax rates file: %s", err)
	}

	if len(files) > 0 {
		// Truncate table if there are files to import
		db.Exec("TRUNCATE TABLE tax_rates")
		if db.Error != nil {
			log.Printf("Error truncating table: %s", db.Error)
		}
	}

	err = filepath.Walk(taxRatesFileDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".csv" {
			mysql.RegisterLocalFile(path)

			fmt.Println("Processing file:", path)
			stmt := fmt.Sprintf(`
				LOAD DATA LOCAL INFILE '%s' INTO TABLE tax_rates 
				FIELDS TERMINATED BY ',' 
				LINES TERMINATED BY '\n' 
				IGNORE 1 LINES 
				(tax_state, zip_code, estimated_combined_rate) 
				SET created_at = NOW(), updated_at = NOW(), csv='%s'
			`, path, path)
			db.Exec(stmt)
			if db.Error != nil {
				log.Printf("Error load tax rates file: %s", db.Error)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
