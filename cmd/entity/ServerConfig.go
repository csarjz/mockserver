package entity

type ServerConfig struct {
	Port    uint16  `json:"port"`
	BaseUrl string  `json:"baseUrl"`
	Routes  []Route `json:"routes"`
}
