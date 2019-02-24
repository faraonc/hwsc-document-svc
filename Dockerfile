FROM golang:1.11.5
WORKDIR $GOPATH/
RUN git clone https://github.com/hwsc-org/hwsc-document-svc.git
WORKDIR $GOPATH/hwsc-document-svc
RUN go install
ENTRYPOINT ["/go/bin/hwsc-document-svc"]
EXPOSE 50051
