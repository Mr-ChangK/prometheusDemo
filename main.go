package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/counter", CounterHandler)
	http.HandleFunc("/gauge", GaugeHandler)
	http.HandleFunc("/histogram", HistogramHandler)
	http.HandleFunc("/summary", SummaryHandler)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":18000", nil)
}

var counter = promauto.NewCounterVec(prometheus.CounterOpts{Name: "demo_counter"}, []string{"version"})

func CounterHandler(w http.ResponseWriter, r *http.Request) {
	version := r.PostFormValue("version")
	counter.With(prometheus.Labels{"version": version}).Inc()
}

var gauge = promauto.NewGauge(prometheus.GaugeOpts{Name: "demo_gauge"})

func GaugeHandler(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().Unix())
	randomNum := rand.Float64()
	gauge.Set(randomNum)
}

var histogram = promauto.NewHistogramVec(prometheus.HistogramOpts{Name: "demo_histogram", Buckets: []float64{0.2, 0.5, 1, 2, 5, 10, 30}}, []string{"version"})

func HistogramHandler(w http.ResponseWriter, r *http.Request) {
	version := r.PostFormValue("version")
	histogram.With(prometheus.Labels{"version": version}).Observe(float64(rand.Intn(30)))
}

var summary = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Name:       "demo_summary",
		Objectives: map[float64]float64{0.1: 0, 0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	},
	[]string{"version"},
)

func SummaryHandler(w http.ResponseWriter, r *http.Request) {
	version := r.PostFormValue("version")
	a := float64(rand.Intn(3000))
	fmt.Println(a)
	summary.With(prometheus.Labels{"version": version}).Observe(a)
}

func init() {
	prometheus.MustRegister(summary)
}
