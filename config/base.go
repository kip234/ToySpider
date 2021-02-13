package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const ConfigPath = "config.json"

type Config struct {
	Directory string //`json:"directory"`
	LinkHead string //`json:"linkhead"`
	RxImg string	//`json:"rximg"`
	RxHtml string
	RxCss string
	RxJs string
	Goal string
	Link string
}

func (c *Config)Init() (err error) {
	var file *os.File
	file,err=os.Open(ConfigPath)
	var buf []byte
	buf,err=ioutil.ReadAll(file)
	err=json.Unmarshal(buf,c)
	return
}