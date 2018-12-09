# Cloud Connector
Connect to various cloud storage providers, transfer files between them or upload.

# API Guideline

## /container - POST
Lists containers from a specific provider

### JSON Body:

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

## /items - POST
Lists objects/files from a specific container/bucket

### JSON Body:
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

## /copy - POST
Transfer a single file from one provider to another. For multiple files, invoke separately

### JSON Body:
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

## /upload - POST
Upload file directly to container.
The request must be created using Form Data. Refer to /public/index.html for example.

### Form Data:
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

## /getjsonstring
Returns supplied JSON object/array as a single string value by escaping required characters. Can be usesful for the "json" attribute of the "config_map" for connecting to Google Cloud Storage account.

### JSON Body:
```js
{
    "key1": "value",
    "key2": "value",
    "key3": "value"
}
```

# config_map
config_map properties for different providers
## config_map for AWS S3
```js
{
    "access_key_id": "",
    "secret_key":   "",
    "region":      ""   // example: us-east-1
}
```
## config_map for Google Cloud Storage
```js
{
    "project_id": "", // project id
    "json":""         // entire json file downloaded from Google as string with double quotes escaped using backslash. If unsure use the API function "/getjsonstring" below
}
```
## config_map for Azure
```js
{
    "account": "",
    "key":   ""
}
```
## config_map for other providers
Please refer to the supported providers https://github.com/graymeta/stow. Required config_map attributes can be found in the confg.go file of the respective package insite the "const" declarations