package main

import (
	"fmt"
	"github.com/innotech/hydra/vendors/github.com/coreos/etcd/third_party/github.com/rcrowley/go-metrics"
	"time"
)

func main() {
	r := metrics.NewRegistry()
	for i := 0; i < 10000; i++ {
		r.Register(fmt.Sprintf("counter-%d", i), metrics.NewCounter())
		r.Register(fmt.Sprintf("gauge-%d", i), metrics.NewGauge())
		r.Register(fmt.Sprintf("histogram-uniform-%d", i), metrics.NewHistogram(metrics.NewUniformSample(1028)))
		r.Register(fmt.Sprintf("histogram-exp-%d", i), metrics.NewHistogram(metrics.NewExpDecaySample(1028, 0.015)))
		r.Register(fmt.Sprintf("meter-%d", i), metrics.NewMeter())
	}
	time.Sleep(600e9)
}