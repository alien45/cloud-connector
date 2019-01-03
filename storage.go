package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/graymeta/stow"
	_ "github.com/graymeta/stow/azure"
	_ "github.com/graymeta/stow/google"
	_ "github.com/graymeta/stow/s3"
)

// Connection describes connection details for supporeted cloud providers as well as query data
type Connection struct {
	Kind          string         `json:"kind"`
	ContainerName string         `json:"container_name"`
	ItemID        string         `json:"item_id"`
	ItemName      string         `json:"item_name"`
	ConfigMap     stow.ConfigMap `json:"config_map"`
	Pagination
}

// Dial initates a connection to cloud storage provider. Will close connection automatically.
func (c Connection) Dial() (loc stow.Location, err error) {
	loc, err = stow.Dial(c.Kind, c.ConfigMap)
	if err == nil {
		defer loc.Close()
	}
	return
}

func connect(w http.ResponseWriter, r *http.Request) (con Connection, loc stow.Location, success bool) {
	bytes, err := ioutil.ReadAll(r.Body)
	if respondIfError(err, w, fmt.Sprintf("Failed to read request body. Error: %v", err), err400) {
		return
	}
	return connectJSON(w, bytes)
}
func connectJSON(w http.ResponseWriter, jsonB []byte) (con Connection, loc stow.Location, success bool) {
	err = json.Unmarshal(jsonB, &con)
	if respondIfError(err, w, fmt.Sprintf("Valid JSON body required. Error: %v", err), err400) {
		return
	}

	log.Println("Dialing ", con.Kind)
	loc, err := con.Dial()
	if respondIfError(err, w, fmt.Sprintf("Connection to %s failed. Error: %v", con.Kind, err), err500) {
		return
	}
	success = true
	log.Println("Connected to", con.Kind, con.ContainerName)
	return
}

// Pagination stores data requried for Stow pagination
type Pagination struct {
	Cursor string `json:"cursor"`
	Count  int    `json:"count"`
}

// ContainersResult describes the API result containing list fo containers
type ContainersResult struct {
	Pagination
	Containers map[string]string `json:"containers"`
}

// ItemsResult describes the API result containing list for items/objects/files
type ItemsResult struct {
	Pagination
	Items []Item `json:"items"`
}

// Item describes the API result for a single item/object/file
type Item struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	Size     int64                  `json:"size"`
	URL      string                 `json:"url"`
	Metadata map[string]interface{} `json:"metadata"`
}
