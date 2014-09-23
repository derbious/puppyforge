package main

import (
  "fmt"
  "flag"
  "os"
  "bufio"
  "strings"
  "errors"
)


func InitConfig() (string, string, error) {
  var clPort, clModuledir, conffile, port, moduledir string
  flag.StringVar(&clPort, "port", "", "Port that the webservice runs on")
  flag.StringVar(&clModuledir, "moduledir", "", "Directory containing the module files")
  flag.StringVar(&conffile, "config", "/etc/puppyforge.conf", "Config file location")
  flag.Parse()

  if clPort != "" && clModuledir != "" {
    // no need to read conf file.
    return clPort, clModuledir, nil
  }
  // read the config file.
  f, err := os.Open(conffile)
  if err != nil {
    return "", "", err
  }
  scanner := bufio.NewScanner(f)
  for scanner.Scan() {
    line := strings.Trim(scanner.Text(), " \t")
    fmt.Println(">>>", line)
    if strings.HasPrefix(line, "#") {
      continue  // comment line
    }else{
      splitlines := strings.Split(line, "=")
      if len(splitlines) == 2 {  // k=v pair here
        key := strings.Trim(splitlines[0], " \t")
        val := strings.Trim(splitlines[1], " \t")
        if key == "port" {
          port = val
        } else if key == "moduledir" {
          moduledir = val
        }
      }
    }
  }
  if err := scanner.Err(); err != nil {
    return "", "", err
  }

  // Overwrite the config options with the port
  if clPort != "" {
    port = clPort
  }
  if clModuledir != "" {
    moduledir = clModuledir
  }

  if port == "" || moduledir == "" {
    return "","", errors.New("Incorrect configuration")
  }else{
    return port, moduledir, nil
  }
}

