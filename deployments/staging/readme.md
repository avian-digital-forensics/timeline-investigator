# staging

The staging environment is used to stage new features and bug fixes to the Timeline Investigator.

## deploy

```bash
    # Set variables for the firebase-secrets
    export AUTH_CREDENTIALS_FILE=<path>
    export AUTH_API_KEY=<key>
    
    # Create elasticsearch-instance (fscrawler & ti-api are dependent)
    kubectl apply -f ./elasticsearch/.

    # Create fscrawler-instance (ti-api is dependent)
    kubectl apply -f ./fscrawler/.

    # Create secret for ti-api
    kubectl create secret generic auth \
        --from-file=credentials-file=${AUTH_CREDENTIALS_FILE} \
        --from-literal=api-key=${AUTH_API_KEY} \
        -n staging

    # Create ti-api-instance
    kubectl apply -f ./ti-api/.
```