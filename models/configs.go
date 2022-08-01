package models

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type DBConfigs struct{
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	DBName string `yaml:"db_name"`
	User string `yaml:"user"`
	Password string `yaml:"-"`
	SSLMode string	`yaml:"ssl_mode"`
}

func GetDBConfigs() *DBConfigs {
	godotenv.Load()
	
	data, err := ioutil.ReadFile("configs/db.yaml")
	if err != nil{
		log.Fatal(err)
	}

	var configs DBConfigs

	err = yaml.Unmarshal(data, &configs)
	if err != nil{
		log.Fatal(err)
	}

	configs.Password = os.Getenv("DB_PASSWORD")
	return &configs
}

