package outboxsender

import "enrichstorage/pkg/config"

type Config struct {
	Storage         config.PostgresConnectConfig `json:"Storage"`
	RequestEventBus struct {
		Brokers       []string `json:"Brokers"`
		Topic         string   `json:"Topic"`
		BatchSize     int      `json:"BatchSize"`
		BatchInterval int      `json:"batchInterval"`
		Username      string
		Password      string
	} `json:"RequestEventBus"`
}
