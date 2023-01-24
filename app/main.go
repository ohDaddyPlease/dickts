package main

import (
	"github.com/ohdaddyplease/dickts/config"
	"gopkg.in/yaml.v3"

	"log"
	"os"
)

func main() {
	var dbCfg config.DB
	cfg, err := os.Open("config/db.yml")
	if err != nil {
		log.Fatal("Can't read DB config: ", err)
	}
	defer cfg.Close()
	d := yaml.NewDecoder(cfg)
	if err := d.Decode(&dbCfg); err != nil {
		log.Fatal("Can't parse DB config: ", err)
	}
}
