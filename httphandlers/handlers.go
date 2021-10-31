package httphandlers

import (
	"fmt"
	"github.com/newrelic/go-agent/v3/newrelic"
	"net/http"
	"os"
	"time"
)

var app *newrelic.Application

func SetupHTTP(nrapp *newrelic.Application) {
	app = nrapp
	http.HandleFunc("/healthz", healthz)
	http.HandleFunc("/kill", kill)
	http.ListenAndServe(":8000", nil)
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
