package main

type Requirement struct {
  Name        string     `json:"name"`
	VersionReq  string     `json:"version_requirement"`
}

// Metadata is populated by the metadata.json file in the tar.gz file
type Metadata struct {
	Name          string            `json:"name"`
	Version       string            `json:"version"`
	Author        string            `json:"author"`
	Licence       string            `json:"license"`
  Summary       string            `json:"summary"`
  Source        string            `json:"source"`
  Requirements  []Requirement    `json:"requirements"`
  Dependencies  []Requirement    `json:"dependencies"`
  OSSupport     []struct {
                  Name          string     `json:"operatingsystem"`
                  Releases      []string   `json:"operatingsystemrelease"`
                } `json:"operatingsystem_support"`
}


type Owner struct {
  Username    string  `json:"username"`
}

type Pagination struct {
	Next        *string `json:"next"`
}

type ReleaseSummary struct {
  Uri         string  `json:"uri"`
  Version     string  `json:"version"`
  Supported   bool    `json:"supported"`
}
type ReleaseSummaries []ReleaseSummary
func (r ReleaseSummaries) Len() int { return len(r) }
func (r ReleaseSummaries) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r ReleaseSummaries) Less(i, j int) bool {
  res,_ := CompareVersion(r[i].Version, r[j].Version)
  return res < 0
}

type ModuleRelease struct {
  Uri    string      `json:"uri"`
  Module struct {
    Uri    string      `json:"uri"`
    Name   string      `json:"name"`
    Owner  Owner       `json:"owner"`
  }                  `json:"module"`
  Version  string    `json:"version"`
  Metadata Metadata  `json:"metadata"`
  Tags     []string  `json:"tags"`
  FileUri  string    `json:"file_uri"`
  FileMd5  string    `json:"file_md5"`
}

type ModulesResult struct {
	Uri             string            `json:"uri"`
  Name            string            `json:"name"`
  Owner           Owner             `json:"owner"`
  CurrentRelease  ModuleRelease     `json:"current_release"`
  Releases        ReleaseSummaries  `json:"releases"`
}

type ModulesResponse struct {
	Pagination struct {
    Next *string `json:"next"`
  } `json:"pagination"`
	Results  []ModulesResult `json:"results"`
}

type ReleasesResponse struct {
	Pagination struct {
    Next *string `json:"next"`
  } `json:"pagination"`
	Results  []ModuleRelease  `json:"results"`
}

