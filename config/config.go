package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type DbConfig struct {
	Host	string 	`yaml:"host"`
	Port	string	`yaml:"port"`
	User	string	`yaml:"user"`
	Pwd		string	`yaml:"pwd"`
	Name	string	`yaml:"name"`
}

type ApiConfig struct {
	Port 		string 	`yaml:"port"`
	JwtSecret	string 	`yaml:"jwtSecret"`
}

// Representation of YAMAL config file as struct
type Config struct {
	DB 	DbConfig  	`yaml:"db"`

	API ApiConfig 	`yaml:"api"`
}

// Package Variable to use in entire project
var Cfg Config

// LoadConfig reads a YAML configuration file and unmarshals it into a Config struct
func LoadConfig() error {
	fileName := "conf.yaml"
	data, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("Faile to read config file.", err)
	}
	if err := yaml.Unmarshal(data, &Cfg); err != nil {
		return fmt.Errorf("Faile to unmarshal data.", err)
	}
	return nil
}
