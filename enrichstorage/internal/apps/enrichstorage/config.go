package enrichstorage

import (
	"enrichstorage/pkg/config"
)

type Config struct {
	Storage  config.PostgresConnectConfig `json:"Storage"`
	GRPCPort int                          `json:"GRPCPort"`
	HTTPPort int                          `json:"HTTPPort"`
}
