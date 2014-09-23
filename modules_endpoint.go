package main

import (
  "encoding/json"
	"fmt"
	"net/http"
  "strings"
  "bytes"
  "sort"
)


func ModulesHandler(w http.ResponseWriter, r *http.Request, modulePath string) {
  // read all of the release information for all of the files in the modPath
  releases, err := ReadModuleReleases(modulePath)
  if err != nil {
    http.Error(w, "Unable to read module files", 500)
  }
  if r.URL.Path == "/v3/modules" || r.URL.Path == "/v3/modules/" {
    queryParams := r.URL.Query()
    if queryParams["query"] != nil {
      queryTerm := queryParams["query"][0]
      moduleResults := make(map[string]ModulesResult)

      // go through each release (file) and determine if it is matched by the query
      for _, release := range releases {
        if strings.Contains(release.Metadata.Name, queryTerm) {
          // go ahead and generate the release summary.
          rs := ReleaseSummary {
            Uri: fmt.Sprintf("/v3/modules/%s-%s",
                             release.Module.Name,
                             release.Version),
            Version: release.Version,
            Supported: true,
          }
          existingResult, exists := moduleResults[release.Metadata.Name]
          if !exists {
            //new module
            var results = ModulesResult {
              Uri:   fmt.Sprintf("/v3/modules/%s", release.Metadata.Name),
              Name:  release.Module.Name,
              Owner: Owner { strings.Split(release.Metadata.Name, "-")[0] },
              CurrentRelease: release,
              Releases: make(ReleaseSummaries, 0),
            }
            results.Releases = append(results.Releases, rs)
            sort.Sort(sort.Reverse(results.Releases))
            moduleResults[release.Metadata.Name] = results
          } else { // existing module
            res, err := CompareVersion(release.Version, existingResult.CurrentRelease.Version)
            if err != nil { http.Error(w, "invalid version number", 500) }
            if res > 0 {
              //replace the currentVersion
              existingResult.CurrentRelease = release
              existingResult.Releases = append(existingResult.Releases, rs)
              sort.Sort(sort.Reverse(existingResult.Releases))
            } else {
              //just another lesser version, just add it to the release summaries
              existingResult.Releases = append(existingResult.Releases, rs)
              sort.Sort(sort.Reverse(existingResult.Releases))
            }
            //Store the existingResult back
            moduleResults[release.Metadata.Name] = existingResult
          }
        }
      }

      var resp ModulesResponse  // the empty response
      res := make([]ModulesResult, 0)
      for  _, v := range moduleResults {
        res = append(res, v)
      }
      resp.Results = res

      // write the response json out to the http Response
      b, err := json.Marshal(resp)
      if err != nil {
        http.Error(w, "json Marshalling error", 500)
      }
      var buff bytes.Buffer
      json.Indent(&buff, b, "", "\t")
      w.Header().Set("Content-Type", "application/json; charset=utf-8")
      fmt.Fprintf(w, string(buff.Bytes()))
    } else {
      http.Error(w, "I'm not going to list all of the modules.", 500)
    }
  }else{
    // find a specific one
    moduleName := strings.Split(r.URL.Path, "/")[3]
    fmt.Fprintf(w, "modulename: %s", moduleName)
  }
}
