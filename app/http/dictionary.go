package http

import (
	"fmt"
	"github.com/gorilla/schema"
	"log"
	"net/http"
	"time"
)

type Dictionary struct {
	Name string `json:"name"`
}

func (d Dictionary) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var updateTime time.Time
	decoder := schema.NewDecoder()
	err := decoder.Decode(&d, r.URL.Query())
	if err != nil {
		log.Fatalf("Decode params error: %s", err)
	}
	_, err = w.Write([]byte(fmt.Sprintf("%s has updated: %s", d.Name, updateTime)))
	if err != nil {
		log.Fatalf("Request error: %s", err)
	}
}
