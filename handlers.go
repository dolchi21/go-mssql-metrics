package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dolchi21/first-go-app/util"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Handlers struct {
	logger     *log.Logger
	collectors *Collectors
}

func (h *Handlers) Info(w http.ResponseWriter, r *http.Request) {
	db := NewDBConn()
	rows, err := db.Query(SQLSysTable)
	must(err)
	rs := util.RowsToMap(rows)
	json.NewEncoder(w).Encode(rs)
}

func (h *Handlers) Metrics(w http.ResponseWriter, r *http.Request) {
	h.collectors.Update()
	promhttp.Handler().ServeHTTP(w, r)
}
