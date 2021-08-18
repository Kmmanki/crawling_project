package utils

import (
	"encoding/json"
	"log"
	"os"
)

func LoadConfig() (PrivateConfig, error) {
	var config PrivateConfig
	file, err := os.Open("privateConfig.json")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
	return config, err

}

func ErrChecker(err error) {
	if err != nil {
		log.Fatalln("err: ", err)
	}
}
