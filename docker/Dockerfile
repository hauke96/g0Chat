# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /

# Build the server command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go build -o g0chat_server /src/g0Ch@_server/*.go

# Run the server command by default when the container starts.
ENTRYPOINT ./g0chat_server -p=44494

# Document that the service listens on port 44494.
EXPOSE 44494
