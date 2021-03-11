package config

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
)

const defaultConfigFile = "/tmp/harbor-operator/config.yaml"

type ConfigFile struct {
	DB struct {
		Type  string `json:"type"`
		PGSQL PGSQL  `json:"pgsql"`
	} `json:"db"`
	Redis struct {
		Addr     string `json:"redisAddr"`
		Password string `json:"redisPassword"`
	} `json:"redis"`
	Minio struct {
		Region   string `json:"region"`
		EndPoint string `json:"endpoint"`
		Proto    string `json:"proto"`
	} `json:"minio"`
	Kubernetes struct {
		Apiserver string `json:"kubeApiserver"`
		Token     string `json:"kubeToken"`
	} `json:"kubernetes"`
	HarborHelmPath struct {
		HarborV213 string `json:"harbor213"`
	} `json:"harborHelmPath"`
}

type PGSQL struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	SSLMode  string `json:"sslMode"`
}

// 解析配置文件config.yaml，生成配置ConfigFile
func ParseConfigFile() (*ConfigFile, error) {
	configData, err := ioutil.ReadFile(defaultConfigFile)
	if err != nil {
		return nil, err
	}

	var c ConfigFile
	if err := yaml.Unmarshal(configData, &c); err != nil {
		return nil, err
	}

	return &c, nil
}
