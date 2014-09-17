package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
  "fmt"
  "strconv"
  "math"
)


func ReadModuleReleases(modulePath string) ([]ModuleRelease, error) {
  releases := make([]ModuleRelease, 0)
  //see if this directory exists.
  files, err := ioutil.ReadDir(modulePath)
  if err != nil {
    return releases, err
  }

  for _, file := range files {
    if strings.HasSuffix(file.Name(), ".tar.gz") {
      m, err := ExtractMetadata(file, modulePath)
      if err != nil {
        log.Fatal("Error reading file %s: %s\n", file.Name(), err)
      }
      splitName := strings.Split(m.Name, "-")
      modUser := splitName[0]
      modName := splitName[1]

      // gen the Md5 sum
      md5, err := Checksum(filepath.Join(modulePath, file.Name()))
      if err != nil {
        log.Fatal("checksum failed")
      }
      var release ModuleRelease
      release.Uri = fmt.Sprintf("/v3/releases/%s-%s", m.Name, m.Version)
      release.Module.Uri = fmt.Sprintf("/v3/modules/%s", m.Name)
      release.Module.Name = modName
      release.Module.Owner = Owner { modUser }
      release.Version = m.Version
      release.Metadata = m
      release.Tags = make([]string, 0)
      release.FileUri = fmt.Sprintf("/v3/files/%s-%s.tar.gz", m.Name, m.Version)
      release.FileMd5 = md5

      releases = append(releases, release)
    }
  }
  return releases, nil
}

//Extract metadata from module file
func ExtractMetadata(module os.FileInfo, path string) (Metadata, error) {
	moduleFile := filepath.Join(path, module.Name())

  var m = Metadata{}

  //open the file, dealing with all of the errors
  fi, err := os.Open(moduleFile)
  if err != nil { return m, err }
  defer fi.Close()

  gzr, err := gzip.NewReader(fi)
  if err != nil { return m, err}
  defer gzr.Close()

  tr := tar.NewReader(gzr)

  // Iterate through the files in the archive.
  for {
    header, err := tr.Next()
    if err == io.EOF {
      break
    }
    if err != nil {
      return m, err
    }
    if strings.HasSuffix(header.Name, "/metadata.json") {
      //extract info
      data, err := ioutil.ReadAll(tr)
      if err != nil { continue }
      json.Unmarshal(data, &m)
      return m, nil
    }
  }
  // we have not found the manifest.json file
  return m, errors.New(fmt.Sprintf("No manifest.json found in %s", moduleFile))
}


func MaxOf(n []string) (uint32, error) {
  var m, max uint32
  var tmp int64
  tmp,  err := strconv.ParseInt(n[0], 10, 32)
  if err != nil { return 0, err }
  max = uint32(tmp)
  for i := 1; i<len(n); i++ {
    tmp, err := strconv.ParseInt(n[i], 10, 32)
    if err != nil { return 0, err }
    m = uint32(tmp)
    if m > max {
      max = m
    }
  }
  return max, nil
}


func CompareVersion(a, b string) (int, error) {
  if a == b {
    return 0, nil
  }else{
    aParts := strings.Split(a, ".")
    bParts := strings.Split(b, ".")
    // find the max int.
    nums := append(aParts, bParts...)
    mb, err := MaxOf(nums)
    if err != nil { return 0, err }
    var base uint32 = 10
    if mb > base {
      base = mb+1
    }
    maxParts := len(aParts)
    if len(bParts) > maxParts{
      maxParts = len(bParts)
    }

    // calculate A value
    var aVal float64 = 0.0
    for i := 0; i< len(aParts); i++ {
      exponent := maxParts-(i+1)
      tmp,_ := strconv.ParseInt(aParts[i], 10, 32)  //ignore error, it has already pasese
      aVal += float64(tmp)*math.Pow(float64(base), float64(exponent))
    }
    // calculate B value
    var bVal float64 = 0.0
    for i := 0; i< len(bParts); i++ {
      exponent := maxParts-(i+1)
      tmp,_ := strconv.ParseInt(bParts[i], 10, 32)  //ignore error, it has already pasese
      bVal += float64(tmp)*math.Pow(float64(base), float64(exponent))
    }
    if aVal > bVal {
      return 1, nil
    }else if aVal < bVal {
      return -1, nil
    }else{
      return 0, nil
    }
  }
}

