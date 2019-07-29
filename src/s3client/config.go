package s3client

import (
	"encoding/json"
	"io/ioutil"
)

type S3Config struct {
	ID       string `json:"id"`
	Secret   string `json:"secret"`
	EndPoint string `json:"endPoint"`
	Region   string `json:"region"`
}

func Load(file string) (config S3Config, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(file); err != nil {
		return config, err
	}
	if err = json.Unmarshal(data, &config); err != nil {
		return config, err
	} else {
		return config, nil
	}
}
