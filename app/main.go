package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/ohdaddyplease/dickts/app/config"
	dictionary "github.com/ohdaddyplease/dickts/app/database"
	server "github.com/ohdaddyplease/dickts/app/http"
	"github.com/ohdaddyplease/dickts/app/updater"
	"github.com/ohdaddyplease/dickts/pkg/logger"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"os"
)

var (
	updateTimeFirstDict  = flag.String("update_time_first_dictionary", "2s", "set update time for first dictionary")
	updateTimeSecondDict = flag.String("update_time_second_dictionary", "3s", "set update time for first dictionary")
)

func main() {
	var cfg config.Config
	cfgFile, err := os.Open("config/config.yml")
	if err != nil {
		log.Fatalf("Can't read DB config: %s", err)
	}
	defer func() {
		err = cfgFile.Close()
		if err != nil {
			log.Fatalf("Closing cfg file err: %s", err)
		}
	}()

	d := yaml.NewDecoder(cfgFile)
	if err = d.Decode(&cfg); err != nil {
		log.Fatalf("Can't parse DB config: %s", err)
	}

	l := logger.New(cfg.LogLevel, "")
	l.Info("Starting server...")

	flag.Parse()

	durationFirstDict, err := time.ParseDuration(*updateTimeFirstDict)
	if err != nil {
		l.Fatal("update_time_first_dictionary parsing error: %s", err)
	}
	durationSecondDict, err := time.ParseDuration(*updateTimeSecondDict)
	if err != nil {
		l.Fatal("update_time_second_dictionary parsing error: %s", err)
	}

	l.Info("[Prepare server] Connecting to database")
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Address, cfg.Port, cfg.UserName, cfg.Password, cfg.DatabaseName,
	)
	conn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		l.Fatal("Open conntion to database error: %s", err)
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			l.Fatal("Closing database connection error: %s", err)
		}
	}()

	l.Info("[Prepare server] Ping database")
	if err = conn.Ping(); err != nil {
		l.Fatal("Can't connect to database: %s", err)
	}

	l.Info("[Prepare server] Making query to database")

	var dict1 = make([]dictionary.Dictionary, 0, 0)
	var dict2 = make([]dictionary.Dictionary, 0, 0)

	l.Info("[Prepare server] Starting updater")
	upd := updater.New(conn, l)
	upd.Add(dict1, "dict_1", durationFirstDict)
	upd.Add(dict2, "dict_2", durationSecondDict)
	upd.Run()

	l.Info("[Prepare server] Registering handlers")

	http.Handle("/api/v1/dictionary", server.Dictionary{Logger: l, Updater: upd})

	sgnls := make(chan os.Signal)
	signal.Notify(sgnls, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt, os.Kill)
	go func() {
		<-sgnls
		l.Info("Shutdown server...")
		l.Info("Clearing dictionaries")
		upd.ClearDicts()
		os.Exit(1)
	}()

	if err = http.ListenAndServe(":81", nil); err != nil {
		l.Fatal("Starting server error: %s", err)
	}
}
