package utils

import (
	"crypto/md5"
	"encoding/hex"
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

func GetMD5Hash(title string, postDate string) string {
	hash := md5.Sum([]byte(title + postDate))
	return hex.EncodeToString(hash[:])
}
