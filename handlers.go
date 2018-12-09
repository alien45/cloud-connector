package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func getJSONStringHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[Request] %s", r.URL.Path)
	var data interface{}
	err = json.NewDecoder(r.Body).Decode(&data)
	if respondIfError(err, w, fmt.Sprintf("Valid JSON body required. Error: %v", err), err400) {
		return
	}
	b, err := json.Marshal(data)
	if respondIfError(err, w, fmt.Sprintf("Valid JSON body required. Error: %v", err), err500) {
		return
	}
	respondJSON(w, string(b), ok200)
}

func containersHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[Request] %s", r.URL.Path)
	con := Connection{}
	err = json.NewDecoder(r.Body).Decode(&con)
	if respondIfError(err, w, fmt.Sprintf("Valid JSON body required. Error: %v", err), err400) {
		return
	}

	log.Println("Dialing ", con.Kind)
	loc, err := dial(con.Kind, con.ConfigMap)
	if respondIfError(err, w, fmt.Sprintf("Connection to %s failed. Error: %v", con.Kind, err), err500) {
		return
	}
	defer loc.Close()

	containers, cursor, err := loc.Containers(con.Cursor, "", con.Count)
	if respondIfError(err, w, fmt.Sprintf("Failed to retrieve containers. Error: %v", err), err500) {
		return
	}

	result := ContainersResult{}
	result.Count = con.Count
	result.Cursor = cursor
	result.Containers = map[string]string{}
	for i := 0; i < len(containers); i++ {
		result.Containers[containers[i].ID()] = containers[i].Name()
	}
	respondJSON(w, result, ok200)
}

func itemsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[Request] %s", r.URL.Path)
	con := Connection{}
	err = json.NewDecoder(r.Body).Decode(&con)
	if respondIfError(err, w, fmt.Sprintf("Valid JSON body required. Error: %v", err), err400) {
		return
	}

	log.Printf("Dialing %s/%s\n", con.Kind, con.ContainerName)
	loc, err := dial(con.Kind, con.ConfigMap)
	if respondIfError(err, w, fmt.Sprintf("Connection to %s failed. %v", con.Kind, err), err500) {
		return
	}
	defer loc.Close()

	container, err := loc.Container(con.ContainerName)
	if respondIfError(err, w, fmt.Sprintf("Failed to retrieve container. Error: %v", err), err500) {
		return
	}

	items, cursor, err := container.Items("", con.Cursor, con.Count)
	result := ItemsResult{}
	result.Count = con.Count
	result.Cursor = cursor
	for i := 0; i < len(items); i++ {
		size, _ := items[i].Size()
		metadata, _ := items[i].Metadata()
		result.Items = append(result.Items, Item{
			ID:       items[i].ID(),
			Name:     items[i].Name(),
			Size:     size,
			URL:      items[i].URL().Host + items[i].URL().Path,
			Metadata: metadata,
		})
	}
	respondJSON(w, result, ok200)
}

func copyHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[Request] %s", r.URL.Path)
	con := struct {
		From Connection `json:"from"`
		To   Connection `json:"to"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&con)
	if respondIfError(err, w, fmt.Sprintf("Valid JSON body required. Error: %v", err), err400) {
		return
	}

	log.Printf("Dialing source: %s/%s", con.From.Kind, con.From.ContainerName)
	sourceLoc, err := dial(con.From.Kind, con.From.ConfigMap)
	if respondIfError(err, w, fmt.Sprintf("Connection to %s failed. %v", con.From.Kind, err), err500) {
		return
	}
	defer sourceLoc.Close()

	sourceContainer, err := sourceLoc.Container(con.From.ContainerName)
	if respondIfError(err, w, fmt.Sprintf("Failed to retrieve container. Error: %v", err), err500) {
		return
	}
	itemID := con.From.ItemID
	if itemID == "" {
		itemID = con.From.ItemName
	}
	item, err := sourceContainer.Item(itemID)

	log.Printf("Dialing destination: %s/%s", con.To.Kind, con.To.ContainerName)
	destLoc, err := dial(con.To.Kind, con.To.ConfigMap)
	if respondIfError(err, w, fmt.Sprintf("Connection to %s failed. %v", con.To.Kind, err), err500) {
		return
	}
	destContainer, err := destLoc.Container(con.To.ContainerName)
	if respondIfError(err, w, fmt.Sprintf("Failed to retrieve container. Error: %v", err), err500) {
		return
	}
	reader, err := item.Open()
	if respondIfError(err, w, fmt.Sprintf("Failed to open file. Error: %v", err), err500) {
		return
	}
	size, err := item.Size()
	if respondIfError(err, w, fmt.Sprintf("Failed to retrieve file size. Error: %v", err), err500) {
		return
	}
	metadata, err := item.Metadata()
	if respondIfError(err, w, fmt.Sprintf("Failed to retrieve file metadata. Error: %v", err), err500) {
		return
	}
	name := con.To.ItemName
	if name == "" {
		name = con.To.ItemID
		if name == "" {
			name = item.Name()
		}
	}
	copiedItem, err := destContainer.Put(name, reader, size, metadata)
	if respondIfError(err, w, fmt.Sprintf("File transfer failed. Error: %v", err), err500) {
		return
	}

	resultItem := Item{
		ID:       copiedItem.ID(),
		Name:     copiedItem.Name(),
		Size:     size,
		URL:      copiedItem.URL().Host + copiedItem.URL().Path,
		Metadata: metadata,
	}
	respondJSON(w, resultItem, created201)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[Request] %s", r.URL.Path)
	con := Connection{}

	err = json.Unmarshal([]byte(r.FormValue("to")), &con)
	if respondIfError(err, w, fmt.Sprintf("Valid JSON required. Error: %v", err), err400) {
		return
	}

	file, fileHandle, err := r.FormFile("file")
	if respondIfError(err, w, fmt.Sprintf("Invalid file or file read failed. Error: %v", err), err400) {
		return
	}
	defer file.Close()

	log.Printf("Dialing destination: %s/%s", con.Kind, con.ContainerName)
	destLoc, err := dial(con.Kind, con.ConfigMap)
	if respondIfError(err, w, fmt.Sprintf("Connection to %s failed. %v", con.Kind, err), err500) {
		return
	}
	destContainer, err := destLoc.Container(con.ContainerName)
	if respondIfError(err, w, fmt.Sprintf("Failed to retrieve container. Error: %v", err), err500) {
		return
	}
	name := con.ItemName
	if name == "" {
		name = fileHandle.Filename
	}
	log.Printf("Uploading %s to %s", name, con.Kind)
	item, err := destContainer.Put(name, file, fileHandle.Size, map[string]interface{}{})
	if respondIfError(err, w, fmt.Sprintf("Failed to upload file to %s. Error: %v", con.Kind, err), err500) {
		return
	}
	log.Printf("Upload complete: %s to %s", name, con.Kind)
	respondJSON(w, item, created201)
}
