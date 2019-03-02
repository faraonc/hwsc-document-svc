FROM golang:1.12.0
WORKDIR $GOPATH/
RUN git clone https://github.com/hwsc-org/hwsc-document-svc.git
WORKDIR $GOPATH/hwsc-document-svc
RUN go mod download
RUN go install
ENTRYPOINT ["/go/bin/hwsc-document-svc"]
EXPOSE 50051
