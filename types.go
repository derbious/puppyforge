package main


// Define all of the types that will be using in other files
type Dependencies struct {
	Name         string `json:"name"`
	VersionRange string `json:"version_requirement"`
}

// This is the required amount of metadata for any puppet forge module
type Metadata struct {
	Name         string         `json:"name"`
	Version      string         `json:"version"`
	Author       string         `json:"author"`
	Licence      string         `json:"license"`
  Summary      string         `json:"summary"`
  Source       string         `json:"source"`
	Dependencies []Dependencies `json:"dependencies"`
}

type Owner struct {
  Username    string  `json:"username"`
}

type Release struct {
  Uri       string  `json:"uri"`
  Version   string  `json:"version"`
  Supported bool    `json:"supported"`
}

type CurrentRelease struct {
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

type Result struct {
	Uri      string    `json:"uri"`
  Name     string    `json:"name"`
  Owner    Owner     `json:"owner"`
  CurrentRelease     `json:"current_release"`
  Releases []Release `json:"releases"`
}

type Pagination struct {
	Next *string `json:"next"`
}

type Response struct {
	Pagination struct {
    Next *string `json:"next"`
  } `json:"pagination"`
	Results  []Result `json:"results"`
}
