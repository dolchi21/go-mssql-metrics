package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

func main() {
	loadConfig()

	M := NewMetrics()
	M.Register(nil)

	h := Handlers{
		collectors: M,
	}

	http.HandleFunc("/", h.Info)
	http.HandleFunc("/metrics", h.Metrics)

	log.Fatal(http.ListenAndServe(viper.GetString("server.listen"), nil))

}

func must(err error) {
	if nil != err {
		panic(err)
	}
}

func loadConfig() {
	viper.SetDefault("server.listen", ":8081")

	viper.SetConfigFile("config.yml")
	err := viper.ReadInConfig()
	if nil != err {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}
