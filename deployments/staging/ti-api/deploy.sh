#!/usr/bin/env bash

# Add the correct variables
AUTH_CREDENTIALS_FILE= # path for the credentials-file
AUTH_API_KEY= # API-key for firebase

kubectl create secret generic auth \
    --from-file=credentials-file=${AUTH_CREDENTIALS_FILE} \
    --from-literal=api-key=${AUTH_API_KEY}

kubectl apply -f .