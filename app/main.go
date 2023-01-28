package main

import (
	"flag"
	"github.com/ohdaddyplease/dickts/app/config"
	dictionary "github.com/ohdaddyplease/dickts/app/http"
	"gopkg.in/yaml.v3"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"log"
	"os"
)

var (
	updateTimeFirstDict  = flag.String("update_time_first_dictionary", "5m", "set update time for first dictionary")
	updateTimeSecondDict = flag.String("update_time_second_dictionary", "3m", "set update time for first dictionary")
)

func main() {
	log.Println("Starting server...")

	flag.Parse()

	durationFirstDict, err := time.ParseDuration(*updateTimeFirstDict)
	if err != nil {
		log.Fatalf("update_time_first_dictionary parsing error: %s", err)
	}
	durationSecondDict, err := time.ParseDuration(*updateTimeSecondDict)
	if err != nil {
		log.Fatalf("update_time_second_dictionary parsing error: %s", err)
	}

	_, _ = durationFirstDict, durationSecondDict

	sgnls := make(chan os.Signal)
	signal.Notify(sgnls, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt, os.Kill)
	go func() {
		<-sgnls
		log.Println("Shutdown server")
		os.Exit(1)
	}()

	log.Println("[Prepare server] Registering handlers")
	http.Handle("/api/v1/dictionary", dictionary.Dictionary{})

	log.Println("[Prepare server] Opening config")
	var dbCfg config.DB
	cfg, err := os.Open("config/db.yml")
	if err != nil {
		log.Fatalf("Can't read DB config: %s", err)
	}
	defer func() {
		err = cfg.Close()
		if err != nil {
			log.Fatalf("Closing cfg file err: %s", err)
		}
	}()

	log.Println("[Prepare server] Parsing config")
	d := yaml.NewDecoder(cfg)
	if err = d.Decode(&dbCfg); err != nil {
		log.Fatalf("Can't parse DB config: %s", err)
	}

	if err = http.ListenAndServe(":81", nil); err != nil {
		log.Fatalf("Starting server error: %s", err)
	}
}
