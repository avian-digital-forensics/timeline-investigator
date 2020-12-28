#!/usr/bin/env bash

# Set variables
GCP_KEY_FILE="/home/simon/Downloads/timeline-investigator-d12be0e0a07b.json" # NOTE: Change this to the correct location
GCR_CI_REG="https://eu.gcr.io"
CONTAINER_IMAGE_STAGING="eu.gcr.io/timeline-investigator/test-console:staging"

# Login to GCR, build and push the image
docker login -u _json_key -p "$(cat ${GCP_KEY_FILE})" ${GCR_CI_REG}
cd ../../
docker build . -f ./console.Dockerfile -t test-console:local
cd ./cmd/test-console
docker tag test-console:local ${CONTAINER_IMAGE_STAGING}-latest
docker push ${CONTAINER_IMAGE_STAGING}-latest

# Deploy to kubernetes
#kubectl set image deployment/test-console -n staging test-console=${CONTAINER_IMAGE_STAGING}-latest