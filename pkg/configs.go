package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
)

type DB_config struct {
	DSN       string `json:"-"`
	DSNConfig struct {
		Addr   string `json:"addr"`
		Port   string `json:"port"`
		DBNAME string `json:"DB_NAME"`
	} `json:"DSN_config"`
	DBUsers struct {
		Admin      string `json:"admin"`
		AdminPass  string `json:"-"`
		Search     string `json:"search"`
		SearchPass string `json:"-"`
	} `json:"DB_users"`
}

func LoadConfig(path string) (*DB_config, error) {
	//os.Create("test")
	wd, _ := os.Getwd()
	path = filepath.Join(wd, "DB_config.json")
	fmt.Printf("%v\n", path)
	DB := &DB_config{}
	if filepath.Ext(path) != ".json" {
		return nil, errors.New("file is not json or doesnt exist")
	}
	file, err := os.ReadFile(path)
	if err != nil {
		log.Println("error reading file ", err)
		return nil, err
	}
	if environment := os.Getenv("ENVIRONMENT"); environment == "" {
		if err := godotenv.Load(".env"); err != nil {
			return nil, err
		}
	}
	if err := json.Unmarshal(file, DB); err != nil {
		return nil, errors.New("error unmarshalling file " + err.Error())
	}
	DB.DBUsers.AdminPass = os.Getenv("ADMIN_PASS")
	DB.DBUsers.SearchPass = os.Getenv("SEARCH_PASS")

	DB.DSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&tls=skip-verify&charset=utf8mb4", DB.DBUsers.Admin, DB.DBUsers.AdminPass, DB.DSNConfig.Addr, DB.DSNConfig.Port, DB.DSNConfig.DBNAME)
	return DB, nil
}
