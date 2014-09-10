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
)


// Read all of the metadata out of all of the .tar.gz files in the modulePath
func ReadMetadata(modulePath string) ([]Metadata, error) {
  metadata := make([]Metadata, 0)
  log.Println("in func ReadMetadata")
  //see if this directory exists.
  files, err := ioutil.ReadDir(modulePath)
  if err != nil {
    log.Println(err)
    return metadata, err
  }
  for _, file := range files {
    if strings.HasSuffix(file.Name(), ".tar.gz") {
      m, err := ExtractMetadata(file, modulePath)
      if err != nil {
        log.Printf("Error reading file %s: %s\n", file.Name(), err)
      }else{
        log.Println(m)
      }
      metadata = append(metadata, m)
    }
  }
  return metadata, nil
}


// ListModules returns all tar.gz files
//func ListModules(path string) []Metadata {
//	var result []Metadata
//	files, err := ioutil.ReadDir(path)
//	if err != nil {
//		log.Println(err)
//	}
//	for _, file := range files {
//		if strings.HasSuffix(file.Name(), ".tar.gz") {
//			err := ExtractMetadata(file, path)
//			if err != nil {
//				log.Println(err)
//				continue
//			}
//			metadata, err := readMetadata(filepath.Join(path, file.Name()+".metadata"))
//			if err != nil {
//				log.Println(err)
//				continue
//			}
//			result = append(result, metadata)
//		}
//	}
//	return result
//}

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
      log.Println("found the manifest!")
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


// Iterate through the files in the archive.
//	for {
//		hdr, err := tr.Next()
//		if err == io.EOF {
//			return errors.New("metadata.json not found in " + moduleFile)
//		}
//		if err != nil {
//			return err
//		}
//		// Found metadata.json, no need to read any further.
//		if hdr.Name == strings.TrimRight(module.Name(), "tar.gz")+"/metadata.json" {
//			f, err := os.Create(metadataPath)
//			defer f.Close()
//			if err != nil {
//				return err
//			}
//			io.Copy(f, tr)
//			return nil
//		}
//	}
//}

//func readMetadata(file string) (Metadata, error) {
//	var m Metadata
//	data, err := ioutil.ReadFile(file)
//	if err != nil {
//		return m, err
//	}
//	json.Unmarshal(data, &m)
//	return m, nil
//}

// ListModules returns all tar.gz files
//func ListModules(path string) []Metadata {
//	var result []Metadata
//	files, err := ioutil.ReadDir(path)
//	if err != nil {
//		log.Println(err)
//	}
//	for _, file := range files {
//		if strings.HasSuffix(file.Name(), ".tar.gz") {
//			err := ExtractMetadata(file, path)
//			if err != nil {
//				log.Println(err)
//				continue
//			}
//			metadata, err := readMetadata(filepath.Join(path, file.Name()+".metadata"))
//			if err != nil {
//				log.Println(err)
//				continue
//			}
//			result = append(result, metadata)
//		}
//	}
//	return result
//}

