package app

// Config это конфигурация приложения
type Config struct {
	RequestEventBus struct {
		Brokers  []string `json:"Brokers"`
		Topic    string   `json:"Topic"`
		Username string
		Password string
	} `json:"RequestEventBus"`
	AgifyAddress         string `json:"AgifyAddress"`
	GenderizeAddress     string `json:"GenderizeAddress"`
	NationalityAddress   string `json:"NationalityAddress"`
	EnrichStorageAddress string `json:"EnrichStorageAddress"`
}
