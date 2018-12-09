#!/bin/bash

docker build -t dev-hwsc-document-svc .
docker tag dev-hwsc-document-svc hwsc/dev-hwsc-document-svc
docker push hwsc/dev-hwsc-document-svc