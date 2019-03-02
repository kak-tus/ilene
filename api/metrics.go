package api

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func (a *Type) metrics() {
	cnt := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ami_queue_items_total",
		Help: "Ami queue lengths",
	}, []string{"queue", "type"})

	a.updateMetrics(cnt)

	go func() {
		for {
			<-a.tick.C
			a.updateMetrics(cnt)
		}
	}()
}

func (a *Type) updateMetrics(cnt *prometheus.GaugeVec) {
	obj.lock.Lock()

	info := a.mdl.ListStreams()

	for _, inf := range info.List {
		cnt.With(prometheus.Labels{"queue": inf.Key, "type": inf.Type}).Set(float64(inf.Length))
	}

	obj.lock.Unlock()
}
