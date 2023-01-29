package http

import (
	"fmt"
	"github.com/gorilla/schema"
	"github.com/ohdaddyplease/dickts/app/updater"
	"github.com/ohdaddyplease/dickts/pkg/logger"
	"net/http"
	"strings"
)

type Dictionary struct {
	Logger  logger.LoggerI `json:"-"`
	Name    string         `json:"name"`
	Updater *updater.Updater
}

func (d Dictionary) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	decoder := schema.NewDecoder()
	err := decoder.Decode(&d, r.URL.Query())
	if err != nil {
		d.Logger.Fatal("Decode params error: %s", err)
	}

	d.Logger.Debug("Jobs in http: %v", d.Updater.Jobs)
	job, ok := d.Updater.Jobs[d.Name]
	if ok {
		sb := strings.Builder{}
		for _, dictFromDb := range job.Dictionary {
			sb.WriteString(fmt.Sprintf("Name: %s\nDescription: %s\nUpdated: %s\n\n", dictFromDb.Name,
				dictFromDb.Description,
				dictFromDb.Updated.String()))
		}
		_, err = w.Write([]byte(sb.String()))
		if err != nil {
			d.Logger.Fatal("Write error: %s", err)
		}
	} else {
		_, err = w.Write([]byte(fmt.Sprintf("No dictionary with name %s", d.Name)))
		if err != nil {
			d.Logger.Fatal("Write error: %s", err)
		}
	}
}
