package main

import (
  "encoding/json"
	"fmt"
	"log"
	"net/http"
  "strings"
  "bytes"
)

func ModulesHandler(w http.ResponseWriter, r *http.Request, modulePath string) {
  // read all of the metadata of all of the files in the modulePath
  metadata, err := ReadMetadata(modulePath)
  if err != nil {
    log.Println("ERROR:", err)
    return
  }
  log.Println(metadata)

  if r.URL.Path == "/v3/modules" || r.URL.Path == "/v3/modules/" {
    queryParams := r.URL.Query()
    if queryParams["query"] != nil {
      moduleName := queryParams["query"][0]
      log.Println("querying for a module", moduleName)


      var resp Response  // the empty response
      resp.Results = make([]Result, 0)

      for _, m := range metadata {
        if strings.Contains(m.Name, moduleName){
          user := strings.Split(m.Name, "-")[0]
          module := strings.Split(m.Name, "-")[1]
          moduleUri := fmt.Sprintf("/v3/modules/%s", m.Name)
          releaseUri := fmt.Sprintf("/v3/releases/%s-%s", m.Name, m.Version)
          // add a new result to the response
          var result Result

          //populate the owner info
          owner := Owner{user}

          // populate the current_release
          var cr CurrentRelease
          cr.Uri = releaseUri
          cr.Module.Uri = moduleUri
          cr.Module.Name = module
          cr.Module.Owner = owner
          cr.Version = m.Version
          cr.Metadata = m
          cr.Tags = make([]string, 0)
          cr.FileUri = fmt.Sprintf("v3/files/%s-%s.tar.gz", m.Name, m.Version)
          cr.FileMd5 = "foo_md5"

          // populate the Result struct
          result.Uri = moduleUri
          result.Name = module
          result.Owner = owner
          result.CurrentRelease = cr
          result.Releases = append(result.Releases,
            Release{
              Uri: releaseUri,
              Version: m.Version,
              Supported: true,
           })

          resp.Results = append(resp.Results, result)

        }
      }
      // write the response json out to the http Response
      b, err := json.Marshal(resp)
      if err != nil {
        log.Println(err)
      }
      var buff bytes.Buffer
      json.Indent(&buff, b, "", "\t")
      w.Header().Set("Content-Type", "application/json; charset=utf-8")
      fmt.Fprintf(w, string(buff.Bytes()))
    }
  }else{
    // find a specific one
    moduleName := strings.Split(r.URL.Path, "/")[3]
    fmt.Fprintf(w, "modulename: %s", moduleName)
  }
}
