package main

import (
  "encoding/json"
	"fmt"
	"net/http"
	"strings"
  "bytes"
)

func ReleasesHandler(w http.ResponseWriter, r *http.Request, modulePath string) {
	moduleName := r.URL.Query().Get("module")
  if moduleName == "" || !strings.Contains(moduleName, "-") {
    http.Error(w, "request must be /v3/releases?module=user-module", 400)
    return
  }
  modnameArray := strings.Split(moduleName, "-")
	user := modnameArray[0]
	mod := modnameArray[1]
  validReleases := make([]ModuleRelease, 0)

	releases, err := ReadModuleReleases(modulePath)
  if err != nil {
    http.Error(w, "Could not read module metadata.", 500)
  }

  for _, rel := range releases {
    if rel.Module.Name == mod && rel.Module.Owner.Username == user {
      validReleases = append(validReleases, rel)
    }
  }
	response := new(ReleasesResponse)
	response.Pagination = Pagination{Next: nil}
  response.Results = validReleases

	jsonData, _ := json.Marshal(response)
  var buff bytes.Buffer
  json.Indent(&buff, jsonData, "", "\t")
  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  fmt.Fprintf(w, string(buff.Bytes()))
}
