package main

import (
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

// dial initates a connection to cloud storage provider
func dial(kind string, config stow.Config) (stow.Location, error) {
	return stow.Dial(kind, config)
}
