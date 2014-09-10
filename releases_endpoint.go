package main

import (
  "encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

func ReleasesHandler(w http.ResponseWriter, r *http.Request, modulePath string) {
	moduleName := r.URL.Query().Get("module")
	if !strings.Contains(moduleName, "-") {
		http.Error(w, "request must be /v3/releases?module=user-module", 400)
		return
	}
  modnameArray := strings.Split(moduleName, "-")
	user := modnameArray[0]
	mod := modnameArray[1]

	modules, err := ReadMetadata(filepath.Join(modulePath, user, mod))
  if err != nil {
    log.Println("Error reading metadata")
    http.Error(w, "Could not read module metadata.", 500)
  }
  log.Println(modules)
	// No modules, return minimal json response.
	if len(modules) == 0 {
		fmt.Fprintf(w, `{"pagination":{"next":null},"results":[]}`)
		return
	}
	response := new(Response)
	response.Pagination = Pagination{Next: nil}

	for _, metadata := range modules {
		//checksum, err := Checksum(filepath.Join(modulePath, user, mod, moduleName+"-"+metadata.Version+".tar.gz"))
		//if err != nil {
			// Unable to checksum modulefile, log and skip.
		log.Println(metadata)
		//	continue
		//}
		//var result = Result{
	//		Uri:     fmt.Sprintf("/v3/release/%s/%s", metadata.Name, metadata.Version),
		//	Version: metadata.Version,
		//	FileUri: fmt.Sprintf("/v3/files/%s/%s/%s-%s.tar.gz", user, mod, moduleName, metadata.Version),
		//	Md5:     checksum}
		//result.Metadata = metadata
		//response.Results = append(response.Results, result)
	}
	jsonData, _ := json.Marshal(response)
	fmt.Fprintf(w, string(jsonData))
}
