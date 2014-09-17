package main

import (
  "encoding/json"
	"fmt"
	//"log"
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
      moduleName := queryParams["query"][0]

      var releaseSummaries ReleaseSummaries
      modulesLatest := make(map[string]string)
      var currentRelease ModuleRelease

      for _, release := range releases {
        // if the version is greater than the current release, upgadea
        // the current release. otherwise, add a new release summary.
        if release.Metadata.Name == moduleName {
          rs := ReleaseSummary {
            Uri: fmt.Sprintf("/v3/modules/%s-%s",
                             release.Module.Name,
                             release.Version),
            Version: release.Version,
            Supported: true,
          }

          if modulesLatest[release.Module.Name] == "" {
            //new module
            modulesLatest[release.Module.Name] = release.Version
            currentRelease = release
          } else {
            res, err := CompareVersion(release.Version, currentRelease.Version)
            if err != nil { http.Error(w, "invalid version number", 500) }
            if res > 0 {
              //replace the currentVersion
              modulesLatest[release.Module.Name] = release.Version
              currentRelease = release
            }
          }
          //just another lesser version
          releaseSummaries = append(releaseSummaries, rs)
        }
      }
      var resp ModulesResponse  // the empty response
      sort.Sort(sort.Reverse(releaseSummaries))
      var results = ModulesResult {
        Uri:   fmt.Sprintf("/v3/modules/%s", currentRelease.Metadata.Name),
        Name:  currentRelease.Module.Name,
        Owner: Owner { strings.Split(currentRelease.Metadata.Name, "-")[0] },
        CurrentRelease: currentRelease,
        Releases: releaseSummaries,
      }
      resp.Results = []ModulesResult{results,}

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
