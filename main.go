package main

import (
	"log"
	"net/http"
	"path/filepath"
)


func main() {
  port, modulepath, err := InitConfig()
  if err != nil {
    log.Fatal(err)
  }

	log.Println("Starting puppyforge on port", port, "serving modules from", modulepath)

  // declare the webservice endpoints
	http.HandleFunc("/v3/modules", func(w http.ResponseWriter, r *http.Request) {
		ModulesHandler(w, r, modulepath)
	})
  http.HandleFunc("/v3/modules/", func(w http.ResponseWriter, r *http.Request) {
    ModulesHandler(w, r, modulepath)
  })

	http.HandleFunc("/v3/files/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(modulepath, r.URL.Path[10:]))
	})
	http.HandleFunc("/v3/releases", func(w http.ResponseWriter, r *http.Request) {
		ReleasesHandler(w, r, modulepath)
	})

	err = http.ListenAndServe(":"+port, nil)
  if err != nil {
    log.Fatal(err)
  }
}
