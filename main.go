package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

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
	teamName := os.Getenv("TEAM")
	config, err := summary.SetupConfig(refreshIntervalString, groupsJSON, hostsJSON, skipSSLValidationString, teamName)
	if err != nil {
		log.Fatal(err)
	}
	config.Templates = templates

	server := summary.CreateServer(config)
	router := server.Start()

	fmt.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
