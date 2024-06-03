package app

import (
	"fmt"
	"strings"
)

type Config struct {
	Storage  PostgresConnectConfig `json:"Storage"`
	GRPCPort int                   `json:"GRPCPort"`
	DataBus  struct {
		Brokers       []string `json:"Brokers"`
		Topic         string   `json:"Topic"`
		BatchSize     int      `json:"BatchSize"`
		BatchInterval int      `json:"BatchInterval"`
	} `json:"DataBus"`
}

type PostgresConnectConfig struct {
	User     string `json:"User"`
	Host     string `json:"Host"`
	Port     uint16 `json:"Port"`
	Database string `json:"Database"`
	Password string `json:"Password"`
}

func (pc PostgresConnectConfig) GetPostgresDSN() string {
	host := fmt.Sprintf("host=%s", pc.Host)
	user := fmt.Sprintf("user=%s", pc.User)
	password := fmt.Sprintf("password=%s", pc.Password)
	port := fmt.Sprintf("port=%d", pc.Port)
	db := fmt.Sprintf("dbname=%s", pc.Database)
	sslmode := fmt.Sprintf("sslmode=%s", "disable")
	return strings.Join([]string{host, user, password, port, db, sslmode}, " ")
}
