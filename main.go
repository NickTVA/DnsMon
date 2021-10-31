package main

import (
	"encoding/json"
	"github.com/NickTVA/DnsMon/dnsresolver"
	"github.com/NickTVA/DnsMon/httphandlers"
	"github.com/newrelic/go-agent/v3/newrelic"
	"os"
	"strconv"
	"time"
)

var app *newrelic.Application
var hostnames []string

func main() {

	hostnames = getHostsFromEnv()

	app = initNewRelic()
	println("Waiting up to a minute for NR connection")
	app.WaitForConnection(time.Minute)

	go httphandlers.SetupHTTP(app)

	app.RecordCustomEvent("DNSMonStarted", map[string]interface{}{
		"NumHosts": len(hostnames),
	})

	go monitorDNS(hostnames)
	tickMonitor()
}

func getHostsFromEnv() []string {
	hostnames := make([]string, 0)

	//Read hostnames from environment
	for i := 0; i < 100; i++ {

		hostname := os.Getenv("dns.host." + strconv.Itoa(i))
		if len(hostname) < 1 {
			continue
		}
		hostnames = append(hostnames, hostname)

	}

	if len(hostnames) < 1 {
		println("No hostnames in ENV")
		os.Exit(-1)
	}
	return hostnames
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
	monitorName := os.Getenv("MONITOR_NAME")

	for true {
		event := map[string]interface{}{
			"monitor_name": monitorName,
			"num_hosts":    len(hostnames),
		}

		println("Tick")

		app.RecordCustomEvent("DNSMonTick", event)
		time.Sleep(5 * time.Minute)

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
