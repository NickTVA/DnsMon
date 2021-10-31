package main

import (
	"encoding/json"
	"github.com/NickTVA/DnsMon/dnsresolver"
	"github.com/NickTVA/DnsMon/httphandlers"
	"github.com/newrelic/go-agent/v3/newrelic"
	"os"
	"time"
)

var app *newrelic.Application

func main() {

	hostnames := make([]string, 0)
	hostnames = append(hostnames, "www.motortrend.com")
	hostnames = append(hostnames, "www.google.com")
	hostnames = append(hostnames, "www.dfsfdsafdsd.com")

	app = initNewRelic()
	hostname, _ := os.Hostname()

	go httphandlers.SetupHTTP(app)

	app.RecordCustomEvent("DNSMonStarted", map[string]interface{}{
		"hostname": hostname,
	})

	go monitorDNS(hostnames)
	tickMonitor()
}

func monitorDNS(hostnames []string) {

	for true {

		for _, hostname := range hostnames {
			dnsinfo := dnsresolver.GetDNSInfo(hostname)
			app.RecordCustomEvent("dns_mon", dnsinfo)
			bytes, _ := json.Marshal(dnsinfo)

			println(string(bytes))

		}

		time.Sleep(90 * time.Second)

	}
}

func tickMonitor() {
	hostname, _ := os.Hostname()
	monitorName := os.Getenv("MONITOR_NAME")

	for true {
		time.Sleep(60 * time.Second)
		event := map[string]interface{}{
			"hostname":     hostname,
			"monitor_name": monitorName,
		}

		println("Tick")

		app.RecordCustomEvent("DNSMonTick", event)
	}
}

func initNewRelic() *newrelic.Application {
	newrelicKey := os.Getenv("NEWRELIC_KEY")
	if len(newrelicKey) < 1 {
		print("Must set NEWRELIC_KEY with NewRelic license key")
		os.Exit(-1)
	}

	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("DNSMon"),
		newrelic.ConfigLicense(newrelicKey),
		newrelic.ConfigDistributedTracerEnabled(true),
	)

	if app == nil || err != nil {
		print("NewRelic Not Initialized")
		os.Exit(-1)
	}
	return app
}
