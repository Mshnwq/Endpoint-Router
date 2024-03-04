#Dockerfile
############################
# STEP 1 build executable binary
############################
#From which image we want to build. This is basically our environment.
FROM golang:1.20-alpine as Build
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/mypackage/myapp/
# COPY . .
COPY *.go .
# Fetch dependencies.
RUN go mod init main
RUN go mod tidy
# Using go get.
RUN go get -d -v
# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /custom-router

############################
# STEP 2 build a small image
############################
FROM scratch
# Copy our static executable.
COPY --from=Build /custom-router /custom-router
COPY proxy.log /proxy.log
# Run the hello binary.
ENTRYPOINT ["./custom-router"]