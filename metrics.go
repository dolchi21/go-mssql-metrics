package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Collectors struct {
	TableCount      *prometheus.GaugeVec
	TableTotalSpace *prometheus.GaugeVec
}

func (m *Collectors) Register(reg *prometheus.Registry) {
	var mustRegister func(...prometheus.Collector)
	if nil != reg {
		mustRegister = reg.MustRegister
	} else {
		mustRegister = prometheus.MustRegister
	}
	mustRegister(m.TableCount)
	mustRegister(m.TableTotalSpace)
}

func (m *Collectors) Update() {
	updateTableCount(m.TableCount)
	updateTableTotalSize(m.TableTotalSpace)
}

func NewMetrics() *Collectors {
	m := Collectors{
		TableCount: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "mssql_table_rows",
			Help: "Table row count",
		}, []string{"table", "schema"}),
		TableTotalSpace: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "mssql_table_total_space",
			Help: "Total KB reserved for table",
		}, []string{"table", "schema"}),
	}
	return &m
}

func updateTableCount(g *prometheus.GaugeVec) {
	db := NewDBConn()
	defer db.Close()
	rows, err := db.Query(SQLSysTableRowCount)
	must(err)
	for rows.Next() {
		var schema string
		var tableName string
		var value float64
		must(rows.Scan(&schema, &tableName, &value))
		g.WithLabelValues(tableName, schema).Set(value)
	}
}

func updateTableTotalSize(g *prometheus.GaugeVec) {
	db := NewDBConn()
	defer db.Close()
	rows, err := db.Query(SQLSysTableBytes)
	must(err)
	for rows.Next() {
		var schema, tableName string
		var value float64
		must(rows.Scan(&schema, &tableName, &value))
		g.WithLabelValues(tableName, schema).Set(value)
	}
}
