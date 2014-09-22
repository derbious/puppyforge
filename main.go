package main

import (
	"log"
	"net/http"
	//"os"
	"path/filepath"
  "flag"
)


// Initialize the application. Read the config options in the following hierarchy
//   command line args         <= highest priority
//   supplied config file
//   /etc/puppyforge.conf  <= lowest priority
func ReadConfig() (string, string) {
  // Read the flags, determine if 
  port := flag.String("port", "80", "The port number to listen on.")
  modpath := flag.String("modulepath", "", "The path to search for the modules.")
  configfile := flag.String("config", "/etc/puppyforge.conf", "The location of the config file.")
  flag.Parse()
  log.Println(*port, *modpath, *configfile)
  if configfile != "/etc/puppyforge.conf" {
  }
  return *port, *modpath
}


func main() {
  port, modulepath := ReadConfig()

  log.Fatal("testing")
	if len(port) == 0 {
		log.Fatal("Missing PORT environment variable")
	}
	if len(modulepath) == 0 {
		log.Fatal("Missing MODULEPATH environment variable")
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

	http.ListenAndServe(":"+port, nil)
}
