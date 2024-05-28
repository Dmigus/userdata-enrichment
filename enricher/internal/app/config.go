package app

import (
	"fmt"
	"strings"
)

// Config это конфигурация приложения
type Config struct {
	DataBus struct {
		Brokers []string `json:"Brokers"`
		Topic   string   `json:"Topic"`
	} `json:"DataBus"`
	Repository PostgresConnectConfig `json:"Repository"`
}

// PostgresConnectConfig это конфиг для подключения к PostgreSQL
type PostgresConnectConfig struct {
	User     string `json:"User"`
	Host     string `json:"Host"`
	Port     uint16 `json:"Port"`
	Database string `json:"Database"`
	Password string `json:"Password"`
}

// GetPostgresDSN возвращает Data Source Name, согласно конфигурации
func (pc PostgresConnectConfig) GetPostgresDSN() string {
	host := fmt.Sprintf("host=%s", pc.Host)
	user := fmt.Sprintf("user=%s", pc.User)
	password := fmt.Sprintf("password=%s", pc.Password)
	port := fmt.Sprintf("port=%d", pc.Port)
	db := fmt.Sprintf("dbname=%s", pc.Database)
	sslmode := fmt.Sprintf("sslmode=%s", "disable")
	return strings.Join([]string{host, user, password, port, db, sslmode}, " ")
}
