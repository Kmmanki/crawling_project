package utils

type PrivateConfig struct {
	NaverAPI struct {
		clientCode string `json: "clientCode"`
		secretCode string `json: "secretCode"`
	}
}
