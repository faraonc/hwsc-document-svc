#!/bin/bash

# Ensure you have a Docker Hub account and belongs to hwsc organization
docker build -t dev-hwsc-document-svc .
docker tag dev-hwsc-document-svc hwsc/dev-hwsc-document-svc
docker push hwsc/dev-hwsc-document-svc
