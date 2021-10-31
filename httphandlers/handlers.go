package httphandlers

import (
	"fmt"
	"github.com/newrelic/go-agent/v3/newrelic"
	"net/http"
	"os"
	"strconv"
	"time"
)

var app *newrelic.Application

func SetupHTTP(nrapp *newrelic.Application) {
	app = nrapp
	http.HandleFunc("/healthz", healthz)
	http.HandleFunc("/kill", kill)

	port := 8000
	println("Setting up http connection to handle /healthz on port: " + strconv.Itoa(port))

	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func healthz(w http.ResponseWriter, r *http.Request) {

	txn := app.StartTransaction("healthz")
	defer txn.End()
	txn.SetWebRequestHTTP(r)
	txn.SetWebResponse(w)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}

func kill(w http.ResponseWriter, r *http.Request) {
	txn := app.StartTransaction("kill")
	defer txn.End()
	txn.SetWebRequestHTTP(r)
	txn.SetWebResponse(w)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "exiting...")
	go exitApplication()
}

func exitApplication() {

	time.Sleep(1 * time.Second)
	os.Exit(0)
}
