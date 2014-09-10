package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	port := os.Getenv("PORT")
	modulePath := os.Getenv("MODULEPATH")

	if len(port) == 0 {
		log.Fatal("Missing PORT environment variable")
	}
	if len(modulePath) == 0 {
		log.Fatal("Missing MODULEPATH environment variable")
	}

	log.Println("Starting go-puppet-forge on port", port, "serving modules from", modulePath)

  // declare the webservice endpoints
	http.HandleFunc("/v3/modules", func(w http.ResponseWriter, r *http.Request) {
		ModulesHandler(w, r, modulePath)
	})
  http.HandleFunc("/v3/modules/", func(w http.ResponseWriter, r *http.Request) {
    ModulesHandler(w, r, modulePath)
  })

	http.HandleFunc("/v3/files/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(modulePath, r.URL.Path[10:]))
	})
	http.HandleFunc("/v3/releases", func(w http.ResponseWriter, r *http.Request) {
		ReleasesHandler(w, r, modulePath)
	})

	http.ListenAndServe(":"+port, nil)
}
