package main

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/tsenart/vegeta/lib"
	"io/ioutil"
)

type Config struct {
	Rate     int // calls/sec
	Duration time.Duration
	Request  RequestConfig
}

type RequestConfig struct {
	Url    string
	Method string
	Body   string
	Header map[string][]string
}

func main() {
	configBytes, err := ioutil.ReadFile("test/fixtures/load-test/config.toml")
	if err != nil {
		panic(err)
	}

	cfg := new(Config)
	err = toml.Unmarshal(configBytes, &cfg)
	if err != nil {
		panic(err)
	}

	rate := vegeta.Rate{Freq: cfg.Rate, Per: time.Second}
	duration := cfg.Duration * time.Second
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: cfg.Request.Method,
		URL:    cfg.Request.Url,
		Body:   []byte(cfg.Request.Body),
	})
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	fmt.Printf("Starting load test for %d seconds with a frequency of %d on %s\n", cfg.Duration, cfg.Rate, cfg.Request.Url)
	for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Printf("Latencies: \n99th percentile: %s\nMean: %s\n", metrics.Latencies.P99, metrics.Latencies.Mean)
	fmt.Printf("\nRate of requests per second: %v\n", metrics.Rate)
	fmt.Printf("Total number of requests executed: %v\n", metrics.Requests)

}
