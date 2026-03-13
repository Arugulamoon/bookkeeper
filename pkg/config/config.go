package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// https://www.golinuxcloud.com/golang-parse-yaml-file/
type Config struct {
	Database Database `yaml:"database"`
	Google   Google   `yaml:"google"`
	Server   Server   `yaml:"server"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func (db Database) DSN() string {
	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		db.Host, db.Port, db.Name, db.User, db.Password)
}

type Google struct {
	Auth      GoogleAuth `yaml:"auth"`
	Calendars []Calendar `yaml:"calendars"`
}

func (goog Google) GetCalendars() map[string]string {
	calendars := make(map[string]string)
	for _, cal := range goog.Calendars {
		calendars[cal.Name] = cal.Id
	}
	return calendars
}

type GoogleAuth struct {
	Dir string `yaml:"dir"`
}

type Calendar struct {
	Name string `yaml:"name"`
	Id   string `yaml:"id"`
}

type Server struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (srv Server) Addr() string {
	return fmt.Sprintf("%s:%d", srv.Host, srv.Port)
}

func GetConfig(filename string) (*Config, error) {
	if strings.HasPrefix(filename, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		filename = strings.Replace(filename, "~", homeDir, 1)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
