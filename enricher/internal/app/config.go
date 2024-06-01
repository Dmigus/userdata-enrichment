package app

// Config это конфигурация приложения
type Config struct {
	DataBus struct {
		Brokers []string `json:"Brokers"`
		Topic   string   `json:"Topic"`
	} `json:"DataBus"`
	AgifyAddress         string `json:"AgifyAddress"`
	GenderizeAddress     string `json:"GenderizeAddress"`
	NationalityAddress   string `json:"NationalityAddress"`
	EnrichStorageAddress string `json:"EnrichStorageAddress"`
}
