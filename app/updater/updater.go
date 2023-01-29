package updater

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	dictionary "github.com/ohdaddyplease/dickts/app/database"
	"github.com/ohdaddyplease/dickts/pkg/logger"
	"time"
)

type JobI interface {
	clearDictionary()
}

type Job struct {
	Dictionary []dictionary.Dictionary
	Duration   time.Duration
	Table      string
}

func (j *Job) clearDictionary() {
	j.Dictionary = []dictionary.Dictionary{}
}

type UpdaterI interface {
	Add(map[string][]dictionary.Dictionary, string, *sql.DB, time.Duration)
	Run()
	ClearDicts()
}

type Updater struct {
	Jobs   map[string]Job
	Logger logger.LoggerI
	Query  *sql.DB
}

func New(q *sql.DB, l logger.LoggerI) *Updater {
	return &Updater{Jobs: make(map[string]Job, 0), Logger: l, Query: q}
}

func (u *Updater) Add(dict []dictionary.Dictionary, table string, d time.Duration) {
	u.Logger.Info("[Prepare server] Adding job for updater")
	u.Jobs[table] = Job{
		Dictionary: dict,
		Duration:   d,
		Table:      table,
	}
	u.Logger.Debug("add updater for dict: %p\n", &dict)

}

func (u *Updater) Run() {
	for jobName, job := range u.Jobs {
		job := job
		go func(j Job, jobName string) {
			u.Logger.Debug("run job %s", j.Table)
			var dict dictionary.Dictionary
			for range time.Tick(j.Duration) {
				q, err := u.Query.Query(fmt.Sprintf("SELECT * FROM %s", j.Table))
				if err != nil {
					u.Logger.Fatal("SQL error: %s", err)
				}

				u.Logger.Debug("reset dict: %p\n", &j.Dictionary)
				j.clearDictionary()
				dicts := make([]dictionary.Dictionary, 0)
				for q.Next() {
					err := q.Scan(&dict.Name, &dict.Description, &dict.Updated)
					if err != nil {
						u.Logger.Error("Can't update dictionary")
						continue
					}
					u.Logger.Debug("Built dict: %v", dict)
					dicts = append(dicts, dict)
				}
				job.Dictionary = dicts
				u.Jobs[jobName] = job
			}
		}(job, jobName)
	}
}

func (u *Updater) ClearDicts() {
	for _, job := range u.Jobs {
		job.clearDictionary()
	}
}
