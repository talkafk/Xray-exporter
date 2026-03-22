package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metrics struct {
	userDownLinks prometheus.GaugeVec
	userUpLinks   prometheus.GaugeVec
}

func fetchStats(address string) *map[string]interface{} {
	resp, err := http.Get("http://" + address + "/debug/vars")
	if err != nil {
		log.Println("Error fetching stats:", err)
		return nil
	}
	defer resp.Body.Close()
	var stats map[string]interface{}

	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		log.Println("Error decoding JSON:", err)
		return nil
	}
	xrayStats, ok := stats["stats"].(map[string]interface{})
	if !ok {
		log.Println("Error: stats[\"stats\"] is not a map[string]interface{}")
		return nil
	}
	userStats, ok := xrayStats["user"].(map[string]interface{})
	if !ok {
		log.Println("Error: stats[\"stats\"] is not a map[string]interface{}")
		return nil
	}
	return &userStats

}

func newMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		userDownLinks: *promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
			Name: "xray_user_downlinks",
			Help: "The number of active user downlinks",
		}, []string{"user_id"}),
		userUpLinks: *promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{
			Name: "xray_user_uplinks",
			Help: "The number of active user uplinks",
		}, []string{"user_id"}),
	}
	return m
}

func recordMetrics(m *metrics, flagAdress *string) {
	go func() {
		for {
			userStats := fetchStats(*flagAdress)
			if userStats != nil {
				for name, data := range *userStats {
					dataMap, ok := data.(map[string]interface{})
					if !ok {
						log.Printf("Error: user data for %s is not a map[string]interface{}", name)
						continue
					}
					downlink, ok1 := dataMap["downlink"].(float64)
					uplink, ok2 := dataMap["uplink"].(float64)
					if !ok1 || !ok2 {
						log.Printf("Error: downlink or uplink for user %s is not a float64", name)
						continue
					}
					m.userDownLinks.With(prometheus.Labels{"user_id": name}).Set(float64(downlink))
					m.userUpLinks.With(prometheus.Labels{"user_id": name}).Set(float64(uplink))
				}
				time.Sleep(10 * time.Second)
			}
		}
	}()
}

func main() {
	flagAdress := flag.String("xray-endpoint", "localhost:11111", "The address of the Xray stats endpoint.")
	flagPort := flag.String("port", "9595", "The port to listen on for HTTP requests.")
	flag.Parse()

	reg := prometheus.NewRegistry()
	m := newMetrics(reg)
	recordMetrics(m, flagAdress)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	http.ListenAndServe(":"+*flagPort, nil)
}
