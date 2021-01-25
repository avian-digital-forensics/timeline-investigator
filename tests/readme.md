# tests

## Write tests

Our client-API for tests is generated to `./client/client.gen.go` from `../def/generate/client.go.plush`, and is used to write integrations-tests against our API-services.

Use TestService to perform unauthenticated requests to the API during testing for creating data.

## Run tests

A running instance of the TI-API with testing enabled is required for testing.

Required environment variables:
```bash
    # TEST_URL for where the TI-API is listening at
    export TEST_URL=http://localhost:8080/api/
    # TEST_SECRET for the secret to use when testing the API
    export TEST_SECRET=super-secret
```

Run the tests: 
```bash 
go test -v
```