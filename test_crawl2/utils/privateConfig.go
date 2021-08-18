package utils

type PrivateConfig struct {
	NaverAPI struct {
		ClientCode string `json: "clientCode"`
		SecretCode string `json: "secretCode"`
	}
}
