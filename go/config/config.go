package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var (
	Conf Config
)

type Config struct {

	Mysql struct{
		Host string `yaml:"host"`
		Port int `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"mysql"`
	Redis struct{

	} `yaml:"redis"`

	Login struct{
		Expires int `yaml:"expires"`
		Key string `yaml:"key"`
	} `yaml:"login"`

}


func LoadEnv()  {
	f,err := ioutil.ReadFile("env.yaml")
	if err != nil{
		log.Fatalln(err)
	}
	err = yaml.Unmarshal(f, &Conf)
	if err != nil {
		log.Fatalln(err)
	}
}
