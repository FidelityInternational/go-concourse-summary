package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/FidelityInternational/go-concourse-summary/concourse"
)

var templates *template.Template

func init() {
	templates = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	hostsJSON := os.Getenv("HOSTS")
	groupsJSON := os.Getenv("CS_GROUPS")
	skipSSLValidationString := os.Getenv("SKIP_SSL_VALIDATION")
	refreshIntervalString := os.Getenv("REFRESH_INTERVAL")
	config, err := summary.SetupConfig(refreshIntervalString, groupsJSON, hostsJSON, skipSSLValidationString)
	if err != nil {
		log.Fatal(err)
	}
	config.Templates = templates

	router := mux.NewRouter()
	router.HandleFunc("/", config.Index)
	router.HandleFunc("/host/{host}", config.HostIndex)
	router.HandleFunc("/group/{group}", config.GroupIndex)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./assets/")))
	fmt.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
