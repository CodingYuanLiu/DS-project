FROM golang:1.14
WORKDIR /go/src/app
COPY . .
RUN go env -w GOPROXY=https://goproxy.cn,direct && go env -w GO111MODULE=on
RUN go build -o build/rpc_client ./Client/rpc_client.go

ENTRYPOINT ["/go/src/app/build/rpc_client"]