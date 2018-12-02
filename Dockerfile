FROM golang:1.11.1
WORKDIR $GOPATH/src/github.com/faraonc/
RUN git clone https://github.com/faraonc/hwsc-document-svc.git
ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep
WORKDIR $GOPATH/src/github.com/faraonc/hwsc-document-svc
RUN dep ensure -v
RUN go install
ENTRYPOINT ["/go/bin/hwsc-document-svc"]
EXPOSE 50051