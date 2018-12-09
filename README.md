# Cloud Connector
Connect to various cloud storage providers, transfer files between them or upload.

## How to start application
If using pre-built executable binary from ./builds directory, navigate to the project root directory (where main.go resides) start application in the command line as follows:\
```./builds/builds/cloud-connector-ubuntu64```

## How to build
Building Cloud Connector requires Golang application to be installed and GOROOT & GOPATH set up properly. If you haven't already, do so by follwoing instructions from http://golang.org. 

Now retrieve dependent packages required by the application.
```bash
go get -v github.com/graymeta/stow
go get -v github.com/graymeta/stow/azure
go get -v github.com/graymeta/stow/google
go get -v github.com/graymeta/stow/s3
```
After that navigate to the project root directory and run the following command:\
```go build -o builds/cloud-connector```

PS: By default Cloud Connector is supports the following providers: Google Cloud Storage, Amazon AWS S3, Microsoft Azure. However, the Stow project supports other providers such as B2, Swift and Oracle. If support for those are required, 1) import the required packages in the Storage.go file and retrieve the package by runnig the ```go get``` command rebuild the executable binary as mentioned above.

## API Guideline

### /container - POST
Lists containers from a specific provider

##### JSON Body:

```js
{
    "kind": "",    // s3, google, azure
    "cursor": "",  // for pagination
    "count": 0,    // number of items to retrieve. 0 for all items
    "config_map": {
        // depends on the value provided in "kind"
        // see "config_map" below for specifics 
    }
}
```

### /items - POST
Lists objects/files from a specific container/bucket

##### JSON Body:
```js
{
    "kind": "",
    "container_name": "", // Name of the container/bucket
    "cursor": "",
    "count": 100,
    "config_map": {
    }
}
```

### /copy - POST
Transfer a single file from one provider to another. For multiple files, invoke separately

##### JSON Body:
```js
{
    "from": {
        "kind": "",
        "container_name": "", // Source container
        "item_id": "path/to/file.jpg",
        "config_map": {
        }
    },
    "to": {
        "kind": "",
        "container_name": "", // Destination container
        "item_name": "optional/new/path/to/uploaded-file.jpg",
        "config_map": {
        }
    }
}
```

### /upload - POST
Upload file directly to container.
The request must be created using Form Data. Refer to /public/index.html for example.

##### Form Data:
```js
{
    "file": (file binary),
    "to": {
        "kind": "",
        "container_name": "",
        "item_name": "optional/new/path/to/uploaded-file.jpg",
        "config_map": {
        }
    }
}
```

### /getjsonstring
Returns supplied JSON object/array as a single string value by escaping required characters. Can be usesful for the "json" attribute of the "config_map" for connecting to Google Cloud Storage account.

##### JSON Body:
```js
{
    "key1": "value",
    "key2": "value",
    "key3": "value"
}
```

## config_map
config_map properties for different providers
### config_map for AWS S3
```js
{
    "access_key_id": "",
    "secret_key":   "",
    "region":      ""   // example: us-east-1
}
```
### config_map for Google Cloud Storage
```js
{
    "project_id": "", // project id
    "json":""         // entire json file downloaded from Google as string with double quotes escaped using backslash. If unsure use the API function "/getjsonstring" below
}
```
### config_map for Azure
```js
{
    "account": "",
    "key":   ""
}
```
### config_map for other providers
Please refer to the supported providers https://github.com/graymeta/stow. Required config_map attributes can be found in the confg.go file of the respective package insite the "const" declarations