# timeline-investigator

## ti-api

### requirements
 - service-account for Firebase Admin SDK
 - elasticsearch
 - fscrawler rest-api (with the same elasticsearch cluster)

### build

`CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api ./cmd/main/main.go`

### start the api

Copy the config from `./.local/config.local.yml` and change for your needs.

`./api --cfg=/path/to/config.yml`

### test the api

We have go-written tests in: `./tests` - to run these you need to have the API running with test enabled, test-secret and the API-key from Firebase for authentication. The API-key is needed for the API-server to create test-users in the system.

```yaml
config:
  test:
    run: true
    secret: super-secret
  authentication:
    api_key: 23jr8hf49q8f # api-key from firebase
```

To perform requests against the API an authorization header with a valid JWToken is required. The token can only be created from two places, 

1. Frontend authentication with firebase
2. With the [TestService](https://github.com/avian-digital-forensics/timeline-investigator/tree/main/pkg/services#testservice) (if test-run is enabled for the API). 

___

## improvements

There are mainly two parts that needs improvement and those haven't been prioritized yet because of limited time during development.

### datastore (./pkg/datastore)

The datastore is using elasticsearch as an index-engine, this is not currently done in an efficient way. We need to create mapping-, index- and search-functions to the datastore-service.

### fscrawler (./pkg/fscrawler)

Currently, we are using fscrawlers REST-interface for indexing OCRed documents, this is not optimal in our use case, the reason is that we want to index containers, for example, PST-files, which fscrawler cannot do. One of the solutions for that is to create a seperate API (for horizontal scaling) that will handle extracting with, for example, [libpff](https://github.com/libyal/libpff) and then do the indexing with fscrawler. Another thing is to mount the filestore-volume for the TI-API with the same volume for fscrawler to keep the data consistent.

